package ratelimiter_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	. "github.com/zainul/ark/ratelimiter"
	"github.com/zainul/ark/ratelimiter/rule"
	"github.com/zainul/ark/storage/redis/dummyrds"
)

func TestIsRateLimiterFalse(t *testing.T) {

	rules := []RateLimiterRule{}

	m := dummyrds.Mocker{}

	rds := dummyrds.New(dummyrds.Config{
		MockingMap: m,
	})

	rules = append(rules, rule.NewSlidingWindow("apisearchtrain", 60, rule.SlidingWindowLimiterByRPM, 1, rds))

	rateLimiter := New(rules)

	isRateLimit, err := rateLimiter.IsRateLimit(time.Now())
	assert.False(t, isRateLimit)
	assert.Nil(t, err)

}

func TestIsRateLimiterTrue(t *testing.T) {

	rules := []RateLimiterRule{}

	m := dummyrds.Mocker{}

	rds := dummyrds.New(dummyrds.Config{
		MockingMap: m,
	})

	rules = append(rules, rule.NewSlidingWindow("apisearchtrain", 60, rule.SlidingWindowLimiterByRPM, 0, rds))

	rateLimiter := New(rules)

	isRateLimit, err := rateLimiter.IsRateLimit(time.Now())
	assert.True(t, isRateLimit)
	assert.Nil(t, err)

}

func TestIsRateLimiterInvalidID(t *testing.T) {

	rules := []RateLimiterRule{}

	m := dummyrds.Mocker{}

	rds := dummyrds.New(dummyrds.Config{
		MockingMap: m,
	})

	rules = append(rules, rule.NewSlidingWindow("", 60, rule.SlidingWindowLimiterByRPM, 0, rds))

	rateLimiter := New(rules)

	isRateLimit, err := rateLimiter.IsRateLimit(time.Now())
	assert.False(t, isRateLimit)
	assert.NotNil(t, err)

}
