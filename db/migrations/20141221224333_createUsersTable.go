package main

import (
	"database/sql"
)

// Up is executed when this migration is applied
func Up_20141221224333(txn *sql.Tx) {
	create := `CREATE TABLE users (
			Id serial NOT NULL,
			Name varchar,
			HashedPassword varchar,
			Email varchar UNIQUE
		);`
	if _, err := txn.Exec(create); err != nil {
		panic(err)
	}
	index := `CREATE UNIQUE INDEX ON users (Email);
	CREATE UNIQUE INDEX ON users (Id);`
	if _, err := txn.Exec(index); err != nil {
		panic(err)
	}
}

// Down is executed when this migration is rolled back
func Down_20141221224333(txn *sql.Tx) {
	if _, err := txn.Exec("DROP TABLE users;"); err != nil {
		panic(err)
	}
}
