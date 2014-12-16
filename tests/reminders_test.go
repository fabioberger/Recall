package tests

import (
	"fmt"
	"testing"

	"github.com/fabioberger/recall/config"
	"github.com/fabioberger/recall/models"
)

func TestReminderSave(t *testing.T) {
	defer tearDown()
	setUp()
	reminder := createReminder()
	reminder.Save()
	retrieved_reminder := &models.Reminder{}
	err := retrieved_reminder.FindByTimestamp(1418727704)
	if err != nil {
		fmt.Println(err)
	}
	compareString(t, "ReminderSave", reminder.Reminder, retrieved_reminder.Reminder)
	compareInt(t, "ReminderSave", int64(reminder.Timestamp), int64(retrieved_reminder.Timestamp))
}

func TestRemoveReminder(t *testing.T) {
	defer tearDown()
	setUp()
	reminder := createReminder()
	reminder.Save()
	timestamp := reminder.Timestamp
	reminder.Remove()
	retrieved_reminder := models.Reminder{}
	err := retrieved_reminder.FindByTimestamp(timestamp)
	if err == nil {
		t.Error("Reminder was not successfully removed")
	}
}

// tearDown clears the database and should be called at
// the end of every test, probably using defer.
func tearDown() {
	if models.Db == nil {
		config.Env = "test"
		config.Init()
		models.Init()
	}
	models.Db.TruncateTables()
}

// setUp loads the psql configs and init's the DB connection
func setUp() {
	config.Init()
	models.Init()
}

func createReminder() models.Reminder {
	reminder := models.Reminder{
		Reminder:  "Revise the Dominic Number System",
		Timestamp: 1418727704,
	}
	return reminder
}

func compareInt(t *testing.T, prefix string, expected int64, got int64) {
	if expected != got {
		t.Errorf("%s Expected %d but got %d", prefix, expected, got)
	}
}

func compareString(t *testing.T, prefix string, expected string, got string) {
	if expected != got {
		t.Errorf("%s Expected '%s' but got '%s'", prefix, expected, got)
	}
}
