package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/schema"
	"github.com/mdg-iitr/Codephile/models"
	"github.com/mdg-iitr/Codephile/scripts"
	"github.com/mdg-iitr/Codephile/services/auth"
	"os"
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

// @Title Login
// @Description Logs user into the system
// @Param	username		query 	string	true		"The username for login"
// @Param	password		query 	string	true		"The password for login"
// @Success 200 {string} login success
// @Failure 403 user not exist
// @router /login [post]
func (u *UserController) Login() {
	username := u.Ctx.Request.FormValue("username")
	password := u.Ctx.Request.FormValue("password")
	if user, isValid := models.AutheticateUser(username, password); isValid {
		u.Data["json"] = map[string]string{"token": auth.GenerateToken(user.ID.Hex())}
	} else {
		u.Data["json"] = map[string]string{"error": "invalid user credential"}
	}
	u.ServeJSON()
}

// @Title logout
// @Description Logs out current logged in user session
// @Success 200 {string} logout success
// @router /logout [post]
func (u *UserController) Logout() {
	requestToken, _ := request.ParseFromRequest(u.Ctx.Request, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("HMACKEY")), nil
	})
	if requestToken.Valid && !auth.IsTokenExpired(requestToken) {
		err := auth.BlacklistToken(requestToken)
		if err == nil {
			u.Data["json"] = map[string]string{"status": "Logout successful"}
		}
	} else {
		u.Data["json"] = map[string]string{"status": "Invalid Credentials"}
	}
	u.ServeJSON()
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

// @Title Verify site handles
// @Description verify user handles across different websites
// @Param	site		path 	string	true		"site name"
// @Param	handle		query 	string	true		"handle to verify"
// @Success 200 {string} Handle valid
// @Failure 403 incorrect site or handle
// @router /verify/:site [get]
func (u *UserController) Verify() {
	handle := u.GetString("handle")
	site := u.GetString(":site")
	var valid = false
	switch site {
	case "codechef":
		valid = scripts.CheckCodechefHandle(handle)
		break;
	case "codeforces":
		valid = scripts.CheckCodeforcesHandle(handle)
		break;
	case "spoj":
		valid = scripts.CheckSpojHandle(handle)
		break;
	case "hackerrank":
		valid = scripts.CheckHackerrankHandle(handle)
		break;
	}
	if valid {
		u.Data["json"] = map[string]string{"status": "Handle valid"}
	} else {
		u.Ctx.ResponseWriter.WriteHeader(403)
	}
	u.ServeJSON()
}

// @Title Fetch User Info	
// @Description Fetches user info from different websites and store them into the database
// @Param	site		path 	string	true		"site name"
// @Param	handle		query 	string	true		"handle to fetch data from"
// @Param	uid		query 	string	true		"uid of user"
// @Success 200 {object} Success
// @Failure 403 incorrect site or handle
// @router /fetch/:site [get]
func (u *UserController) Fetch() {
	handle := u.GetString("handle")
	site := u.GetString(":site")
	uid := u.GetString("uid")
	if uid != "" && bson.IsObjectIdHex(uid) {
	_ , err := models.AddorUpdateProfile( bson.ObjectIdHex(uid) , site , handle)
	    if err == nil {
		    u.Data["json"] = map[string]string{"status": "Data fetched"}
	    } else {
		      // handle the error
			  u.Data["json"] = map[string]string{"status": "user invalid or database operation failed"}  
	           }
	} else{
		//handle the error(uid of the user isn't valid)
		u.Data["json"] = map[string]string{"status": "uid is not valid"}
	}
	// u.Data["json"] = profile
	u.ServeJSON()
}
