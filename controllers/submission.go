package controllers

import (
	"github.com/astaxie/beego"
	"github.com/globalsign/mgo/bson"
	. "github.com/mdg-iitr/Codephile/conf"
	"github.com/mdg-iitr/Codephile/models"
)

var valid_sites = []string{HACKERRANK, CODECHEF, CODEFORCES, SPOJ}

type SubmissionController struct {
	beego.Controller
}

// @Title Get
// @Description Get submissions of user(logged-in if uid is empty) across various platforms
// @Security token_auth read:submission
// @Param	uid		path 	string	false		"UID of user"
// @Success 200 {object} types.Submissions
// @Failure 403 user not exist
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
		s.Ctx.ResponseWriter.WriteHeader(403)
		s.ServeJSON()
		return
	}
	subs, err := models.GetSubmissions(uid)
	if err != nil {
		s.Data["json"] = err.Error()
		s.Ctx.ResponseWriter.WriteHeader(403)
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
// @Failure 403 user or site invalid
// @router /:site [post]
func (s *SubmissionController) SaveSubmission() {
	uid := s.Ctx.Input.GetData("uid").(bson.ObjectId)
	site := s.GetString(":site")

	if isSiteValid(site) {
		user, err := models.GetUser(uid)
		if err != nil {
			s.Data["json"] = map[string]string{"error": err.Error()}
			s.ServeJSON()
			return
		}
		err = models.AddSubmissions(user, site)
		if err != nil {
			s.Data["json"] = map[string]string{"error": err.Error()}
		} else {
			s.Data["json"] = map[string]string{"status": "submission successfully saved"}
		}
	} else {
		s.Data["json"] = map[string]string{"error": "user or site invalid"}
		s.Ctx.ResponseWriter.WriteHeader(403)
	}
	s.ServeJSON()
}

func isSiteValid(s string) bool {
	for _, vs := range valid_sites {
		if s == vs {
			return true
		}
	}
	return false
}

// @Title Filter
// @Description Filter submissions of user on the basis of status, site and tags
// @Security token_auth read:submission
// @Param	uid		path 	string	false		"UID of user"
// @Param	site		path 	string	true		"Website name"
// @Param	status		query 	string	false		"Submission status"
// @Param	tag 		query	string	false		"Submission tag"
// @Success 200 {object} submission.CodechefSubmission
// @Failure 403 user not exist
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
		s.Ctx.ResponseWriter.WriteHeader(403)
		s.ServeJSON()
		return
	}
	status := s.GetString("status")
	site := s.GetString(":site")
	tag := s.GetString("tag")
	subs, err := models.FilterSubmission(uid, status, tag, site)
	if err != nil {
		s.Data["json"] = err.Error()
		s.Ctx.ResponseWriter.WriteHeader(403)
	} else {
		s.Data["json"] = subs
	}
	s.ServeJSON()
}
