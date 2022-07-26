package planner

import (
	"time"

	"github.com/robfig/cron/v3"
)

var Cron = cron.New()

func ScheduleForEvery(t time.Duration, cmd func()) cron.EntryID {
	return Cron.Schedule(cron.Every(t), cron.FuncJob(cmd))
}

func RunOnceAfter(t time.Duration, cmd func()) {
	// Library has no human way to do it
	// "Fine. I'll do it myself." (didn't get => Youtube it)
	go func() {
		time.Sleep(t)
		cmd()
	}()
}

func CancelTask(id cron.EntryID) {
	Cron.Remove(id)
}
