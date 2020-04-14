package cache

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/vmihailenco/msgpack"
	"log"
	"time"
)

type goRedis struct {
	client *redis.Client
}

func newGoRedis(cfg Config) Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password, // no password set
		DB:       0,            // use default DB
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
