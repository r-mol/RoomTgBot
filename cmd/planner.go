package main

import (
	"fmt"

	cron "github.com/robfig/cron/v3"
)

func Test() {
	c := cron.New(cron.WithSeconds())

	lines := 0

	_, err := c.AddFunc("@every 1s", func() {
		fmt.Println("line")
		lines++
	})

	if err != nil {
		return
	}

	c.Start()
}
