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

func cronExample() {
	// set repeating action
	const sec = time.Second
	id := ScheduleForEvery(1*sec/2, func() { fmt.Println("repeating") })
	// set one time action. Cannot be canceled (for now?)
	RunOnceAfter(sec/2*3, func() { fmt.Println("Once") })
	// cancel first task
	RunOnceAfter(sec*2, func() {
		fmt.Println("cancel task")
		CancelTask(id)
	})

	// because tasks execute in background
	time.Sleep(3 * time.Second)
}

func testAllServices() {
	// cron
	cronExample()

	// mongodb
	//DBTest()

	// redis
	//ExampleClient()

	// start echo bot
	bot.Setup()
}
