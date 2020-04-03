package dummy

import (
	"fmt"
	"time"

	"github.com/zainul/ark/scheduler"
)

type Scheduler struct{}

func NewScheduler() scheduler.SchedulerInterface {
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
	return
}
