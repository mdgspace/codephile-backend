package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/getsentry/sentry-go"
	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/schema"
	. "github.com/mdg-iitr/Codephile/conf"
	. "github.com/mdg-iitr/Codephile/errors"
	"github.com/mdg-iitr/Codephile/models"
	"github.com/mdg-iitr/Codephile/models/types"
	"github.com/mdg-iitr/Codephile/scrappers"
	"github.com/mdg-iitr/Codephile/services/auth"
	"github.com/mdg-iitr/Codephile/services/firebase"
	"github.com/mdg-iitr/Codephile/services/redis"
	"github.com/mdg-iitr/Codephile/services/worker"
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
// @Success 201 {int} types.User.Id
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
		u.Data["json"] = AlreadyExistsError("User already exists")
		u.ServeJSON()
		return
	} else if err != nil {
		hub := sentry.GetHubFromContext(u.Ctx.Request.Context())
		hub.CaptureException(err)
		log.Println(err.Error())
		u.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		u.Data["json"] = InternalServerError("Internal server error")
		u.ServeJSON()
		return
	}
	// Keep track of which IP creates how many user
	client := redis.GetRedisClient()
	err = client.Incr(u.Ctx.Request.RemoteAddr).Err()
	if err != nil {
		hub := sentry.GetHubFromContext(u.Ctx.Request.Context())
		hub.CaptureException(err)
	}
	u.Ctx.ResponseWriter.WriteHeader(http.StatusCreated)
	u.Data["json"] = map[string]string{"id": id}
	u.ServeJSON()
}

