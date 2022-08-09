package main

import (
	"RoomTgBot/cmd/bot"
	"fmt"
)

func main() {
	fmt.Println("Hello Goland")
	testAllServices()
}

func testAllServices() {
	// cron
	//Test()

	// mongodb
	//DBTest()

	// redis
	//ExampleClient()

	// start echo bot
	bot.Setup()
}
