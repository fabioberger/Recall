package controllers

import (
	"net/http"
	"strings"

	"code.google.com/p/go.crypto/bcrypt"
	"github.com/albrow/go-data-parser"
	"github.com/fabioberger/recall/models"
	"github.com/goincremental/negroni-sessions"
	"github.com/unrolled/render"
)

type Users struct{}

type FrontendUser struct {
	Id    int
	Name  string
	Email string
}

func (u Users) GetOne(res http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{
		IndentJSON: true,
	})

	user, err := getCurrentUser(req)
	if err != nil {
		e := NewRequestError(err.Error())
		r.JSON(res, http.StatusOK, e)
		return
	}
	if user == nil {
		e := NewRequestError("Please Login")
		r.JSON(res, http.StatusOK, e)
		return
	}

	FrontendUser := FrontendUser{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}
	r.JSON(res, http.StatusOK, FrontendUser)
}

func (u Users) Create(res http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{
		IndentJSON: true,
	})

	// Parse data from form
	userData, err := data.Parse(req)
	if err != nil {
		panic(err)
	}

	// Validate
	val := userData.Validator()
	val.Require("Name")
	val.LengthRange("Name", 3, 35)
	val.Require("Email")
	val.MatchEmail("Email")
	val.Require("Password")
	val.MinLength("Password", 8)
	val.Require("ConfirmPassword")
	val.Equal("Password", "ConfirmPassword")
	email, name := userData.Get("Email"), userData.Get("Name")
	models.ValidateUserUnique(val, email, name)
	if val.HasErrors() {
		e := NewRequestError(strings.Join(val.Messages(), " "))
		r.JSON(res, http.StatusOK, e)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword(userData.GetBytes("Password"), bcrypt.DefaultCost)
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
	sess.CreateFromCredentials(s, user.Email, string(userData.GetBytes("Password")), false)
	sess.SetLoggedInCookie(res, "true") //Still needed?

	frontendUser := FrontendUser{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}
	r.JSON(res, http.StatusOK, frontendUser)
}
