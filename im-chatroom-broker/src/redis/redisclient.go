package redis

import (
	"github.com/go-redis/redis/v8"
	"sync"
)

var once sync.Once

var client *redis.Client

func Singleton() *redis.Client {
	once.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:     "47.95.148.121:6379",
			Password: "o1trUmeh", // no password set
			DB:       1,          // use default DB
		})
	})

	return client
}

