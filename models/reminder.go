package models

type Reminder struct {
	Id        int    `db:"id"`
	Reminder  string `db:"reminder"`
	Timestamp int32  `db:"timestamp"`
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

func GetAllReminders() []Reminder {
	var reminders []Reminder
	_, err := Db.Select(&reminders, "select * from reminders")
	if err != nil {
		panic(err)
	}
	return reminders
}
