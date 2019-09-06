package redis

import (
	"github.com/go-redis/redis"
	"log"
	"os"
)

var client *redis.Client

func GetRedisClient() *redis.Client {
	if client == nil {
		client = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDISURL"),
			Password: os.Getenv("REDISPASSWD"),
			DB:       0, // use default DB
		})
		_, err := client.Ping().Result()
		if err != nil {
			log.Fatalf("Could not connect to redis %v", err)
		}
	}
	return client
}
