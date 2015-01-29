package models

import (
	"strconv"
	"time"
)

type Reminder struct {
	Id        int    `db:"id"`
	Reminder  string `db:"reminder"`
	Timestamp int32  `db:"timestamp"`
	Userid    int32  `db:"userid"`
}

type ReadableReminder struct {
	Reminder string
	Date     string
	Id       int
}

func (r *Reminder) Save() error {
	return Db.Insert(r)
}

func (r *Reminder) FindByTimestamp(timestamp int32) error {
	err := Db.SelectOne(r, "SELECT * FROM reminders WHERE timestamp = :ts", map[string]interface{}{"ts": timestamp})
	if err != nil {
		return err
	}
	return nil
}

func (r *Reminder) Remove() error {
	_, err := Db.Delete(r)
	if err != nil {
		return err
	}
	return nil
}

func RemoveById(id int) error {
	query := "DELETE FROM reminders WHERE id= '" + strconv.Itoa(id) + "'"
	_, err := Db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func GetReminders(userid int32) []Reminder {
	var reminders []Reminder
	query := "select * from reminders where userid = '" + strconv.Itoa(int(userid)) + "'"
	_, err := Db.Select(&reminders, query)
	if err != nil {
		panic(err)
	}
	return reminders
}

func GetAllReminders() []Reminder {
	var reminders []Reminder
	_, err := Db.Select(&reminders, "select * from reminders")
	if err != nil {
		panic(err)
	}
	return reminders
}

func MakeReadableReminders(reminders []Reminder) []ReadableReminder {
	readableReminders := []ReadableReminder{}
	for _, reminder := range reminders {
		year, month, day := time.Unix(int64(reminder.Timestamp), 0).Date()
		readableDate := strconv.Itoa(day) + " " + month.String() + ", " + strconv.Itoa(year)
		r := ReadableReminder{
			Reminder: reminder.Reminder,
			Date:     readableDate,
			Id:       reminder.Id,
		}
		readableReminders = append(readableReminders, r)
	}
	return readableReminders
}
