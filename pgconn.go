package pgconn

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"
	"github.com/xtracdev/envinject"
)

const (
	maxConns         = "DB_MAX_OPEN_CONNS"
	idleConns        = "DB_MAX_IDLE_CONNS"
	defaultMaxConns  = 6
	defaultIdleConns = 1
)

type PostgresDB struct {
	*sql.DB
	injectedEnv *envinject.InjectedEnv
}

var ErrRetryCount = errors.New("Retry count must be greater than 1")



func OpenAndConnect(env *envinject.InjectedEnv, retryCount int) (*PostgresDB, error) {
	if retryCount < 1 {
		return nil, ErrRetryCount
	}

	connectString, err := ConnectStringFromInjectedEnv(env)
	if err != nil {
		return nil,err
	}

	log.Infof("Open the database, %d retries", retryCount)
	db, err := sql.Open("postgres", connectString)
	if err != nil {
		return nil, err
	}

	log.Info("Ping the db as open might not actually connect")

	var dbError error
	maxAttempts := retryCount
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		log.Info("ping database...")
		dbError = db.Ping()
		if dbError == nil {
			break
		}

		log.Infof("Ping failed: %s", strings.TrimSpace(dbError.Error()))
		log.Infof("Retry in %d seconds", attempts)
		time.Sleep(time.Duration(attempts) * time.Second)
	}
	if dbError != nil {
		return nil, dbError
	}

	pgdb := &PostgresDB{DB: db, injectedEnv: env}
	pgdb.setMaxOpenConns()
	pgdb.setMaxIdleConns()

	return pgdb, nil
}

//Reconnect to the database. Useful when a loss of connection has been detected
func (pgdb *PostgresDB) Reconnect(retryCount int) error {
	pgdb.Close()
	db, err := OpenAndConnect(pgdb.injectedEnv, retryCount)
	if err != nil {
		return err
	}

	pgdb.DB = db.DB
	return nil
}

func (pgdb *PostgresDB) getIntFromEnv(varName string, defaultVal int) int {
	var val = defaultVal
	env := pgdb.injectedEnv.Getenv(varName)
	if env != "" {
		var err error
		val, err = strconv.Atoi(env)
		if err != nil {
			log.Infof("Failed to convert %s value, setting default value", varName)
			val = defaultVal
		}
	}

	return val
}

func (pgdb *PostgresDB) getMaxConns() int {
	return pgdb.getIntFromEnv(maxConns,defaultMaxConns)
}

func (pgdb *PostgresDB) getIdleConns() int {
	return pgdb.getIntFromEnv(idleConns,defaultIdleConns)
}


func (pgdb *PostgresDB) setMaxOpenConns() {
	var max = pgdb.getMaxConns()
	log.Infof("Setting %s to %d connections...", maxConns, max)
	pgdb.DB.SetMaxOpenConns(max)
}

func (pgdb *PostgresDB) setMaxIdleConns() {
	var idle = pgdb.getIdleConns()
	log.Infof("Setting %s to %d connections...", idleConns, idle)
	pgdb.DB.SetMaxIdleConns(idle)
}

//BuildConnectString builds an Oracle connect string from its constituent parts.
func BuildConnectString(user, password, host, port, dbName string) string {
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		user, password, dbName, host, port)
}

//IsConnectionError returns error if the argument is a connection error
func IsConnectionError(err error) bool {
	//Observed empirically...
	errStr := err.Error()
	return strings.Contains(errStr, "connection refused")
}
