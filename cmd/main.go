package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hello Goland")
	Test()
	time.Sleep(3*time.Second)
	DBTest()
	Setup()
}
