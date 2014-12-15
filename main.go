package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/albrow/go-data-parser"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	_ "github.com/lib/pq"
)

type Reminder struct {
	Id       int    `db:"id"`
	Reminder string `db:"reminder"`
	Sent     int    `db:"sent"`
}

func initDb() *gorp.DbMap {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	db, err := sql.Open("postgres", "dbname=recall sslmode=disable")
	checkErr(err, "sql.Open failed")

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	// add a table, setting the table name to 'posts' and
	dbmap.AddTableWithName(Reminder{}, "reminders").SetKeys(true, "Id")

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

func main() {
	m := martini.Classic()
	m.Map(initDb())
	m.Use(render.Renderer())

	m.Get("/", showForm)
	m.Post("/reminder", createReminder)
	m.Run()
}

func createReminder(r render.Render, req *http.Request, db *gorp.DbMap) {
	reminderData, err := data.Parse(req)
	checkErr(err, "POST data parse error")

	val := reminderData.Validator()
	val.Require("reminder")
	val.LengthRange("reminder", 3, 100)

	reminder := Reminder{
		Reminder: reminderData.Get("reminder"),
		Sent:     0,
	}

	err = db.Insert(&reminder)
	checkErr(err, "Inserting reminder failed")

	r.HTML(200, "index", "success")

}

func showForm(r render.Render, db *gorp.DbMap) {

	var reminders []Reminder
	_, err := db.Select(&reminders, "select * from reminders")
	checkErr(err, "error selecting reminders")
	r.HTML(200, "index", reminders)
}
