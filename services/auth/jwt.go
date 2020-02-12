package auth

import (
	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"
	r "github.com/go-redis/redis"
	"github.com/mdg-iitr/Codephile/services/redis"
	"log"
	"os"
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
	tokenString, err := token.SignedString([]byte(os.Getenv("HMACKEY")))
	if err != nil {
		log.Fatal(err)
	}
	return tokenString
}
func BlacklistToken(token *jwt.Token) error {
	client := redis.GetRedisClient()
	claims := token.Claims.(jwt.MapClaims)
	_, err := client.Set(claims["sub"].(string), int64(claims["iat"].(float64)), getTokenRemainingValidity(token.Claims.(jwt.MapClaims)["exp"])).Result()
	if err != nil {
		log.Println(err.Error())
	}
	return err
}
func IsTokenBlacklisted(token *jwt.Token) bool {
	client := redis.GetRedisClient()
	claims := token.Claims.(jwt.MapClaims)
	iat, err := client.Get(claims["sub"].(string)).Int64()
	if err == r.Nil {
		return false
	}
	if int64(claims["iat"].(float64)) == iat {
		return true
	}
	return false
}

func getTokenRemainingValidity(timestamp interface{}) time.Duration {
	if validity, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(validity), 0)
		remained := time.Until(tm)
		if remained > 0 {
			return remained
		}
	}
	return 0;
}
func IsTokenExpired(token *jwt.Token) bool {
	exp := int64(token.Claims.(jwt.MapClaims)["exp"].(float64))
	return exp <= time.Now().UTC().Unix()
}
