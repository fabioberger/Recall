package models

import (
	"database/sql"
	"fmt"

	"github.com/coopernurse/gorp"
	"github.com/fabioberger/golympus/config"
	_ "github.com/lib/pq"
)

var Db *gorp.DbMap

func Init() {
	// connect to postgres database
	// TODO: use a password here if needed
	fmt.Println("[database] Connecting to database...")
	db, err := sql.Open("postgres", "dbname=recall sslmode=disable")
	if err != nil {
		panic(err)
	}
	fmt.Println("[database] Connected successfully.")

	// construct a gorp DbMap
	Db = &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	// add a table, setting the table name to 'reminders' and
	// specifying that the Id property is an auto incrementing PK
	Db.AddTableWithName(Reminder{}, "reminders").SetKeys(true, "Id")

	if config.Env == "test" {
		// if we're in the test environment, clear the database on startup
		Db.TruncateTables()
	}
}
