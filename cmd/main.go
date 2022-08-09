package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello Goland")
	testAllServices()
}

func testAllServices() {
	// cron
	Test()

	// mongodb
	DBTest()

	// redis
	ExampleClient()

	// start echo bot
	Setup()
}
