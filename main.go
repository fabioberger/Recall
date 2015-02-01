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
	"github.com/martini-contrib/cors"

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
	frontend_path := os.Getenv("GOPATH") + "/src/github.com/fabioberger/recall-frontend"
	s := negroni.NewStatic(http.Dir(frontend_path))
	n.Use(s)
	n.UseHandler(cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	n.Use(recovery.JSONRecovery(config.Env != "production"))
	store := cookiestore.New([]byte(config.Secret))
	n.Use(sessions.Sessions("recall_session", store))

	// Routes
	router := mux.NewRouter()

	reminders := controllers.Reminders{
		Mock: false,
	}

	users := controllers.Users{}
	sess := controllers.Sessions{}
	router.HandleFunc("/users", users.Create).Methods("POST")
	router.HandleFunc("/users/{user_id}", users.GetOne).Methods("GET")
	router.HandleFunc("/sessions", sess.Create).Methods("POST")
	router.HandleFunc("/sessions/{session_id}", sess.Delete).Methods("DELETE")
	router.HandleFunc("/reminders", reminders.GetAllForCurrentUser).Methods("GET")
	router.HandleFunc("/reminders", reminders.Create).Methods("POST")
	router.HandleFunc("/reminders/{reminder_id}", reminders.Delete).Methods("DELETE")

	n.UseHandler(router)
	n.Run(":" + config.Port)
}
