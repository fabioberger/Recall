package controllers

import (
	"errors"
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

type RequestError struct {
	Error string
}

func NewRequestError(msg string) *RequestError {
	e := new(RequestError)
	e.Error = msg
	return e
}

func (s Sessions) Create(res http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{
		IndentJSON: true,
	})
	session := sessions.GetSession(req)
	sessionData, err := data.Parse(req)
	if err != nil {
		panic(err)
	}

	// Validations
	val := sessionData.Validator()
	val.Require("Email")
	val.MatchEmail("Email")
	val.Require("Password")
	if val.HasErrors() {
		sessErr := NewRequestError(strings.Join(val.Messages(), " "))
		r.JSON(res, http.StatusOK, sessErr)
		return
	}

	// Get user from DB
	email := sessionData.Get("Email")
	user, err := models.FindUserByEmail(email)
	if err != nil {
		panic(err)
	} else if user == nil {
		sessErr := NewRequestError("email and/or password is incorrect")
		r.JSON(res, http.StatusOK, sessErr)
		return
	}

	// Make sure password matches
	password := sessionData.Get("Password")
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		sessErr := NewRequestError("email and/or password is incorrect")
		r.JSON(res, http.StatusOK, sessErr)
		return
	}
	rememberMe := sessionData.GetBool("rememberMe")
	s.CreateFromCredentials(session, email, password, rememberMe)
	s.SetLoggedInCookie(res, "true")

	frontendSession := struct {
		Email      string
		RememberMe bool
	}{
		Email:      email,
		RememberMe: rememberMe,
	}
	r.JSON(res, http.StatusOK, frontendSession)
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
	session := sessions.GetSession(req)
	session.Delete(AuthenticationKey)
	s.SetLoggedInCookie(res, "false")
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
	if s == nil {
		return nil, errors.New("You are currently not logged in.")
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
