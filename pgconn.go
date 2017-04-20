package pgconn

import (
	"database/sql"
	_ "github.com/lib/pq"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"strings"
	"time"
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

	return &PostgresDB{DB: db, connectStr: connectString}, nil
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

//BuildConnectString builds an Oracle connect string from its constituent parts.
func BuildConnectString(user, password, host, port, dbName string) string {
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		user, password, dbName, host, port)
}

//IsConnectionError returns error if the argument is a connection error
func IsConnectionError(err error) bool {
	errStr := err.Error()
	return strings.HasPrefix(errStr, "ORA-03114") || strings.HasPrefix(errStr, "ORA-03113")
}
