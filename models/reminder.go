package models

type Reminder struct {
	Id        int    `db:"id"`
	Reminder  string `db:"reminder"`
	Timestamp int32  `db:"timestamp"`
}

func (r *Reminder) Save() error {
	return Db.Insert(r)
}

func (r *Reminder) Remove() error {
	_, err := Db.Delete(r)
	if err != nil {
		return err
	}
	return nil
}

func GetAllReminders() []Reminder {
	var reminders []Reminder
	_, err := Db.Select(&reminders, "select * from reminders")
	if err != nil {
		panic(err)
	}
	return reminders
}
