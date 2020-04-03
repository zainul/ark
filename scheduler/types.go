package scheduler

import "time"

type SchedulerInterface interface {
	GetJobScheduleFromTime(spec time.Time) string
	SetJob(schedule string, functions []func())
}
