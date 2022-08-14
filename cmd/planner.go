package main

import (
	"time"

	"github.com/robfig/cron/v3"
)

var Cron = cron.New()

type Task interface {
	Run()
}

type job struct {
	callback func()
}

func (j job) Run() { j.callback() }

func ScheduleForEvery(t time.Duration, cmd func()) {
	Cron.Schedule(cron.Every(t), job{callback: cmd})
}

