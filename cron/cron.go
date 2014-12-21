package main

import (
	"github.com/fabioberger/recall/config"
	"github.com/fabioberger/recall/controllers"
	"github.com/fabioberger/recall/models"
)

func main() {
	config.Init()
	models.Init()
	reminders := controllers.Reminders{}
	reminders.Check()
}
