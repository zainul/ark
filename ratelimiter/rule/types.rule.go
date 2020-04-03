package rule

import (
	"github.com/zainul/ark/storage/redis"
)

type slidingWindowRule struct {
	uniqueID      string
	expirySeconds int
	limiterType   int
	requestLimit  int
	redis         redis.Redis
}

const (
	SlidingWindowLimiterByRPM = 1
)
