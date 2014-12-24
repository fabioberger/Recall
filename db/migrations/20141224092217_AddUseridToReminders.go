package main

import (
	"database/sql"
)

// Up is executed when this migration is applied
func Up_20141224092217(txn *sql.Tx) {
	alter := `ALTER TABLE reminders ADD userid INT`
	if _, err := txn.Exec(alter); err != nil {
		panic(err)
	}
}

// Down is executed when this migration is rolled back
func Down_20141224092217(txn *sql.Tx) {
	if _, err := txn.Exec("ALTER TABLE reminders DROP userid;"); err != nil {
		panic(err)
	}
}
