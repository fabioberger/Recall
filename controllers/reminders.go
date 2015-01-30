package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/albrow/go-data-parser"
	"github.com/fabioberger/recall/models"
	"github.com/gorilla/mux"
	"github.com/sendgrid/sendgrid-go"
	"github.com/unrolled/render"
)

type Reminders struct {
	Mock bool
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

// CreateReminder takes a reminder form submission and saves it to the DB
func (rs Reminders) Create(res http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{
		IndentJSON: true,
	})

	// Check user is loggedin
	user, err := getCurrentUser(req)
	if err != nil {
		e := NewRequestError(err.Error())
		r.JSON(res, http.StatusOK, e)
		return
	}
	if user == nil {
		e := NewRequestError("Please Login/Signup to create a reminder")
		r.JSON(res, http.StatusOK, e)
		return
	}

	reminderData, err := data.Parse(req)
	checkErr(err, "Reminder POST data parse error")

	val := reminderData.Validator()
	val.Require("Reminder")
	val.LengthRange("Reminder", 3, 100)
	if val.HasErrors() {
		e := NewRequestError(strings.Join(val.Messages(), " "))
		r.JSON(res, http.StatusOK, e)
		return
	}

	reminder := &models.Reminder{
		Reminder:  reminderData.Get("Reminder"),
		Timestamp: int32(time.Now().Unix()),
		Userid:    int32(user.Id),
	}

	err = reminder.Save()
	checkErr(err, "Inserting reminder failed")

	r.JSON(res, http.StatusOK, reminder)
}

func (rs Reminders) Delete(res http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["reminder_id"])
	if err != nil {
		fmt.Printf("%v Is not a number.\n", id)
		return
	}

	err = models.RemoveById(id)
	if err != nil {
		fmt.Println(err)
	}
}

func (rs Reminders) GetAllForCurrentUser(res http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{
		IndentJSON: true,
	})

	user, err := getCurrentUser(req)
	if err != nil {
		// Is there a more suscinct way to create list of error returns?
		e := NewRequestError("You must be logged in to see reminders")
		eList := []*RequestError{}
		eList = append(eList, e)
		r.JSON(res, http.StatusOK, eList)
		return
	}

	reminders := models.GetReminders(int32(user.Id))
	r.JSON(res, http.StatusOK, reminders)
}

func (rs Reminders) CheckAll() string {
	reminders := models.GetAllReminders()
	messages := ""
	for _, r := range reminders {
		newMsg := rs.Check(r)
		messages = newMsg + ", " + messages
	}
	return messages
}

// Check looks for reminders to be sent and removes old reminders
func (rs Reminders) Check(r models.Reminder) string {
	reminderSetAt := time.Unix(int64(r.Timestamp), 0)
	hrsElapsed := int(time.Now().Sub(reminderSetAt).Hours())
	if hrsElapsed == 24 {
		rs.Send("24 hour", r)
		return "24h"
	} else if hrsElapsed == 168 {
		rs.Send("1 week", r)
		return "1w"
	} else if hrsElapsed == 730 {
		rs.Send("1 month", r)
		return "1m"
	} else if hrsElapsed == 2191 {
		rs.Send("3 month", r)
		return "3m"
	} else if hrsElapsed > 2191 {
		err := r.Remove()
		checkErr(err, "Removing reminder failed")
		return "rm"
	}
	return "no reminder"
}

// Send uses to Sendgrid API to send recall emails
func (rs Reminders) Send(recallTime string, reminder models.Reminder) {
	// If its a mock reminder (i.e testing), don't actually send any emails
	// Get user's email for reminder
	if rs.Mock == true {
		return
	}
	user, err := models.FindUserById(reminder.Userid)
	if err != nil {
		return
	}
	sg := sendgrid.NewSendGridClient(os.Getenv("SENDGRID_USER"), os.Getenv("SENDGRID_KEY"))
	message := sendgrid.NewMail()
	message.AddTo(user.Email)
	message.AddToName("Fabio Berger")
	message.SetSubject("Recall: " + reminder.Reminder)
	message.SetText(recallTime + " Reminder: " + reminder.Reminder)
	message.SetFrom("me@fabioberger.com")
	err = sg.Send(message)
	checkErr(err, "Sending Reminder Email Failed")
}
