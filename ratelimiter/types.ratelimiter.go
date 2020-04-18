package ratelimiter

import (
	"time"
)

type rlModule struct {
	rules []Rule
}

type Rule interface {
	IsRateLimit(requestTime time.Time) (bool, error)
}
