package auth

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"
	"github.com/getsentry/sentry-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	r "github.com/go-redis/redis"
	"github.com/mdg-iitr/Codephile/services/redis"
	"log"
	"os"
	"strconv"
	"time"
)

// identifier used to prevent user from logging in again
// to be used if a user is suspicious
var UserBlacklisted = "blacklisted"

func GenerateToken(uid string) string {
	currentTimestamp := time.Now().UTC().Unix()
	var ttl = beego.AppConfig.DefaultInt64("TOKENDURATION", 3600000)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: currentTimestamp + ttl,
		IssuedAt:  currentTimestamp,
		Issuer:    "mdg",
		Subject:   uid,
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("HMACKEY")))
	if err != nil {
		sentry.CaptureException(err)
		log.Fatal(err)
	}
	return tokenString
}
func BlacklistToken(token *jwt.Token) error {
	client := redis.GetRedisClient()
	claims := token.Claims.(jwt.MapClaims)
	_, err := client.Set(claims["sub"].(string), int64(claims["iat"].(float64)), getTokenRemainingValidity(token.Claims.(jwt.MapClaims)["exp"])).Result()
	if err != nil {
		sentry.CaptureException(err)
		log.Println(err.Error())
	}
	return err
}
func IsTokenBlacklisted(token *jwt.Token) bool {
	client := redis.GetRedisClient()
	claims := token.Claims.(jwt.MapClaims)
	val, err := client.Get(claims["sub"].(string)).Result()
	if err == r.Nil {
		return false
	} else if err != nil {
		return true
	}
	iat, _ := strconv.ParseInt(val, 10, 64)
	if int64(claims["iat"].(float64)) == iat || val == UserBlacklisted {
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

//TODO: check if uid.String() is working perfectly or not, earlier string(uid), uid bson.ObjectID

func BlacklistUser(uid primitive.ObjectID) error {
	client := redis.GetRedisClient()
	_, err := client.Set(uid.String(), UserBlacklisted, 0).Result()
	return err
}

func WhitelistUser(uid primitive.ObjectID) error {
	client := redis.GetRedisClient()
	val := client.Get(uid.String()).Val()
	if val != UserBlacklisted {
		return errors.New("already whitelisted")
	}
	_, err := client.Del(uid.String()).Result()
	return err
}
