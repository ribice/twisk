package pgsql

import (
	"log"
	"time"

	"github.com/go-pg/pg"
)

// New creates new database connection to a postgres database
// Function panics if it can't connect to database
func New(psn string, logQueries bool, timeout int) (*pg.DB, error) {
	u, err := pg.ParseURL(psn)
	if err != nil {
		return nil, err
	}

	db := pg.Connect(u)

	_, err = db.Exec("SELECT 1")
	if err != nil {
		return nil, err
	}

	if timeout > 0 {
		db.WithTimeout(time.Second * time.Duration(timeout))
	}

	if logQueries {
		db.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
			query, err := event.FormattedQuery()
			checkErr(err)
			log.Printf("%s | %s", time.Since(event.StartTime), query)
		})
	}

	return db, nil
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
