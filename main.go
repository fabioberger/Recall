package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/fabioberger/recall/controllers"
	"github.com/fabioberger/recall/models"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const version = "0.0.0"

var versionFlag = flag.Bool("version", false, "If provided, server will print out the version number then immediately exit.")

func main() {
	flag.Parse()
	if *versionFlag {
		// print out the version and immediately exit
		fmt.Println(version)
		os.Exit(0)
	}

	// Middleware
	n := negroni.New(negroni.NewLogger())

	// Routes
	router := mux.NewRouter()

	models.Init()

	reminders := controllers.Reminders{}
	router.HandleFunc("/", reminders.GetAll).Methods("GET")
	router.HandleFunc("/reminder", reminders.Create).Methods("POST")
	router.HandleFunc("/check", reminders.Check).Methods("GET")

	n.UseHandler(router)
	n.Run(":3000")
}
