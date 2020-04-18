package rule

import (
	"github.com/zainul/ark/bridge/cache"
)

type slidingWindowRule struct {
	uniqueID     string
	expiry       int
	limiterType  int
	requestLimit int
	redis        cache.Cache
}

const (
	SlidingWindowLimiterByRPM = 1
	SlidingWindowLimiterByRPS = 2
)
