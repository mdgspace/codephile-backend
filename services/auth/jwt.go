package auth

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"
	"github.com/mdg-iitr/Codephile/services/redis"
	r "github.com/go-redis/redis"
	"log"
	"time"
)

func GenerateToken(uid string) string {
	currentTimestamp := time.Now().UTC().Unix()
	var ttl = beego.AppConfig.DefaultInt64("TOKENDURATION", 3600)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: currentTimestamp + ttl,
		IssuedAt:  currentTimestamp,
		Issuer:    "mdg",
		Subject:   uid,
	})
	tokenString, err := token.SignedString([]byte(beego.AppConfig.String("HMACKEY")))
	if err != nil {
		log.Fatal(err)
	}
	return tokenString
}
func BlacklistToken(token *jwt.Token) error {
	client := redis.GetRedisClient()
	_, err := client.Set(token.Claims.(jwt.MapClaims)["sub"].(string), true, getTokenRemainingValidity(token.Claims.(jwt.MapClaims)["exp"])).Result()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(token.Claims.(jwt.MapClaims)["sub"].(string))
	return err
}
func IsTokenBlacklisted(token *jwt.Token) bool {
	client := redis.GetRedisClient()
	_, err := client.Get(token.Claims.(jwt.MapClaims)["sub"].(string)).Result()
	if err == r.Nil {
		return false
	}
	return true
}

func getTokenRemainingValidity(timestamp interface{}) time.Duration {
	if validity, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(validity), 0)
		remained := tm.Sub(time.Now())
		if remained > 0 {
			return remained
		}
	}
	return 0;
}
func IsTokenExpired(token *jwt.Token) bool {
	exp := int64(token.Claims.(jwt.MapClaims)["exp"].(float64))
	if exp > time.Now().UTC().Unix() {
		return false
	}
	return true
}
