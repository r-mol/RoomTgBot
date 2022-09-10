package main

import (
	"RoomTgBot/cmd/bot"
	"RoomTgBot/internal/planner"

	"fmt"
	"time"
)

func init() {
	planner.Cron.Start()
}

type Config struct {
	selfCleaningNotification time.Duration
	cleaningNotificationList []time.Duration
	moneyNotification        time.Duration
}

var config Config = Config{
}

func main() {
	fmt.Println("Read Config")

	fmt.Println("Setup schedule")

	fmt.Println("Start bot")
	bot.Setup()
}
