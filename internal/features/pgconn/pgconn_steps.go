package pgconn

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	. "github.com/gucumber/gucumber"
	"github.com/stretchr/testify/assert"
	"github.com/xtracdev/envinject"
	"github.com/xtracdev/pgconn"
	"os"
	"time"
)

var connectStr = ""
var maskedConnectStr = ""

var bogusConnectStr = "user=luser password=passw0rd dbname=postgres host=localhost port=15432 sslmode=disable"

func init() {
	os.Setenv(envinject.ParamPrefixEnvVar, "")
	env, _ := envinject.NewInjectedEnv()

	maskedConnectStr, _ = pgconn.MaskedConnectStringFromInjectedEnv(env)

	var db *pgconn.PostgresDB
	var noConnectError error

	Given(`^a running postgres instance$`, func() {
		log.Infof("Postgres instance available via %s assumed", maskedConnectStr)
	})

	When(`^provide a connection string for the running instance$`, func() {
		//
	})

	Then(`^a connection is returned$`, func() {
		var err error
		db, err = pgconn.OpenAndConnect(env, 10)
		assert.Nil(T, err)
	})

	And(`^I can select the system timestamp$`, func() {
		rows, err := db.Query("select now from now()")
		if assert.Nil(T, err) {
			defer rows.Close()

			for rows.Next() {
				var ts time.Time
				rows.Scan(&ts)
				log.Infof("select from now is %s", ts.Format(time.RFC3339))
			}

			assert.Nil(T, rows.Err())
		}
	})

	Given(`^a connection string with no listener$`, func() {
		log.Infof("No pg instance available via %s assumed", bogusConnectStr)
	})

	When(`^I connect to no listener$`, func() {

		currentPort := os.Getenv(pgconn.DBPort)
		os.Setenv(pgconn.DBPort, "12345")
		bogusEnv, _ := envinject.NewInjectedEnv()

		db, noConnectError = pgconn.OpenAndConnect(bogusEnv, 3)
		os.Setenv(pgconn.DBPort, currentPort)
	})

	Then(`^an error is returned$`, func() {
		assert.NotNil(T, noConnectError)
	})

	Given(`^a loss of database connectivity$`, func() {
		var err error
		db, err = pgconn.OpenAndConnect(env, 10)
		if assert.Nil(T, err) {
			err = db.Close()
			assert.Nil(T, err)
		}
	})

	When(`^I detect I've lost connectivity$`, func() {
		assert.True(T, pgconn.IsConnectionError(errors.New("connection refused")), "Expected a connection error")
	})

	Then(`^I can reconnect$`, func() {
		err := db.Reconnect(3)
		assert.Nil(T, err)
	})

	And(`^I can select data after reconnecting$`, func() {
		rows, err := db.Query("select now from now()")
		if assert.Nil(T, err) {
			defer rows.Close()

			for rows.Next() {
				var ts time.Time
				rows.Scan(&ts)
				log.Infof("select from now is %s", ts.Format(time.RFC3339))
			}

			assert.Nil(T, rows.Err())
		}
	})

}
