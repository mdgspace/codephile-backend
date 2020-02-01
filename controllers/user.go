package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/schema"
	. "github.com/mdg-iitr/Codephile/conf"
	. "github.com/mdg-iitr/Codephile/errors"
	"github.com/mdg-iitr/Codephile/models"
	"github.com/mdg-iitr/Codephile/models/db"
	"github.com/mdg-iitr/Codephile/models/types"
	"github.com/mdg-iitr/Codephile/scripts"
	"github.com/mdg-iitr/Codephile/services/auth"
	"github.com/mdg-iitr/Codephile/services/firebase"
	"log"
	"net/http"
	"os"
)

var decoder = schema.NewDecoder()

// Operations about Users
type UserController struct {
	beego.Controller
}

// @Title CreateUser
// @Description create users
// @Param	username 		formData	string	true "Username"
// @Param	password		formData 	string	true "Password"
// @Param	fullname		formData 	string	true "Full name of User"
// @Param	institute		formData 	string	false "Name of Institute"
// @Param	handle.codechef	formData	string 	false "Codechef Handle"
// @Param	handle.codeforces	formData	string 	false "Codeforces Handle"
// @Param	handle.hackerrank	formData	string 	false "Hackerrank Handle"
// @Param	handle.spoj		formData	string 	false "Spoj Handle"
// @Success 201 {int} models.types.User.Id
// @Failure 409 username already exists
// @Failure 400 bad request body or blank username/password/full name
// @Failure 500 server_error
// @router /signup [post]
func (u *UserController) CreateUser() {
	user, err := u.parseRequestBody()
	if err != nil {
		u.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		u.Data["json"] = BadInputError("Bad request body")
		u.ServeJSON()
		return
	}
	if user.Username == "" || user.Password == "" || user.FullName == "" {
		u.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		u.Data["json"] = BadInputError("username/password/full name cannot be empty")
		u.ServeJSON()
		return
	}
	id, err := models.AddUser(user)
	if err == UserAlreadyExistError {
		u.Ctx.ResponseWriter.WriteHeader(http.StatusConflict)
		u.Data["json"] = AlreayExistsError("User already exists")
		u.ServeJSON()
		return
	} else if err != nil {
		log.Println(err.Error())
		u.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		u.Data["json"] = InternalServerError("Internal server error")
		u.ServeJSON()
		return
	}
	u.Ctx.ResponseWriter.WriteHeader(http.StatusCreated)
	u.Data["json"] = map[string]string{"id": id}
	u.ServeJSON()
}

// @Title GetAll
// @Description get all Users
// @Security token_auth read:user
// @Success 200 {object} []models.types.User
// @Failure 500 server_error
// @router /all [get]
func (u *UserController) GetAll() {
	users, err := models.GetAllUsers()
	if err != nil {
		log.Println(err.Error())
		u.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		u.Data["json"] = InternalServerError("Internal server error")
		u.ServeJSON()
		return
	}
	u.Data["json"] = users
	u.ServeJSON()
}

// @Title Get
// @Description Get user by uid. Returns logged in user if uid is empty
// @Security token_auth read:user
// @Param	uid		path 	string	false		"uid of user"
// @Success 200 {object} models.types.User
// @Failure 401 : Unauthorized
// @Failure 400 :uid is invalid
// @Failure 404 user not found
// @Failure 500 server_error
// @router / [get]
// @router /:uid [get]
func (u *UserController) Get() {
	uidString := u.GetString(":uid")
	var uid bson.ObjectId
	if bson.IsObjectIdHex(uidString) {
		uid = bson.ObjectIdHex(uidString)
	} else if uidString == "" {
		uid = u.Ctx.Input.GetData("uid").(bson.ObjectId)
	} else {
		u.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		u.Data["json"] = BadInputError("Invalid UID")
		u.ServeJSON()
		return
	}
	user, err := models.GetUser(uid)
	if err != nil {
		u.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		u.Data["json"] = NotFoundError("User not found")
		u.ServeJSON()
		return
	}
	u.Data["json"] = user
	u.ServeJSON()
}

