package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/albrow/negroni-json-recovery"
	"github.com/codegangsta/negroni"
	"github.com/fabioberger/recall/config"
	"github.com/fabioberger/recall/controllers"
	"github.com/fabioberger/recall/models"

	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
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

	config.Init()
	models.Init()

	// Middleware
	n := negroni.New(negroni.NewLogger())
	s := negroni.NewStatic(http.Dir("public"))
	n.Use(s)
	n.Use(recovery.JSONRecovery(config.Env != "production"))
	store := cookiestore.New([]byte(config.Secret))
	n.Use(sessions.Sessions("recall_session", store))

	// Routes
	router := mux.NewRouter()

	reminders := controllers.Reminders{
		Mock: false,
	}
	router.HandleFunc("/", reminders.GetAll).Methods("GET")
	router.HandleFunc("/reminder", reminders.Create).Methods("POST")
	router.HandleFunc("/reminder", reminders.Delete).Methods("DELETE")

	users := controllers.Users{}
	router.HandleFunc("/login", users.Login).Methods("GET")
	sess := controllers.Sessions{}
	router.HandleFunc("/login", sess.Create).Methods("POST")
	router.HandleFunc("/login", sess.Delete).Methods("DELETE")
	router.HandleFunc("/signup", users.Signup).Methods("GET")
	router.HandleFunc("/signup", users.Create).Methods("POST")
	router.HandleFunc("/profile", users.Profile).Methods("GET")

	n.UseHandler(router)
	n.Run(":" + config.Port)
}
