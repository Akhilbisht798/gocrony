package scheduler

import (
	"time"

	"github.com/robfig/cron/v3"
)

// TODO: save errors and response as logs.
func Scheduler() {
}

func GetNextRun(schedule string) (*time.Time, error) {
	scheduler, err := cron.ParseStandard(schedule)
	now := time.Now()
	if err != nil {
		return &now, err
	}
	nextRun := scheduler.Next(now)
	return &nextRun, nil
}