// @Title Update
// @Description update the logged in user
// @Security token_auth write:user
// @Param	username 		formData	string	false "New Username"
// @Param	password		formData 	string	false "New Password"
// @Param	fullname		formData 	string	false "New Full name of User"
// @Param	institute		formData 	string	false "New Name of Institute"
// @Param	handle.codechef	formData	string 	false "New Codechef Handle"
// @Param	handle.codeforces	formData	string 	false "New Codeforces Handle"
// @Param	handle.hackerrank	formData	string 	false "New Hackerrank Handle"
// @Param	handle.spoj		formData	string 	false "New Spoj Handle"
// @Success 202 {object} models.types.User
// @Failure 409 username already exists
// @Failure 400 bad request body or blank username/password/full name
// @Failure 401 : Unauthorized
// @Failure 500 server_error
// @router / [put]
func (u *UserController) Put() {
	uid := u.Ctx.Input.GetData("uid").(bson.ObjectId)
	newUser, err := u.parseRequestBody()
	if err != nil {
		u.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		u.Data["json"] = BadInputError("Bad request body")
		u.ServeJSON()
		return
	}
	uu, err := models.UpdateUser(uid, &newUser)
	if err == UserAlreadyExistError {
		u.Ctx.ResponseWriter.WriteHeader(http.StatusConflict)
		u.Data["json"] = AlreayExistsError("User already exists")
		u.ServeJSON()
		return
	} else if err != nil {
		log.Println(err.Error())
		u.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		u.Data["json"] = InternalServerError("Internal server error")
		u.ServeJSON()
		return
	}
	u.Ctx.ResponseWriter.WriteHeader(http.StatusAccepted)
	u.Data["json"] = uu
	u.ServeJSON()
}

// @Title Login
// @Description Logs user into the system
// @Param	username		formData 	string	true		"The username for login"
// @Param	password		formData 	string	true		"The password for login"
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
		u.Ctx.ResponseWriter.WriteHeader(403)
	}
	u.ServeJSON()
}

// @Title logout
// @Description Logs out current logged in user session
// @Security token_auth write:user
// @Success 200 {string} logout success
// @Failure 401 invalid authentication token
// @Failure 500 server_error
// @router /logout [post]
func (u *UserController) Logout() {
	requestToken, err := request.ParseFromRequest(u.Ctx.Request, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("HMACKEY")), nil
	})
	if err == request.ErrNoTokenInRequest {
		u.Ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		u.Data["json"] = BadInputError("Bad request header")
		u.ServeJSON()
		return
	} else if err != nil {
		log.Println(err.Error())
		u.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		u.Data["json"] = InternalServerError("Internal server error")
		u.ServeJSON()
		return
	}
	if requestToken.Valid && !auth.IsTokenExpired(requestToken) {
		err := auth.BlacklistToken(requestToken)
		if err != nil {
			log.Println(err.Error())
			u.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
			u.Data["json"] = InternalServerError("Internal server error")
			u.ServeJSON()
			return
		}
		u.Data["json"] = map[string]string{"status": "Logout successful"}
	} else {
		u.Ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		u.Data["json"] = map[string]string{"status": "Invalid Credentials"}
	}
	u.ServeJSON()
}

func (u *UserController) parseRequestBody() (types.User, error) {
	var (
		user types.User
		err  error
	)
	if u.Ctx.Request.Header.Get("content-type") == "application/json" {
		err = json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	} else {
		decoder.IgnoreUnknownKeys(true)
		err = decoder.Decode(&user, u.Ctx.Request.PostForm)
	}
	if err != nil {
		log.Println(err.Error())
		return types.User{}, err
	}
	return user, err
}

// @Title Verify site handles
// @Description verify user handles across different websites
// @Security token_auth read:user
// @Param	site		path 	string	true		"site name"
// @Param	handle		query 	string	true		"handle to verify"
// @Success 200 {string} Handle valid
// @Failure 400 invalid contest site
// @Failure 403 incorrect handle
// @router /verify/:site [get]
func (u *UserController) Verify() {
	handle := u.GetString("handle")
	site := u.GetString(":site")
	var valid = false
	if !IsSiteValid(site) {
		u.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		u.Data["json"] = BadInputError("Invalid contest site")
		u.ServeJSON()
		return
	}
	switch site {
	case CODECHEF:
		valid = scripts.CheckCodechefHandle(handle)
		break
	case CODEFORCES:
		valid = scripts.CheckCodeforcesHandle(handle)
		break
	case SPOJ:
		valid = scripts.CheckSpojHandle(handle)
		break
	case HACKERRANK:
		valid = scripts.CheckHackerrankHandle(handle)
		break
	}
	if valid {
		u.Data["json"] = map[string]string{"status": "Handle valid"}
	} else {
		u.Ctx.ResponseWriter.WriteHeader(http.StatusForbidden)
	}
	u.ServeJSON()
}

