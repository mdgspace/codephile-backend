package controllers

import (
	"github.com/astaxie/beego"
	"github.com/globalsign/mgo/bson"
	"github.com/mdg-iitr/Codephile/models"
)

var valid_sites = []string{"codechef", "codeforces", "spoj", "hackerrank"}

type SubmissionController struct {
	beego.Controller
}

// @Title Get
// @Description Get submissions of user across various platforms
// @Param	uid		path 	string	true		"UID of user"
// @Success 200 {object} submission.Submissions
// @Failure 403 user not exist
// @router /:uid [get]
func (s *SubmissionController) GetSubmission() {
	uid := s.GetString(":uid")
	if uid != "" && bson.IsObjectIdHex(uid) {
		subs, err := models.GetSubmissions(bson.ObjectIdHex(uid))
		if err != nil {
			s.Data["json"] = err.Error()
			s.Ctx.ResponseWriter.WriteHeader(403)
		} else {
			s.Data["json"] = subs
		}
	} else {
		s.Data["json"] = "user not exist"
		s.Ctx.ResponseWriter.WriteHeader(403)
	}
	s.ServeJSON()
}

// @Title Post
// @Description Triggers saving of user's submissions across a particular platform into database
// @Param	uid		path 	string	true		"UID of user"
// @Param	site		path 	string	true		"Platform site name"
// @Success 200 submission successfully saved
// @Failure 403 user or site invalid
// @router /:site/:uid [post]
func (s *SubmissionController) SaveSubmission() {
	uid := s.GetString(":uid")
	site := s.GetString(":site")

	if uid != "" && bson.IsObjectIdHex(uid) && isSiteValid(site) {
		user, err := models.GetUser(bson.ObjectIdHex(uid))
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
