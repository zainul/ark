package ratelimiter

import (
	"time"

	"github.com/zainul/ark/logging/datadog"
)

type rlModule struct {
	rules []RateLimiterRule
	ddog  datadog.Datadog
}

type Ratelimiter interface {
	IsRateLimit(requestTime time.Time) (bool, error)
}

type RateLimiterRule interface {
	IsRateLimit(requestTime time.Time) (bool, error)
}
