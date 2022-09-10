package main

import (
	"RoomTgBot/cmd/bot"

	"fmt"
	// "time"
)

func init() {
	Cron.Start()
}

func main() {
	fmt.Println("Hello Goland")
	bot.Setup()
}
