package rule

import (
	"errors"
	"strconv"
	"time"

	"github.com/zainul/ark/bridge/cache"
	"github.com/zainul/ark/ratelimiter"
)

func NewSlidingWindow(uniqueID string, expiry int, limiterType int, requestLimit int, rds cache.Cache) ratelimiter.Rule {
	return &slidingWindowRule{
		uniqueID:     uniqueID,
		limiterType:  limiterType,
		requestLimit: requestLimit,
		expiry:       expiry,
		redis:        rds,
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
	} else if s.limiterType == SlidingWindowLimiterByRPS {
		d = 1 * time.Second
	}

	redisKey := "ark:rl:slidingwindow:" + s.uniqueID + ":" + strconv.FormatInt(requestTime.Truncate(d).Unix(), 10)

	// Get current hits value
	var currentHits int
	err := s.redis.Get(redisKey, &currentHits)

	if err != nil {
		return false, err
	}

	// Exceed maximum hits
	if currentHits >= s.requestLimit {
		return true, nil
	}

	// Increment total hits
	if count, _ := s.redis.TxIncr(redisKey); int(count) >= s.requestLimit {
		_ = s.redis.Expire(redisKey, time.Duration(s.expiry)*time.Second)
		return true, nil
	}

	// Set redis expiry
	_ = s.redis.Expire(redisKey, time.Duration(s.expiry)*time.Second)

	return false, nil
}
