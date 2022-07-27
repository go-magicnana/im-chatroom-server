package redis

import (
	"github.com/go-redis/redis/v8"
	"im-chatroom-gateway/config"
	"sync"
)

var once sync.Once

var Rdb *redis.Client

const Nil = redis.Nil

func singleton() *redis.Client {
	once.Do(func() {
		Rdb = redis.NewClient(&redis.Options{
			Addr:     config.OP.Redis.Address,
			Password: config.OP.Redis.Password, // no password set
			DB:       config.OP.Redis.Db,       // use default DB
		})
	})

	return Rdb
}

func NewZSetMember(score float64, member interface{}) *redis.Z {
	return &redis.Z{
		Score: score, Member: member,
	}
}

func init() {
	singleton()
}
