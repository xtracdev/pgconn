package pgconn

import (
	. "github.com/gucumber/gucumber"
	"github.com/xtracdev/pgconn"
	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"time"
	"errors"
)

var connectStr = ""
var maskedConnectStr = ""

var bogusConnectStr = "user=luser password=passw0rd dbname=postgres host=localhost port=15432 sslmode=disable"



func init() {
	envConnect,envError := pgconn.NewEnvConfig()
	if envError != nil {
		log.Fatal("No config specified for gucumber test")
	}


	connectStr = envConnect.ConnectString()
	maskedConnectStr = envConnect.MaskedConnectString()

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
		db, err = pgconn.OpenAndConnect(connectStr, 10)
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
		db, noConnectError = pgconn.OpenAndConnect(bogusConnectStr, 3)
	})

	Then(`^an error is returned$`, func() {
		assert.NotNil(T, noConnectError)
	})

	Given(`^a loss of database connectivity$`, func() {
		var err error
		db, err = pgconn.OpenAndConnect(connectStr, 10)
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
