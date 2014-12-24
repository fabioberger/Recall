package controllers

import (
	"fmt"
	"net/http"

	"code.google.com/p/go.crypto/bcrypt"
	"github.com/albrow/go-data-parser"
	"github.com/fabioberger/recall/models"
	"github.com/goincremental/negroni-sessions"
	"github.com/unrolled/render"
)

type Users struct{}

func (u Users) Profile(res http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{
		Layout: "layout",
	})

	user, err := getCurrentUser(req)
	if err != nil {
		r.HTML(res, 200, "login", err)
		return
	}
	if user == nil {
		r.HTML(res, 200, "login", "Please Login")
		return
	}
	fmt.Println("HERE!")
	fmt.Println("userid: ", user.Id)
	reminders := models.GetReminders(int32(user.Id))
	readableReminders := models.MakeReadableReminders(reminders)

	profile := struct {
		Name      string
		Email     string
		Reminders []models.ReadableReminder
	}{
		Name:      user.Name,
		Email:     user.Email,
		Reminders: readableReminders,
	}

	r.HTML(res, 200, "profile", profile)
}

func (u Users) Create(res http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{
		Layout: "layout",
	})

	// Parse data from form
	userData, err := data.Parse(req)
	if err != nil {
		panic(err)
	}

	// Validate
	val := userData.Validator()
	val.Require("name")
	val.LengthRange("name", 3, 35)
	val.Require("email")
	val.MatchEmail("email")
	val.Require("password")
	val.MinLength("password", 8)
	val.Require("confirmPassword")
	val.Equal("password", "confirmPassword")
	email, name := userData.Get("email"), userData.Get("name")
	models.ValidateUserUnique(val, email, name)
	if val.HasErrors() {
		r.HTML(res, 200, "signup", val.Messages())
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword(userData.GetBytes("password"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	// Create and save user
	user := &models.User{
		Name:           name,
		Email:          email,
		HashedPassword: string(hashedPassword),
	}
	if err := user.Save(); err != nil {
		panic(err)
	}

	sess := Sessions{}
	s := sessions.GetSession(req)
	sess.CreateFromCredentials(s, user.Email, string(userData.GetBytes("password")), false)
	sess.SetLoggedInCookie(res, "true")

	// Redirect to profile
	http.Redirect(res, req, "/profile", 301)
}

func (u Users) Login(res http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{
		Layout: "layout",
	})
	r.HTML(res, 200, "login", nil)
}

func (u Users) Signup(res http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{
		Layout: "layout",
	})
	r.HTML(res, 200, "signup", nil)
}
