package redis

import (
	"github.com/getsentry/sentry-go"
	"github.com/go-redis/redis"
	"log"
	"os"
)

var client *redis.Client

func init() {
	initRedisClient()
}

func initRedisClient() {
	opt, err := redis.ParseURL(os.Getenv("REDISURL"))
	if err != nil {
		panic(err.Error())
	}
	client = redis.NewClient(opt)
	_, err = client.Ping().Result()
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Could not connect to redis %v", err)
	}
}

func GetRedisClient() *redis.Client {
	if client == nil {
		initRedisClient()
	}
	return client
}
