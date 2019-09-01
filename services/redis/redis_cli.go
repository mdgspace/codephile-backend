package redis

import (
	"github.com/astaxie/beego"
	"github.com/go-redis/redis"
	"log"
)

var client *redis.Client

func GetRedisClient() *redis.Client {
	if client == nil {
		client = redis.NewClient(&redis.Options{
			Addr:     beego.AppConfig.DefaultString("REDISURL", "localhost:6379"),
			Password: beego.AppConfig.DefaultString("REDISPASSWD", ""),
			DB:       0, // use default DB
		})
		_, err := client.Ping().Result()
		if err != nil {
			log.Fatalf("Could not connect to redis %v", err)
		}
	}
	return client
}