// @Title GetAll
// @Description get all Users
// @Security token_auth read:user
// @Success 200 {object} []types.User
// @Failure 500 server_error
// @router /all [get]
func (u *UserController) GetAll() {
	users, err := models.GetAllUsers()
	if err != nil {
		hub := sentry.GetHubFromContext(u.Ctx.Request.Context())
		hub.CaptureException(err)
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
// @Success 200 {object} types.User
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
// @Param	fullname		formData 	string	false "New Full name of User"
// @Param	institute		formData 	string	false "New Name of Institute"
// @Param	handle.codechef	formData	string 	false "New Codechef Handle"
// @Param	handle.codeforces	formData	string 	false "New Codeforces Handle"
// @Param	handle.hackerrank	formData	string 	false "New Hackerrank Handle"
// @Param	handle.spoj		formData	string 	false "New Spoj Handle"
// @Success 202 {object} types.User
// @Failure 409 username already exists
// @Failure 400 bad request body
// @Failure 401 : Unauthorized
// @Failure 404 : User not found
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
		u.Data["json"] = AlreadyExistsError("User already exists")
		u.ServeJSON()
		return
	} else if err != nil {
		hub := sentry.GetHubFromContext(u.Ctx.Request.Context())
		hub.CaptureException(err)
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
	if user, isValid := models.AuthenticateUser(username, password); isValid {
		u.Data["json"] = map[string]string{"token": auth.GenerateToken(user.ID.Hex())}
	} else {
		u.Data["json"] = map[string]string{"error": "invalid user credential"}
		u.Ctx.ResponseWriter.WriteHeader(403)
	}
	u.ServeJSON()
}

// @Title Password-Reset
// @Description Resets the password of the user
// @Param	email		formData 	string	true		"The email for of the user"
// @Success 200 {string} email sent
// @Failure 403 user not exist
// @router /password-reset [post]
func (u *UserController) PasswordReset() {
	email := u.Ctx.Request.FormValue("email")
	if isValid := models.PasswordResetEmail(email); isValid {
		u.Data["json"] = map[string]string{"email": "sent"}
	} else {
		u.Data["json"] = map[string]string{"error": "invalid email"}
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
		hub := sentry.GetHubFromContext(u.Ctx.Request.Context())
		hub.CaptureException(err)
		log.Println(err.Error())
		u.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		u.Data["json"] = InternalServerError("Internal server error")
		u.ServeJSON()
		return
	}
	if requestToken.Valid && !auth.IsTokenExpired(requestToken) {
		err := auth.BlacklistToken(requestToken)
		if err != nil {
			hub := sentry.GetHubFromContext(u.Ctx.Request.Context())
			hub.CaptureException(err)
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
	scrapper, err := scrappers.NewScrapper(site, handle)
	if err != nil {
		u.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		u.Data["json"] = BadInputError("Invalid contest site")
		u.ServeJSON()
		return
	}
	valid := scrapper.CheckHandle()
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
	job := worker.NewJob(uid, site, models.AddOrUpdateProfile)
	err := worker.Enqueue(job)
	if err != nil {
		u.Ctx.ResponseWriter.WriteHeader(http.StatusServiceUnavailable)
		u.Data["json"] = UnavailableError("slow down cowboy")
		u.ServeJSON()
		return
	}
	u.Ctx.ResponseWriter.WriteHeader(http.StatusCreated)
	u.Data["json"] = map[string]string{"status": "data will be fetched in a moment"}
	u.ServeJSON()
}

// @Title Fetch All User Profiles And returns them
// @Description Returns info of user(logged in user if uid is empty) from different websites
// @Security token_auth read:user
// @Param	uid		path 	string	false		"UID of user"
// @Success 200 {object} types.AllProfiles
// @Failure 401 Unauthenticated
// @Failure 400 invalid user
// @Failure 500 server_error
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
		u.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		u.Data["json"] = BadInputError("Invalid UID")
		u.ServeJSON()
		return
	}
	user, err := models.GetProfiles(uid)
	if err != nil {
		u.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		u.Data["json"] = NotFoundError("Profile not found")
		u.ServeJSON()
		return
	}
	u.Data["json"] = user

	u.ServeJSON()
}

// @Title Update Profile Pic
// @Description update the profile picture of logged in user
// @Security token_auth write:user
// @Param	image		formData 	file	true		"profile image"
// @Success 201  successful
// @Failure 401 Unauthenticated
// @Failure 400 could not get image
// @router /picture [put]
func (u *UserController) ProfilePic() {
	uid := u.Ctx.Input.GetData("uid").(bson.ObjectId)
	f, fh, err := u.GetFile("image")
	if err != nil {
		u.Data["json"] = "could not get image"
		u.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		u.ServeJSON()
		return
	}
	newPic, err := firebase.AddFile(f, fh, models.GetPicture(uid))
	if err != nil {
		hub := sentry.GetHubFromContext(u.Ctx.Request.Context())
		hub.CaptureException(err)
		log.Println(err.Error())
		u.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		u.Data["json"] = InternalServerError("Internal server error")
		u.ServeJSON()
		return
	}
	err = models.UpdatePicture(uid, newPic)
	if err != nil {
		hub := sentry.GetHubFromContext(u.Ctx.Request.Context())
		hub.CaptureException(err)
		log.Println(err.Error())
		u.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		u.Data["json"] = InternalServerError("Internal server error")
		u.ServeJSON()
		return
	}
	u.Ctx.ResponseWriter.WriteHeader(http.StatusCreated)
	u.Data["json"] = "successful"
	u.ServeJSON()
}

// @Title username available
// @Description checks if username is available
// @Param	username		query 	string	true		"Username"
// @Success 200  available
// @Failure 400 empty username
// @Failure 403 unavailable
// @router /available [get]
func (u *UserController) IsAvailable() {
	username := u.GetString("username")
	if username == "" {
		u.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	exists, err := models.UserExists(username)
	if err != nil {
		hub := sentry.GetHubFromContext(u.Ctx.Request.Context())
		hub.CaptureException(err)
		log.Println(err.Error())
		u.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		u.Data["json"] = InternalServerError("Internal server error")
		u.ServeJSON()
		return
	}
	if !exists {
		u.Data["json"] = map[string]string{"status": "available"}
		u.ServeJSON()
		return
	}
	u.Ctx.ResponseWriter.WriteHeader(403)
	u.Data["json"] = "unavailable"
	u.ServeJSON()
}

// @Title Password Change
// @Description Changes password of the user
// @Security token_auth write:user
// @Param	data body types.UpdatePassword  true "JSON body containing old and new password
// @Success 200 {string} success
// @Failure 401 Unauthenticated
// @Failure 400 bad request
// @Failure 403 old password incorrect
// @Failure 500 server error
// @router /password-reset [post]
func (u *UserController) PasswordChange() {
	uid := u.Ctx.Input.GetData("uid").(bson.ObjectId)
	var passwordUpdateRequest types.UpdatePassword
	err := json.Unmarshal(u.Ctx.Input.RequestBody, &passwordUpdateRequest)
	if err != nil {
		u.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		u.Data["json"] = BadInputError("json body is malformed")
		u.ServeJSON()
		return
	}
	err = models.UpdatePassword(uid, passwordUpdateRequest)
	if err == PasswordIncorrectError {
		u.Ctx.ResponseWriter.WriteHeader(http.StatusForbidden)
		u.Data["json"] = BadInputError("old password is incorrect or new password is empty")
		u.ServeJSON()
		return
	} else if err != nil {
		hub := sentry.GetHubFromContext(u.Ctx.Request.Context())
		hub.CaptureException(err)
		log.Println(err.Error())
		u.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		u.Data["json"] = InternalServerError("server error.. report to admin")
		u.ServeJSON()
		return
	}
	u.Data["json"] = "success"
	u.ServeJSON()
}

// @Title Filter
// @Description Filter users on basis of institute name
// @Security token_auth read:user
// @Param	institute		query 	string	true		"Institute Name"
// @Success 200 {object} []types.SearchDoc
// @Failure 404 "No users found"
// @Failure 500 server_error
// @router /filter [get]
func (u *UserController) FilterUsers() {
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
		u.Data["json"] = NotFoundError("no such user")
		u.ServeJSON()
		return
	}
	u.Data["json"] = res
	u.ServeJSON()
}
