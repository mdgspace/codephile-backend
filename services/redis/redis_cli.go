package redis

import (
	"github.com/go-redis/redis"
	"log"
	"os"
)

var client *redis.Client

func GetRedisClient() *redis.Client {
	if client == nil {
		opt, err := redis.ParseURL(os.Getenv("REDISURL"))
		if err != nil {
			panic(err.Error())
		}
		client = redis.NewClient(opt)
		_, err = client.Ping().Result()
		if err != nil {
			log.Fatalf("Could not connect to redis %v", err)
		}
	}
	return client
}
