package main

import (
	"github.com/xtracdev/envinject"
	"github.com/xtracdev/pgconn"
	"log"
	"time"
)

func main() {
	envConnect, _ := envinject.NewInjectedEnv()

	db, err := pgconn.OpenAndConnect(envConnect, 10)
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

		time.Sleep(5 * time.Second)
	}
}
