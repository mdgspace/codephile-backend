package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/getsentry/sentry-go"
	. "github.com/mdg-iitr/Codephile/errors"
	"github.com/mdg-iitr/Codephile/models"
)

// @Title Search
// @Description Endpoint to search users
// @Security token_auth read:user
// @Param	count		query 	string	false		"No of search objects to be returned"
// @Param	query		query 	string	true		"Search query"
// @Param	path		query	string	false		"Field to search in"
// @Success 200 {object} []types.SearchDoc
// @Failure 400 "search query too small"
// @Failure 500 server_error
// @router /search [get]
func (u *UserController) Search() {
	query := u.GetString("query")
	if len(query) < 4 {
		u.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		u.Data["json"] = BadInputError("Search query too small")
		u.ServeJSON()
		return
	}
	count := u.GetString("count")
	c, err := strconv.Atoi(count)
	//Default query response size
	if err != nil {
		c = 500
	}
	path := u.GetString("path")
	//Checking for possible fields, default is "username"
	var possible_path []string
	var good_path bool = false
	possible_path = []string{"username", "fullname", "handle.codechef", "handle.codeforces", "handle.hackerearth", "handle.hackerrank", "handle.spoj"}
	for _, itr := range possible_path {
		if itr == path {
			good_path = true
		}
	}
	if good_path == false {
		path = "username"
	}
	results, err := models.SearchUser(query, c, path)
	if err != nil {
		hub := sentry.GetHubFromContext(u.Ctx.Request.Context())
		hub.CaptureException(err)
		log.Println(err.Error())
		u.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		u.Data["json"] = InternalServerError("Internal server error")
		u.ServeJSON()
		return
	}
	u.Data["json"] = results
	u.ServeJSON()
}
