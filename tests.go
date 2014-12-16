package main

import (
	"testing"

	"github.com/fabioberger/recall/models"
)

func TestReminderSave(t *testing.T) {
	reminder := &models.Reminder{
		Reminder:  "Revise the Dominic Number System",
		Timestamp: 1418727704,
	}
	reminder.Save()
	retrieved_reminder := &models.Reminder{}
	retrieved_reminder.FindByTimestamp(1418727704)
	compareString(t, "ReminderSave", retrieved_reminder.Reminder, reminder.Reminder)
	compareInt(t, "ReminderSave", int64(retrieved_reminder.Timestamp), 1418727704)
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
