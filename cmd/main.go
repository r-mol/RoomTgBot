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
	bot.Setup()
}
