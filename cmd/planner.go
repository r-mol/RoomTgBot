package main

import (
	"github.com/robfig/cron/v3"
	"fmt"
)

func Test() {
	c := cron.New(cron.WithSeconds())
	// defer c.Stop()
	lines := 0
	c.AddFunc("@every 1s", func() {
		fmt.Println("line")
		lines++
	})
	c.Start()

	// for lines < 10 {}

}
