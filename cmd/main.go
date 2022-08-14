package main

import (
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

func testAllServices() {
	// cron

	// mongodb
	//DBTest()

	// redis
	ExampleClient()

	// start echo bot
	//bot.Setup()
}
