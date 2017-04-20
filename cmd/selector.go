package main

import (
	"github.com/xtracdev/pgconn"
	"time"
	"log"
)


func main() {
	envConnect,_ := pgconn.NewEnvConfig()
	connectStr := envConnect.ConnectString()

	db, err := pgconn.OpenAndConnect(connectStr, 10)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Print("connected")

	for {
		log.Print("select...")
		rows, err := db.Query("select now from now()")
		if err == nil {

			for rows.Next() {
				var ts time.Time
				rows.Scan(&ts)
				log.Printf("select from now is %s", ts.Format(time.RFC3339))
			}

			rows.Close()
		} else {
			log.Printf(err.Error())
		}

		time.Sleep(5*time.Second)
	}
}