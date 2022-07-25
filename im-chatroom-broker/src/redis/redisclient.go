package redis

import (
	"github.com/go-redis/redis/v8"
	"sync"
)

var once sync.Once

var Rdb *redis.Client

func singleton() *redis.Client {
	once.Do(func() {
		Rdb = redis.NewClient(&redis.Options{
			Addr:     "47.95.148.121:6379",
			Password: "o1trUmeh", // no password set
			DB:       1,          // use default DB
		})
	})

	return Rdb
}

const Nil = redis.Nil

func NewZSetMember(score float64, member interface{}) *redis.Z {
	return &redis.Z{
		Score: score, Member: member,
	}
}

func init() {
	singleton()
}