// @Title Fetch User Info	
// @Description Fetches user info from different websites and store them into the database
// @Security token_auth write:user
// @Param	site		path 	string	true		"site name"
// @Success 201 Success
// @Failure 400 incorrect site or handle
// @Failure 401 Unauthenticated
// @Failure 500 server_error
// @router /fetch/:site [post]
func (u *UserController) Fetch() {
	site := u.GetString(":site")
	uid := u.Ctx.Input.GetData("uid").(bson.ObjectId)
	if !IsSiteValid(site) {
		u.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		u.Data["json"] = BadInputError("Invalid contest site")
		u.ServeJSON()
		return
	}
	_, err := models.AddorUpdateProfile(uid, site)
	if err == UserNotFoundError {
		u.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		u.Data["json"] = BadInputError("Invalid user id")
		u.ServeJSON()
		return
	} else if err != nil {
		log.Println(err.Error())
		u.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		u.Data["json"] = InternalServerError("Internal server error")
		u.ServeJSON()
		return
	}
	u.Ctx.ResponseWriter.WriteHeader(http.StatusCreated)
	u.Data["json"] = map[string]string{"status": "Data fetched"}
	u.ServeJSON()
}

// @Title Fetch All User Profiles And returns them
// @Description Returns info of user(logged in user if uid is empty) from different websites
// @Security token_auth read:user
// @Param	uid		path 	string	false		"UID of user"
// @Success 200 {object} profile.AllProfiles
// @Failure 401 Unauthenticated
// @Failure 403 invalid user
// @router /fetch/ [get]
// @router /fetch/:uid [get]
func (u *UserController) ReturnAllProfiles() {
	uidString := u.GetString(":uid")
	var uid bson.ObjectId
	if bson.IsObjectIdHex(uidString) {
		uid = bson.ObjectIdHex(uidString)
	} else if uidString == "" {
		uid = u.Ctx.Input.GetData("uid").(bson.ObjectId)
	} else {
		u.Ctx.ResponseWriter.WriteHeader(403)
		u.ServeJSON()
		return
	}
	user, err := models.GetProfiles(uid)
	if err != nil {
		u.Data["json"] = err.Error()
		u.Ctx.ResponseWriter.WriteHeader(403)
		log.Println(err.Error())
	} else {
		u.Data["json"] = user
	}
	u.ServeJSON()
}

// @Title Update Profile Pic
// @Description update the profile picture of logged in user
// @Security token_auth write:user
// @Param	image		formData 	file	true		"profile image"
// @Success 200  successful
// @Failure 401 Unauthenticated
// @Failure 403 could not get image
// @router /picture [put]
func (u *UserController) ProfilePic() {
	uid := u.Ctx.Input.GetData("uid").(bson.ObjectId)
	f, fh, err := u.GetFile("image")
	if err != nil {
		u.Data["json"] = "could not get image"
		u.Ctx.ResponseWriter.WriteHeader(403)
		u.ServeJSON()
		return
	}
	newPic, err := firebase.AddFile(f, fh, models.GetPicture(uid))
	if err != nil {
		log.Println(err)
		return
	}
	err = models.UpdatePicture(uid, newPic)
	if err != nil {
		u.Data["json"] = err.Error()
		u.Ctx.ResponseWriter.WriteHeader(403)
		u.ServeJSON()
		return
	}
	u.Data["json"] = "successful"
	u.ServeJSON()
}

// @Title username available
// @Description checks if username is available
// @Param	username		query 	string	true		"Username"
// @Success 200  available
// @Failure 403 unavailable
// @router /available [get]
func (u *UserController) IsAvailable() {
	username := u.GetString("username")
	if username == "" {
		u.Ctx.ResponseWriter.WriteHeader(403)
		return
	}
	collection := db.NewUserCollectionSession()
	defer collection.Close()
	c, err := collection.Collection.Find(bson.M{"username": username}).Count()
	if err != nil {
		log.Println(err.Error())
		u.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
	}
	//username available
	if c == 0 {
		u.Data["json"] = "available"
		u.ServeJSON()
		return
	}
	u.Ctx.ResponseWriter.WriteHeader(403)
	u.Data["json"] = "unavailable"
	u.ServeJSON()
}
