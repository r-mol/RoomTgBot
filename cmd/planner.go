package main

import (
	"time"

	"github.com/robfig/cron/v3"
)

var Cron = cron.New()

func ScheduleForEvery(t time.Duration, cmd func()) {
	Cron.Schedule(cron.Every(t), cron.FuncJob(cmd))
}

