package threshold

import (
	"github.com/zainul/ark/storage/redis"
)

// renew
type threshold struct {
	timeToLive  int
	max         int
	rds         redis.Redis
	renewExpire bool
}

type ThresholdParam struct {
	TimeToLive int
	Max        int
	Redis      redis.Redis
	// sample renewAttemptExpired => true
	// if expire in 60 mins
	// then at 59 mins to expire have attempt again to that block code (Attempt)
	// so will be renewthe attempt Expired to be 1 + 60 mins
	// More Detail:
	// if you making scheduler to check something set RenewExpire as false
	//   - this will make your key holds expire until its expire it self
	// if you making spam prevention set RenewExpire as true
	//	 - this will make your key keep renewing expire each time Attamp is called
	RenewExpire bool
}

func NewThreshold(t ThresholdParam) Threshold {
	return &threshold{
		timeToLive:  t.TimeToLive,
		max:         t.Max,
		rds:         t.Redis,
		renewExpire: t.RenewExpire,
	}
}

func (t *threshold) Attempt(key string) error {
	val, err := t.rds.IncrSingle(key)
	if err != nil {
		return err
	}

	// conditional 1 - val == t.max
	// 		if reach threshold set expire
	// conditional 2 - val > t.max && t.renewExpire
	// 		renewing expire or not decision
	if val == t.max || (val > t.max && t.renewExpire) {
		t.rds.Expire(key, t.timeToLive)
	}

	return nil
}

func (t *threshold) doFailOver(action ...func()) {
	for _, f := range action {
		f()
	}
}

func (t *threshold) IsAllow(key string, failAction ...func()) bool {
	res := t.rds.Get(key)
	if res.Error != nil {
		t.doFailOver(failAction...)
		return false
	}

	if res.Int64() <= int64(t.max) {
		return true
	} else {
		t.doFailOver(failAction...)
		return false
	}
}
