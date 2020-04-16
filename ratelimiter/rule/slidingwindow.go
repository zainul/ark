package rule

import (
	"errors"
	"strconv"
	"time"

	"github.com/zainul/ark/ratelimiter"
	"github.com/zainul/ark/storage/redis"
)

func NewSlidingWindow(uniqueID string, expirySeconds int, limiterType int, requestLimit int, rds redis.Redis) ratelimiter.RateLimiterRule {
	return &slidingWindowRule{
		uniqueID:      uniqueID,
		limiterType:   limiterType,
		requestLimit:  requestLimit,
		expirySeconds: expirySeconds,
		redis:         rds,
	}
}

// IsRateLimit checking rateLimit
func (s *slidingWindowRule) IsRateLimit(requestTime time.Time) (bool, error) {

	// Failed to get value
	if s.uniqueID == "" {
		return false, errors.New("ID is missing")
	}

	d := 60 * time.Second

	if s.limiterType == SlidingWindowLimiterByRPM {
		d = 60 * time.Second
	}

	redisKey := "ark:rl:slidingwindow:" + s.uniqueID + ":" + strconv.FormatInt(requestTime.Truncate(d).Unix(), 10)

	// Get current hits value
	currentHits := s.redis.Get(redisKey).Int()

	// Exceed maximum hits
	if currentHits >= s.requestLimit {
		return true, nil
	}

	// Increment total hits
	if count, _ := s.redis.IncrSingle(redisKey); count >= s.requestLimit {
		s.redis.Expire(redisKey, s.expirySeconds)
		return true, nil
	}

	// Set redis expiry
	s.redis.Expire(redisKey, s.expirySeconds)

	return false, nil
}
