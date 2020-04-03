package scheduler_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zainul/ark/scheduler"
)

func TestSetCronSchedule(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		time     time.Time
	}{
		{
			"Success",
			"0 6 17 7 5 1",
			time.Date(2018, 5, 7, 17, 6, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sch := scheduler.NewScheduler()
			actual := sch.GetJobScheduleFromTime(tt.time)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestSetCronJob(t *testing.T) {
	tests := []struct {
		name     string
		schedule string
		fn       []func()
	}{
		{
			"Success start time cron",
			"0 6 17 7 5 1",
			[]func(){
				func() {},
			},
		},
		{
			"Success end time cron",
			"0 6 17 7 5 1",
			[]func(){
				func() {},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sch := scheduler.NewScheduler()
			sch.SetJob(tt.schedule, tt.fn)
		})
	}
}
