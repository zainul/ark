package cache

import (
	"time"
)

type DriverImplementor string

const (
	GoRedis DriverImplementor = "go-redis"
)

type (
	Config struct {
		Host     string
		Port     int
		Password string
	}
)

var (
	driversImplementors map[DriverImplementor]func(cfg Config) Cache
)

func init() {
	ds := make(map[DriverImplementor]func(cfg Config) Cache)
	ds[GoRedis] = newGoRedis
	driversImplementors = ds
}

func New(cfg Config, driverPlugin DriverImplementor) Cache {
	return driversImplementors[driverPlugin](cfg)
}

type Cache interface {
	SetNX(key string, value interface{}, ttl time.Duration) error
	Get(key string, targeted interface{}) error
}
