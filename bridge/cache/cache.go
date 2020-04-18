package cache

import (
	"time"
)

//DriverImplementor is implementor of driver, can go-redis or redigo
type DriverImplementor string

const (
	GoRedis DriverImplementor = "go-redis"
)

type (
	//Config configuration of cache
	Config struct {
		Host     string
		Port     int
		Password string
	}
)

var (
	driversImplementors map[DriverImplementor]func(cfg Config, options ...Option) Cache
)

func init() {
	ds := make(map[DriverImplementor]func(cfg Config, options ...Option) Cache)
	ds[GoRedis] = newGoRedis
	driversImplementors = ds
}

type cacher struct {
	Host        string
	Port        int
	Password    string
	ReTryOption ReTryOption
}

//Option is type of any option that passing to cache
type Option func(c *cacher)

//ReTryOption structure for retry mechanism
type ReTryOption struct {
	// Maximum number of retries before giving up.
	// Default is to not retry failed commands.
	MaxRetries int
	// Minimum backoff between each retry.
	// Default is 8 milliseconds; -1 disables backoff.
	MinRetryBackoff time.Duration
	// Maximum backoff between each retry.
	// Default is 512 milliseconds; -1 disables backoff.
	MaxRetryBackoff time.Duration
}

//WithRetry is function to handle retry mechanism
func WithRetry(retry ReTryOption) Option {
	return func(c *cacher) {
		if retry.MaxRetries > 0 {
			c.ReTryOption.MaxRetries = retry.MaxRetries
		}

		if retry.MinRetryBackoff != 0 {
			c.ReTryOption.MinRetryBackoff = retry.MinRetryBackoff
		}

		if retry.MaxRetryBackoff != 0 {
			c.ReTryOption.MaxRetryBackoff = retry.MaxRetryBackoff
		}
	}
}

//New is initiate of instance cache
func New(driverPlugin DriverImplementor, cfg Config, options ...Option) Cache {
	return driversImplementors[driverPlugin](cfg, options...)
}

//Cache abstraction of caching
type Cache interface {
	//SetNX is set key value in cache with expired time duration
	SetNX(key string, value interface{}, ttl time.Duration) error
	//Get is getting value base on passing key the result in targeted
	Get(key string, targeted interface{}) error
	//TxIncr is increment using transaction pipeline
	TxIncr(keys string) (int64, error)
	//Expire is set the keys will expire in some second
	Expire(keys string, ttl time.Duration) error
}
