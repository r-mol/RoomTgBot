package main

import (
	"RoomTgBot/cmd/bot"

	"fmt"
	"time"
)

func init() {
	Cron.Start()
}

func main() {
	fmt.Println("Hello Goland")
	testAllServices()
}

const (
	two   = 2
	three = 3
)

func cronExample() {
	// set repeating action
	const sec = time.Second
	id := ScheduleForEvery(sec/two, func() { fmt.Println("repeating") })
	// set one time action. Cannot be canceled (for now?)
	RunOnceAfter(sec/two*three, func() { fmt.Println("Once") })
	// cancel first task
	RunOnceAfter(sec*two, func() {
		fmt.Println("cancel task")
		CancelTask(id)
	})

	// because tasks execute in background
	time.Sleep(three * time.Second)
}

func testAllServices() {
	// cron
	cronExample()

	// mongodb
	DBTest()

	// redis
	ExampleClient()

	// start echo bot
	bot.Setup()
}
