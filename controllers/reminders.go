package controllers

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/albrow/go-data-parser"
	"github.com/fabioberger/recall/models"
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
		Layout: "layout",
	})

	// Check user is loggedin
	user, err := getCurrentUser(req)
	if err != nil {
		r.HTML(res, 200, "login", err)
		return
	}
	if user == nil {
		r.HTML(res, 200, "login", "Please Login/Signup to create a reminder")
		return
	}

	reminderData, err := data.Parse(req)
	checkErr(err, "POST data parse error")

	val := reminderData.Validator()
	val.Require("reminder")
	val.LengthRange("reminder", 3, 100)

	reminder := &models.Reminder{
		Reminder:  reminderData.Get("reminder"),
		Timestamp: int32(time.Now().Unix()),
		Userid:    int32(user.Id),
	}

	err = reminder.Save()
	checkErr(err, "Inserting reminder failed")

	http.Redirect(res, req, "/profile", 301)
}

// GetAll returns all existing reminders and new reminder form
func (rs Reminders) GetAll(res http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{
		Layout: "layout",
	})
	reminders := models.GetAllReminders()

	r.HTML(res, 200, "home", reminders)
}

func (rs Reminders) CheckAll() string {
	reminders := models.GetAllReminders()
	messages := ""
	for _, r := range reminders {
		messages = rs.Check(r) + ", " + messages
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
