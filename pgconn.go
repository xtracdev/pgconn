package pgconn

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"
)

const (
	maxConns         = "DB_MAX_OPEN_CONNS"
	idleConns        = "DB_MAX_IDLE_CONNS"
	defaultMaxConns  = 6
	defaultIdleConns = 3
)

type PostgresDB struct {
	*sql.DB
	connectStr string
}

var ErrRetryCount = errors.New("Retry count must be greater than 1")

func OpenAndConnect(connectString string, retryCount int) (*PostgresDB, error) {
	if retryCount < 1 {
		return nil, ErrRetryCount
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

	pgdb := &PostgresDB{DB: db, connectStr: connectString}
	pgdb.SetMaxOpenConns()
	pgdb.SetMaxIdleConns()

	return pgdb, nil
}

//Reconnect to the database. Useful when a loss of connection has been detected
func (pgdb *PostgresDB) Reconnect(retryCount int) error {
	pgdb.Close()
	db, err := OpenAndConnect(pgdb.connectStr, retryCount)
	if err != nil {
		return err
	}

	pgdb.DB = db.DB
	return nil
}

func (pgdb *PostgresDB) SetMaxOpenConns() {
	var max int
	max = defaultMaxConns

	env := os.Getenv(maxConns)
	if env != "" {
		var err error
		max, err = strconv.Atoi(env)
		if err != nil {
			log.Infof("Failed to convert %s value, setting default value", maxConns)
			max = defaultMaxConns
		}
	}
	log.Infof("Setting %s to %d connections...", maxConns, max)
	pgdb.DB.SetMaxOpenConns(max)
}

func (pgdb *PostgresDB) SetMaxIdleConns() {
	var idle int
	idle = defaultIdleConns

	env := os.Getenv(idleConns)
	if env != "" {
		var err error
		idle, err = strconv.Atoi(env)
		if err != nil {
			log.Infof("Failed to convert %s value, setting default value", idleConns)
			idle = defaultIdleConns
		}
	}
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
