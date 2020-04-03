package rule_test

import (
	"testing"
	"time"

	"github.com/bouk/monkey"
	"github.com/stretchr/testify/assert"
	. "github.com/zainul/ark/ratelimiter/rule"
	"github.com/zainul/ark/storage/redis/dummyrds"
)

func TestSlidingWindow(t *testing.T) {
	m := dummyrds.Mocker{}
	m.AddMock("GET travel:rl:slidingwindow:apisearchtrain:1573725600", "100", false)
	m.AddMock("GET travel:rl:slidingwindow:apisearchtrain:1576317600", "5", false)

	rds := dummyrds.New(dummyrds.Config{
		MockingMap: m,
	})

	pload, _ := time.LoadLocation("Europe/London")
	wayback := time.Date(2010, time.May, 19, 16, 05, 0, 4, pload)
	patch := monkey.Patch(time.Now, func() time.Time { return wayback })
	defer patch.Unpatch()

	r1 := NewSlidingWindow("apisearchtrain", 120, 1, 10, rds)
	r2 := NewSlidingWindow("apisearchtrain", 120, 1, 100, rds)

	// Exceed maximum hits
	isBot, err := r1.IsRateLimit(time.Date(2019, 11, 14, 10, 0, 0, 0, time.UTC))
	assert.True(t, isBot)
	assert.Nil(t, err)

	// Not maximum hits
	isBot, err = r2.IsRateLimit(time.Date(2019, 12, 14, 10, 0, 0, 0, time.UTC))
	assert.False(t, isBot)
	assert.Nil(t, err)
}
