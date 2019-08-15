package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/globalsign/mgo/bson"
	"github.com/mdg-iitr/Codephile/models"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

// Operations about Users
type UserController struct {
	beego.Controller
}

// @Title CreateUser
// @Description create users
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.Id
// @Failure 403 body is empty
// @router /signup [post]
func (u *UserController) CreateUser() {
	user := u.parseRequestBody()
	id, err := models.AddUser(user)
	if err != nil {
		u.Data["json"] = map[string]string{"error": err.Error()}
	} else
	{
		u.Data["json"] = map[string]string{"id": id}
	}
	u.ServeJSON()
}

// @Title GetAll
// @Description get all Users
// @Success 200 {object} models.User
// @router /all [get]
func (u *UserController) GetAll() {
	users := models.GetAllUsers()
	u.Data["json"] = users
	u.ServeJSON()
}

// @Title Get
// @Description get user by uid
// @Param	uid		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 403 :uid is empty
// @router /:uid [get]
func (u *UserController) Get() {
	uid := u.GetString(":uid")
	if uid != "" && bson.IsObjectIdHex(uid) {
		user, err := models.GetUser(bson.ObjectIdHex(uid))
		if err != nil {
			u.Data["json"] = map[string]string{"error": err.Error()}
		} else {
			u.Data["json"] = user
		}
	}
	u.ServeJSON()
}

// @Title Update
// @Description update the user
// @Param	uid		path 	string	true		"The uid you want to update"
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {object} models.User
// @Failure 403 :uid is not int
// @router /:uid [put]
func (u *UserController) Put() {
	uid := u.GetString(":uid")
	if uid != "" && bson.IsObjectIdHex(uid) {
		newUser := u.parseRequestBody()
		uu, err := models.UpdateUser(bson.ObjectIdHex(uid), &newUser)
		if err != nil {
			u.Data["json"] = map[string]string{"error": err.Error()}
		} else {
			u.Data["json"] = uu
		}
	}
	u.ServeJSON()
}

// // @Title Delete
// // @Description delete the user
// // @Param	uid		path 	string	true		"The uid you want to delete"
// // @Success 200 {string} delete success!
// // @Failure 403 uid is empty
// // @router /:uid [delete]
// func (u *UserController) Delete() {
// 	uid := u.GetString(":uid")
// 	models.DeleteUser(uid)
// 	u.Data["json"] = "delete success!"
// 	u.ServeJSON()
// }

// @Title Login
// @Description Logs user into the system
// @Param	username		query 	string	true		"The username for login"
// @Param	password		query 	string	true		"The password for login"
// @Success 200 {string} login success
// @Failure 403 user not exist
// @router /login [post]
func (u *UserController) Login() {
	// username := u.GetString("username")
	// password := u.GetString("password")
	// if models.Login(username, password) {
	// 	u.Data["json"] = "login success"
	// } else {
	// 	u.Data["json"] = "user not exist"
	// }
	// u.ServeJSON()
}

// @Title logout
// @Description Logs out current logged in user session
// @Success 200 {string} logout success
// @router /logout [post]
func (u *UserController) Logout() {
	// u.Data["json"] = "logout success"
	// u.ServeJSON()
}
func (u *UserController) parseRequestBody() models.User {
	var (
		user models.User
		err  error
	)
	if u.Ctx.Request.Header.Get("content-type") == "application/json" {
		err = json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	} else {
		decoder.IgnoreUnknownKeys(true)
		err = decoder.Decode(&user, u.Ctx.Request.PostForm)
	}
	if err != nil {
		panic(err)
	}
	return user
}
