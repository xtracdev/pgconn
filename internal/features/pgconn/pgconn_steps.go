package pgconn

import (
	. "github.com/gucumber/gucumber"
	"github.com/xtracdev/pgconn"
	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"time"
)

var connectStr = ""
var maskedConnectStr = ""





func init() {
	envConnect,_ := pgconn.NewEnvConfig()
	connectStr = envConnect.ConnectString()
	maskedConnectStr = envConnect.MaskedConnectString()

	var db *pgconn.PostgresDB
	//var noConnectError error

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
		T.Skip() // pending
	})

	When(`^I connect to no listener$`, func() {
		T.Skip() // pending
	})

	Then(`^an error is returned$`, func() {
		T.Skip() // pending
	})

	Given(`^a loss of database connectivity$`, func() {
		T.Skip() // pending
	})

	When(`^I detect I've lost connectivity$`, func() {
		T.Skip() // pending
	})

	Then(`^I can reconnect$`, func() {
		T.Skip() // pending
	})

	And(`^I can select data after reconnecting$`, func() {
		T.Skip() // pending
	})

}
