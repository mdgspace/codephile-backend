package controllers

import (
	"github.com/astaxie/beego"
	"github.com/globalsign/mgo/bson"
	. "github.com/mdg-iitr/Codephile/conf"
	. "github.com/mdg-iitr/Codephile/errors"
	"github.com/mdg-iitr/Codephile/models"
	"log"
	"net/http"
)

type SubmissionController struct {
	beego.Controller
}

// @Title Get
// @Description Get submissions of user(logged-in if uid is empty) across various platforms
// @Security token_auth read:submission
// @Param	uid		path 	string	false		"UID of user"
// @Success 200 {object} types.Submissions
// @Failure 400 invalid uid
// @Failure 404 User/Submission not found
// @router / [get]
// @router /:uid [get]
func (s *SubmissionController) GetSubmission() {
	uidString := s.GetString(":uid")
	var uid bson.ObjectId
	if bson.IsObjectIdHex(uidString) {
		uid = bson.ObjectIdHex(uidString)
	} else if uidString == "" {
		uid = s.Ctx.Input.GetData("uid").(bson.ObjectId)
	} else {
		s.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		s.Data["json"] = BadInputError("Invalid UID")
		s.ServeJSON()
		return
	}
	subs, err := models.GetSubmissions(uid)
	if err != nil {
		s.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		s.Data["json"] = NotFoundError("User/Submission not found")
		s.ServeJSON()
		return
	} else {
		s.Data["json"] = subs
	}
	s.ServeJSON()
}

// @Title Post
// @Description Triggers saving of user's submissions across a particular platform into database
// @Security token_auth write:submission
// @Param	site		path 	string	true		"Platform site name"
// @Success 200 submission successfully saved
// @Failure 400 site invalid
// @Failure 404 user/handle found
// @Failure 500 server_error
// @router /:site [post]
func (s *SubmissionController) SaveSubmission() {
	uid := s.Ctx.Input.GetData("uid").(bson.ObjectId)
	site := s.GetString(":site")
	if !IsSiteValid(site) {
		s.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		s.Data["json"] = BadInputError("Invalid contest site")
		s.ServeJSON()
		return
	}

	err := models.AddSubmissions(uid, site)
	if err == UserNotFoundError || err == HandleNotFoundError {
		s.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		s.Data["json"] = NotFoundError("User/Handle not found")
		s.ServeJSON()
		return
	} else if err != nil {
		log.Println(err.Error())
		s.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		s.Data["json"] = InternalServerError("Internal server error")
		s.ServeJSON()
		return
	}

	s.Data["json"] = map[string]string{"status": "submission successfully saved"}
	s.ServeJSON()
}

// @Title Filter
// @Description Filter submissions of user on the basis of status, site and tags
// @Security token_auth read:submission
// @Param	uid		path 	string	false		"UID of user"
// @Param	site		path 	string	true		"Website name"
// @Param	status		query 	string	false		"Submission status"
// @Param	tag 		query	string	false		"Submission tag"
// @Success 200 {object} submission.CodechefSubmission
// @Failure 400 user not exist
// @Failure 500 server_error
// @router /:site/filter [get]
// @router /:site/:uid/filter [get]
func (s *SubmissionController) FilterSubmission() {
	uidString := s.GetString(":uid")
	var uid bson.ObjectId
	if bson.IsObjectIdHex(uidString) {
		uid = bson.ObjectIdHex(uidString)
	} else if uidString == "" {
		uid = s.Ctx.Input.GetData("uid").(bson.ObjectId)
	} else {
		s.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		s.Data["json"] = BadInputError("Invalid UID")
		s.ServeJSON()
		return
	}
	status := s.GetString("status")
	site := s.GetString(":site")
	tag := s.GetString("tag")
	subs, err := models.FilterSubmission(uid, status, tag, site)
	if err != nil {
		s.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		s.Data["json"] = InternalServerError("Internal server error")
		s.ServeJSON()
		return
	}
	s.Data["json"] = subs
	s.ServeJSON()
}
