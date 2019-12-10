package middleware

import (
	"errors"
	"github.com/astaxie/beego/context"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/globalsign/mgo/bson"
	"github.com/mdg-iitr/Codephile/services/auth"
	"os"
	"strings"
)

//Checks if token is valid and put valid token in context
func Authenticate(ctx *context.Context) {
	// signup and login endpoints
	if (strings.HasPrefix(ctx.Request.RequestURI, "/v1/user/login") && ctx.Request.Method == "POST") ||
		(strings.HasPrefix(ctx.Request.RequestURI, "/v1/user/signup") && ctx.Request.Method == "POST") ||
		(strings.HasPrefix(ctx.Request.RequestURI, "/v1/user/available") && ctx.Request.Method == "GET") ||
		(strings.HasPrefix(ctx.Request.RequestURI, "/v1/user/verify/" && ctx.Request.Method == "GET"){
		return
	}
	requestToken, err := request.ParseFromRequest(ctx.Request, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unauthorized")
		}
		return []byte(os.Getenv("HMACKEY")), nil
	})
	if err != nil {
		ctx.ResponseWriter.WriteHeader(401)
		ctx.ResponseWriter.Write([]byte("401 Unauthorized\n"))
		return
	}
	if requestToken.Valid && !auth.IsTokenExpired(requestToken) && !auth.IsTokenBlacklisted(requestToken) {
		claim := requestToken.Claims.(jwt.MapClaims)
		uid := bson.ObjectIdHex(claim["sub"].(string))
		ctx.Input.SetData("uid", uid)
	} else {
		ctx.ResponseWriter.WriteHeader(401)
		ctx.ResponseWriter.Write([]byte("401 Unauthorized\n"))
	}
}
