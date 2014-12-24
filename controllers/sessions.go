package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"code.google.com/p/go.crypto/bcrypt"

	"github.com/albrow/go-data-parser"
	"github.com/fabioberger/recall/models"
	"github.com/goincremental/negroni-sessions"
	"github.com/unrolled/render"
)

const AuthenticationKey = "recallAuth"

type Sessions struct{}

func (s Sessions) Create(res http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{
		Layout: "layout",
	})
	session := sessions.GetSession(req)
	sessionData, err := data.Parse(req)
	if err != nil {
		panic(err)
	}

	// Validations
	val := sessionData.Validator()
	val.Require("email")
	val.Require("password")
	if val.HasErrors() {
		r.HTML(res, 200, "login", val.Messages())
		return
	}

	// Get user from DB
	email := sessionData.Get("email")
	user, err := models.FindUserByEmail(email)
	if err != nil {
		panic(err)
	} else if user == nil {
		r.HTML(res, 200, "login", "email and/or password incorrect")
		return
	}

	// Make sure password matches
	password := sessionData.Get("password")
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		r.HTML(res, 200, "login", "email and/or password was incorrect.")
		return
	}

	s.CreateFromCredentials(session, email, password, sessionData.GetBool("rememberMe"))
	s.SetLoggedInCookie(res, "true")

	http.Redirect(res, req, "/profile", 301)
}

func (s Sessions) CreateFromCredentials(session sessions.Session, email string, password string, rememberMe bool) {
	// Set data in the session store
	session.Set(AuthenticationKey, generateAuthString(email, password))
	if rememberMe {
		session.Options(sessions.Options{
			MaxAge: 60 * 60 * 24 * 30, // 30 days
		})
	} else {
		session.Options(sessions.Options{
			MaxAge: 0,
		})
	}
}

func (s Sessions) Delete(res http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{
		Layout: "layout",
	})
	session := sessions.GetSession(req)
	session.Delete(AuthenticationKey)
	s.SetLoggedInCookie(res, "false")
	r.HTML(res, 200, "home", nil)
}

// authString consists of email and password separated by a colon
// this works because valid email cannot contain a colon
func generateAuthString(email string, password string) string {
	return email + ":" + password
}

func parseAuthString(authString string) (email string, password string) {
	authSlice := strings.SplitN(authString, ":", 2)
	return authSlice[0], authSlice[1]
}

// gets the current user by looking at session data
// any error returned here should be considered a 400 error
func getCurrentUser(req *http.Request) (*models.User, error) {
	// Parse session data
	session := sessions.GetSession(req)
	s := session.Get(AuthenticationKey)
	fmt.Println(s)
	if s == nil {
		return nil, nil
	}
	authString, ok := s.(string)
	if !ok {
		panic("Couldn't convert session to string")
	}
	email, password := parseAuthString(authString)

	// Get user from database
	user, err := models.FindUserByEmail(email)
	if err != nil {
		panic(err)
	} else if user == nil {
		return nil, errors.New("email and/or password was incorrect.")
	}

	// Make sure password matches
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return nil, errors.New("email and/or password was incorrect.")
	}

	// If we've reached here, the credentials are correct and we can return the
	// current user
	return user, nil
}

func (s Sessions) SetLoggedInCookie(res http.ResponseWriter, status string) {
	cookie := http.Cookie{
		Name:  "isLoggedIn",
		Value: status,
	}
	http.SetCookie(res, &cookie)
}
