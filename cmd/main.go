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

// store values in days for each notification type
type Config struct {
	selfCleaningNotification int64
	// cleaningNotificationList []int64
	removeAbsent      int64
	moneyNotification int64
}

func day2duration(x int64) time.Duration {
	return time.Duration(x) * 24 * time.Hour
}

func main() {
	fmt.Println("(not) Read Config")
	// TODO add YAML support

	config := Config{
		selfCleaningNotification: 7,
		removeAbsent:             3,
		moneyNotification:        30,
	}

	fmt.Println("Setup schedule")
	planner.ScheduleForEvery(
		day2duration(config.selfCleaningNotification),
		func() { bot.FindInitCleanMan() },
	)
	planner.ScheduleForEvery(
		day2duration(config.moneyNotification),
		func() { bot.NotifyAboutMoney() },
	)
	planner.ScheduleForEvery(
		day2duration(config.removeAbsent),
		func() { bot.PutNotAbsentToAllUsers() },
	)
	// for _, day := range config.cleaningNotificationList{
	// }

	fmt.Println("Start bot")
	bot.Setup()
}
