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

type Reminders struct{}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

// CreateReminder takes a reminder form submission and saves it to the DB
func (rs Reminders) Create(res http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{})

	reminderData, err := data.Parse(req)
	checkErr(err, "POST data parse error")

	val := reminderData.Validator()
	val.Require("reminder")
	val.LengthRange("reminder", 3, 100)

	reminder := &models.Reminder{
		Reminder:  reminderData.Get("reminder"),
		Timestamp: int32(time.Now().Unix()),
	}

	err = reminder.Save()
	checkErr(err, "Inserting reminder failed")

	r.HTML(res, 200, "index", "success")
}

// GetAll returns all existing reminders and new reminder form
func (rs Reminders) GetAll(res http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{})
	reminders := models.GetAllReminders()

	r.HTML(res, 200, "index", reminders)
}

// Check looks for reminders to be sent and removes old reminders
func (rs Reminders) Check() {
	reminders := models.GetAllReminders()
	for _, r := range reminders {
		reminderSetAt := time.Unix(int64(r.Timestamp), 0)
		hrsElapsed := int(time.Now().Sub(reminderSetAt).Hours())
		if hrsElapsed == 24 {
			rs.Send("24 hour", r.Reminder)
		} else if hrsElapsed == 168 {
			rs.Send("1 week", r.Reminder)
		} else if hrsElapsed == 730 {
			rs.Send("1 month", r.Reminder)
		} else if hrsElapsed == 2191 {
			rs.Send("3 month", r.Reminder)
		} else if hrsElapsed > 2191 {
			err := r.Remove()
			checkErr(err, "Removing reminder failed")
		}
	}
}

// Send uses to Sendgrid API to send recall emails
func (rs Reminders) Send(recallTime string, reminder string) {
	sg := sendgrid.NewSendGridClient(os.Getenv("SENDGRID_USER"), os.Getenv("SENDGRID_KEY"))
	message := sendgrid.NewMail()
	message.AddTo(os.Getenv("GMAIL_EMAIL"))
	message.AddToName("Fabio Berger")
	message.SetSubject("Recall: " + reminder)
	message.SetText(recallTime + " Reminder: " + reminder)
	message.SetFrom("me@fabioberger.com")
	err := sg.Send(message)
	checkErr(err, "Sending Reminder Email Failed")
}
