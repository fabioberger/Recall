package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/albrow/negroni-json-recovery"
	"github.com/codegangsta/negroni"
	"github.com/fabioberger/recall/config"
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
	n.Use(recovery.JSONRecovery(config.Env != "production"))

	// Routes
	router := mux.NewRouter()

	config.Init()
	models.Init()

	reminders := controllers.Reminders{}
	router.HandleFunc("/", reminders.GetAll).Methods("GET")
	router.HandleFunc("/reminder", reminders.Create).Methods("POST")

	n.UseHandler(router)
	n.Run(":4000")
}
