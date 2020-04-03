package scheduler

import (
	"fmt"
	"time"

	"github.com/robfig/cron"
)

type Scheduler struct{}

func NewScheduler() SchedulerInterface {
	return &Scheduler{}
}

func (s Scheduler) GetJobScheduleFromTime(spec time.Time) string {
	hours, mins, seconds := spec.Clock()
	_, month, day := spec.Date()
	dow := spec.Weekday()

	schedule := fmt.Sprintf(
		"%d %d %d %d %d %d",
		seconds, mins, hours, day, month, dow,
	)

	return schedule
}

func (s Scheduler) SetJob(schedule string, functions []func()) {
	// Log purpose
	fmt.Println(time.Now().Format("2006/01/02 15:04:05"), "Set cronjob for", schedule)

	// Return new cron
	c := cron.New()
	for _, k := range functions {
		c.AddFunc(schedule, func() {
			k()
		})
	}
	c.Start()
}
