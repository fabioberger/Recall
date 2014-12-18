package main

import (
	"database/sql"
)

// Up is executed when this migration is applied
func Up_20141216110042(txn *sql.Tx) {
	create := `CREATE TABLE reminders (
		Id serial NOT NULL,
		Reminder varchar UNIQUE,
		Timestamp int NOT NULL
	);`
	if _, err := txn.Exec(create); err != nil {
		panic(err)
	}
	index := `CREATE UNIQUE INDEX ON reminders (Id);`
	if _, err := txn.Exec(index); err != nil {
		panic(err)
	}
}

// Down is executed when this migration is rolled back
func Down_20141216110042(txn *sql.Tx) {
	if _, err := txn.Exec("DROP TABLE reminders;"); err != nil {
		panic(err)
	}
}
