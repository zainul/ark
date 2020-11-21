package cache

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/vmihailenco/msgpack"
	"log"
	"time"
)

type goRedis struct {
	client *redis.Client
}

func newGoRedis(cfg Config, options ...Option) Cache {

	c := cacher{
		Host:     cfg.Host,
		Port:     cfg.Port,
		Password: cfg.Password,
	}

	for _, opt := range options {
		opt(&c)
	}

	client := redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%d", c.Host, c.Port),
		Password:        c.Password, // no password set
		DB:              0,          // use default DB
		MaxRetries:      c.ReTryOption.MaxRetries,
		MaxRetryBackoff: c.ReTryOption.MaxRetryBackoff,
		MinRetryBackoff: c.ReTryOption.MinRetryBackoff,
	})

	_, err := client.Ping().Result()

	if err != nil {
		log.Print("Failed to connect to redis")
	}

	return &goRedis{client: client}
}

func (g *goRedis) SetNX(key string, value interface{}, ttl time.Duration) error {
	bt, err := msgpack.Marshal(value)

	if err != nil {
		return err
	}

	return g.client.SetNX(key, bt, ttl).Err()
}

func (g *goRedis) Get(key string, targeted interface{}) error {
	bt, err := g.client.Get(key).Bytes()

	if err != nil {
		return err
	}

	err = msgpack.Unmarshal(bt, targeted)

	if err != nil {
		return err
	}

	return nil
}

func (g *goRedis) TxIncr(keys string) (int64, error) {
	pipe := g.client.TxPipeline()
	incr := pipe.Incr(keys)
	_, err := pipe.Exec()

	if incr == nil {
		return 0, errors.New("increment failed")
	}

	return incr.Val(), err
}

func (g *goRedis) Expire(key string, ttl time.Duration) error {
	return g.client.Expire(key, ttl).Err()
}
