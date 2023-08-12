package controllers

import (
	// "context"
	// "encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"

	// "os"
	// "time"

	// "github.com/google/uuid"
	// "github.com/mdg-iitr/Codephile/services/mail"

	"github.com/astaxie/beego"
	// "github.com/dgrijalva/jwt-go"
	// "github.com/dgrijalva/jwt-go/request"
	"github.com/getsentry/sentry-go"
	// "github.com/globalsign/mgo/bson"
	// "github.com/gorilla/schema"
	// . "github.com/mdg-iitr/Codephile/conf"
	. "github.com/mdg-iitr/Codephile/errors"
	"github.com/mdg-iitr/Codephile/models"
	"github.com/mdg-iitr/Codephile/models/types"
	// "github.com/mdg-iitr/Codephile/scrappers"
	// "github.com/mdg-iitr/Codephile/services/auth"
	// "github.com/mdg-iitr/Codephile/services/firebase"
	// "github.com/mdg-iitr/Codephile/services/redis"
	// "github.com/mdg-iitr/Codephile/services/worker"
)



// Operations about the User's CP institute level rank
type RankController struct {
	beego.Controller
}

// @Title Codechef rank
// @Description Codechef inistitute level rank based on codechef worldrank
// @Security token_auth read:user
// @Param	institute 			query	string	true "institute"
// @Success 200 {object} []types.SearchDoc
// @Failure 409 no user to belongs to this institute or user doesnt have codechef handle
// @Failure 400 bad request no institute parameter
// @Failure 500 server_error
// @router /codechef[get]
func (u *UserController) codechefRank() {
	instituteName := u.GetString("institute")
	res, err := models.FilterUsers(instituteName)
	if err != nil {
		hub := sentry.GetHubFromContext(u.Ctx.Request.Context())
		hub.CaptureException(err)
		log.Println(err.Error())
		u.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		u.Data["json"] = InternalServerError("server error.. report to admin")
		u.ServeJSON()
		return
	}
	if len(res) == 0 {
		u.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		u.Data["json"] = NotFoundError("no user belongs to this institute")
		u.ServeJSON()
		return
	}
	userWithCodechefHandle:= []types.SearchDoc{}
	
	for _,user := range res {
		if(user.Handle.Codechef!=""){
            userWithCodechefHandle= append(userWithCodechefHandle,user)
		}
	}
    if len(userWithCodechefHandle) == 0 {
		u.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		u.Data["json"] = NotFoundError("no user have codechef handle")
		u.ServeJSON()
		return
	}
    sort.Slice(userWithCodechefHandle,func(i,j int)bool{
        p1,_ := models.GetProfiles(userWithCodechefHandle[i].ID)
		p2,_ := models.GetProfiles(userWithCodechefHandle[i].ID)
		rank1,_ := strconv.Atoi(p1.CodechefProfile.WorldRank) 
		rank2,_ := strconv.Atoi(p2.CodechefProfile.WorldRank)
		return rank1>rank2
	})
	
	u.Data["json"] = userWithCodechefHandle
	u.ServeJSON()	
}

// @Title Codeforces rank
// @Description Codechef inistitute level rank based on codeforces worldrank
// @Security token_auth read:user
// @Param	institute 			query	string	true "institute"
// @Success 200 {object} []types.SearchDoc
// @Failure 409 no user to belongs to this institute or user doesnt have codechef handle
// @Failure 400 bad request no institute parameter
// @Failure 500 server_error
// @router /codeforces[get]
func (u *UserController) codeforcesRank() {
	instituteName := u.GetString("institute")
	res, err := models.FilterUsers(instituteName)
	if err != nil {
		hub := sentry.GetHubFromContext(u.Ctx.Request.Context())
		hub.CaptureException(err)
		log.Println(err.Error())
		u.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		u.Data["json"] = InternalServerError("server error.. report to admin")
		u.ServeJSON()
		return
	}
	if len(res) == 0 {
		u.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		u.Data["json"] = NotFoundError("no user belongs to this institute")
		u.ServeJSON()
		return
	}
	userWithCodeforcesHandle:= []types.SearchDoc{}
	
	for _,user := range res {
		if(user.Handle.Codechef!=""){
            userWithCodeforcesHandle= append(userWithCodeforcesHandle,user)
		}
	}
    if len(userWithCodeforcesHandle) == 0 {
		u.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		u.Data["json"] = NotFoundError("no user have codeforces handle")
		u.ServeJSON()
		return
	}
    sort.Slice(userWithCodeforcesHandle,func(i,j int)bool{
        p1,_ := models.GetProfiles(userWithCodeforcesHandle[i].ID)
		p2,_ := models.GetProfiles(userWithCodeforcesHandle[i].ID)
		rank1,_ := strconv.Atoi(p1.CodeforcesProfile.WorldRank) 
		rank2,_ := strconv.Atoi(p2.CodeforcesProfile.WorldRank)
		return rank1>rank2
	})
	
	u.Data["json"] = userWithCodeforcesHandle
	u.ServeJSON()	
}