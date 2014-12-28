package main

import (
	"fmt"

	"github.com/fabioberger/recall/config"
	"github.com/fabioberger/recall/controllers"
	"github.com/fabioberger/recall/models"
)

func main() {
	config.Init()
	models.Init()
	reminders := controllers.Reminders{
		Mock: false,
	}
	log_messages := reminders.CheckAll()
	fmt.Println(log_messages)
}
