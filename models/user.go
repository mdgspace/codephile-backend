package models

import (
	"bytes"
	"context"
	"html/template"
	"log"
	"math/rand"
	"time"

	"github.com/astaxie/beego"
	"github.com/getsentry/sentry-go"
	"github.com/mdg-iitr/Codephile/services/firebase"

	"github.com/mdg-iitr/Codephile/services/mail"

	// "github.com/globalsign/mgo"
	// "github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/google/uuid"
	. "github.com/mdg-iitr/Codephile/conf"
	. "github.com/mdg-iitr/Codephile/errors"
	"github.com/mdg-iitr/Codephile/models/db"
	"github.com/mdg-iitr/Codephile/models/types"
	"github.com/mdg-iitr/Codephile/services/redis"
	"golang.org/x/crypto/bcrypt"
)

var (
	getCodechefSolvesQuery = bson.M{
		"$size": bson.M{
			"$filter": bson.M{
				"input": "$submissions",
				"as":    "sub",
				"cond": bson.M{
					"$and": []bson.M{
						{
							"$regexMatch": bson.M{
								"input": "$$sub.url",
								"regex": primitive.Regex{Pattern: "^" + "https://www.codechef.com"},
							},
						},
						{"$eq": []string{"$$sub.status", StatusCorrect}}},
				},
			},
		},
	}
	getCodeforcesSolvesQuery = bson.M{
		"$size": bson.M{
			"$filter": bson.M{
				"input": "$submissions",
				"as":    "sub",
				"cond": bson.M{
					"$and": []bson.M{
						{
							"$regexMatch": bson.M{
								"input": "$$sub.url",
								"regex": primitive.Regex{Pattern: "^" + "http://codeforces.com"},
							},
						},
						{"$eq": []string{"$$sub.status", StatusCorrect}}},
				},
			},
		},
	}
	getHackerrankSolvesQuery = bson.M{
		"$size": bson.M{
			"$filter": bson.M{
				"input": "$submissions",
				"as":    "sub",
				"cond": bson.M{
					"$and": []bson.M{
						{
							"$regexMatch": bson.M{
								"input": "$$sub.url",
								"regex": primitive.Regex{Pattern: "^" + "https://www.hackerrank.com"},
							},
						},
						{"$eq": []string{"$$sub.status", StatusCorrect}}},
				},
			},
		},
	}
	getSpojSolvesQuery = bson.M{
		"$size": bson.M{
			"$filter": bson.M{
				"input": "$submissions",
				"as":    "sub",
				"cond": bson.M{
					"$and": []bson.M{
						{
							"$regexMatch": bson.M{
								"input": "$$sub.url",
								"regex": primitive.Regex{Pattern: "^" + "https://www.spoj.com"},
							},
						},
						{"$eq": []string{"$$sub.status", StatusCorrect}}},
				},
			},
		},
	}
	getFollowingCountQuery = bson.M{
		"$size": "$followingUsers",
	}
)

func AddUser(u types.User) (string, error) {
	u.ID = primitive.NewObjectID()
	u.Verified = false
	defaultPic := beego.AppConfig.Strings("DEFAULT_PICS")
	if len(defaultPic) > 0 {
		u.Picture = firebase.URLFromName(defaultPic[rand.Intn(len(defaultPic))])
	}
	collection := db.NewUserCollectionSession()
	defer collection.Close()
	//hashing the password
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	//data type of hash is []byte
	if err != nil {
		return "", err
	}
	u.Password = string(hash)
	_, err = collection.Collection.InsertOne(context.TODO(), u)
	if err != nil {
		return "", UserAlreadyExistError
	}
	return u.ID.Hex(), nil
}

func GetUser(uid primitive.ObjectID) (*types.User, error) {
	var user types.User
	var err error
	collection := db.NewUserCollectionSession()
	defer collection.Close()
	err = collection.Collection.FindOne(context.TODO(), bson.M{"_id":uid}, options.FindOne().SetProjection(bson.M{
		"_id": 1, "username": 1, "email": 1, "handle": 1, "lastfetched": 1, "profiles": 1,
		"picture": 1, "fullname": 1, "institute": 1, "submissions": bson.M{"$slice": 5}})).Decode(&user)
	//fmt.Println(err.Error())
	if err != nil {
		return nil, err
	}
	pipe, err := collection.Collection.Aggregate(context.TODO(), []bson.M{
		{
			"$match": bson.M{
				"_id": uid,
			},
		},
		{
			"$project": bson.M{
				"_id":              0,
				"following":        getFollowingCountQuery,
				"codechefSolves":   getCodechefSolvesQuery,
				"codeforcesSolves": getCodeforcesSolvesQuery,
				"hackerrankSolves": getHackerrankSolvesQuery,
				"spojSolves":       getSpojSolvesQuery,
			}},
	}, options.Aggregate().SetBatchSize(1))
	var res map[string]int
	err = pipe.Decode(&res)
	if err != nil {
		return nil, err
	}
	user.NoOfFollowing = res["following"]
	user.SolvedProblemsCount = types.SolvedProblemsCount{
		Codechef:   res["codechefSolves"],
		Codeforces: res["codeforcesSolves"],
		Hackerrank: res["hackerrankSolves"],
		Spoj:       res["spojSolves"],
	}
	return &user, nil
}

func GetAllUsers() ([]types.User, error) {
	var users []types.User
	collection := db.NewUserCollectionSession()
	defer collection.Close()
	cursor, err := collection.Collection.Find(context.TODO(), bson.M{}, options.Find().SetProjection(bson.M{"_id": 1, "username": 1, "email": 1,
		"handle": 1, "lastfetched": 1, "profiles": 1,
		"picture": 1, "fullname": 1, "institute": 1, "submissions": bson.M{"$slice": 5}}))
	cursor.All(context.TODO(), &users)
	if err != nil {
		return nil, err
	}
	pipe, err := collection.Collection.Aggregate(context.TODO(), []bson.M{
		{
			"$project": bson.M{
				"_id":              0,
				"following":        getFollowingCountQuery,
				"codechefSolves":   getCodechefSolvesQuery,
				"codeforcesSolves": getCodeforcesSolvesQuery,
				"hackerrankSolves": getHackerrankSolvesQuery,
				"spojSolves":       getSpojSolvesQuery,
			}},
	})
	var res []map[string]int
	err = pipe.All(context.TODO(), &res)
	if err != nil {
		return nil, err
	}
	for i := range users {
		users[i].SolvedProblemsCount = types.SolvedProblemsCount{
			Codechef:   res[i]["codechefSolves"],
			Codeforces: res[i]["codeforcesSolves"],
			Hackerrank: res[i]["hackerrankSolves"],
			Spoj:       res[i]["spojSolves"],
		}
		users[i].NoOfFollowing = res[i]["following"]
	}
	return users, nil
}

func GetHandle(uid primitive.ObjectID) (types.Handle, error) {
	var user types.User
	collection := db.NewUserCollectionSession()
	defer collection.Close()
	err := collection.Collection.FindOne(context.TODO(), bson.M{"_id": uid}, options.FindOne().SetProjection(bson.M{"handle": 1})).Decode(&user)
	if err != nil {
		return types.Handle{}, err
	}
	return user.Handle, nil
}

func UpdateUser(uid primitive.ObjectID, uu *types.User, ctx context.Context) (a *types.User, err error) {
	var updateDoc = bson.M{}
	newHandle, err := GetHandle(uid)
	var UpdatedSites []string
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	if uu.Username != "" {
		updateDoc["username"] = uu.Username
	}
	if uu.Institute != "" {
		updateDoc["institute"] = uu.Institute
	}
	if uu.FullName != "" {
		updateDoc["fullname"] = uu.FullName
	}
	if uu.Handle.Codechef != "" && uu.Handle.Codechef != newHandle.Codechef {
		updateDoc["handle.codechef"] = uu.Handle.Codechef
		newHandle.Codechef = uu.Handle.Codechef
		UpdatedSites = append(UpdatedSites, CODECHEF)
	}
	if uu.Handle.Codeforces != "" && uu.Handle.Codeforces != newHandle.Codeforces {
		updateDoc["handle.codeforces"] = uu.Handle.Codeforces
		newHandle.Codeforces = uu.Handle.Codeforces
		UpdatedSites = append(UpdatedSites, CODEFORCES)
	}
	if uu.Handle.Hackerearth != "" {
		updateDoc["handle.hackerearth"] = uu.Handle.Hackerearth
		newHandle.Hackerearth = uu.Handle.Hackerearth
		// UpdatedSites = append(UpdatedSites, HACKEREARTH)
	}
	if uu.Handle.Hackerrank != "" && uu.Handle.Hackerrank != newHandle.Hackerrank {
		updateDoc["handle.hackerrank"] = uu.Handle.Hackerrank
		newHandle.Hackerrank = uu.Handle.Hackerrank
		UpdatedSites = append(UpdatedSites, HACKERRANK)
	}
	if uu.Handle.Spoj != "" && uu.Handle.Spoj != newHandle.Spoj {
		updateDoc["handle.spoj"] = uu.Handle.Spoj
		newHandle.Spoj = uu.Handle.Spoj
		UpdatedSites = append(UpdatedSites, SPOJ)
	}
	if len(updateDoc) != 0 {
		collection := db.NewUserCollectionSession()
		defer collection.Close()
		_, err = collection.Collection.UpdateOne(context.TODO(), bson.M{"_id": uid}, bson.M{"$set": updateDoc})
		if err == mongo.ErrNoDocuments {
			return nil, UserNotFoundError
		} else if err != nil {
			log.Println(err.Error())
			return nil, UserAlreadyExistError
		}
	}

	go func() {
		for _, value := range UpdatedSites {
			hub := sentry.GetHubFromContext(ctx)
			err = DeleteSubmissions(uid, value)
			if err != nil {
				hub.CaptureException(err)
			}
			err = ResetProfile(uid, value)
			if err != nil {
				hub.CaptureException(err)
			}
			err = AddSubmissions(uid, value, ctx)
			if err != nil {
				hub.CaptureException(err)
			}
			err = AddOrUpdateProfile(uid, value, ctx)
			if err != nil {
				hub.CaptureException(err)
			}
		}
	}()

	u, err := GetUser(uid)
	if err != nil {
		return nil, err
	}
	return u, err
}

func AuthenticateUser(username string, password string) (*types.User, error) {
	var user types.User
	collection := db.NewUserCollectionSession()
	defer collection.Close()
	err := collection.Collection.FindOne(context.TODO(), bson.M{"username": username}, options.FindOne().SetProjection(bson.M{"password": 1, "verified": 1})).Decode(&user)
	//fmt.Println(err.Error())
	if err != nil {
		//log.Println(err)
		return nil, UserNotFoundError
	}

	err2 := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err2 != nil {
		//log.Println(err2)
		return nil, UserNotFoundError
	}
	if !user.Verified {
		return nil, UserUnverifiedError
	}
	return &user, nil
}

func UpdatePicture(uid primitive.ObjectID, url string) error {
	coll := db.NewUserCollectionSession()
	defer coll.Close()
	_, err := coll.Collection.UpdateOne(context.TODO(), bson.M{"_id": uid}, bson.M{"$set": bson.M{"picture": url}})
	return err
}

func VerifyEmail(uid primitive.ObjectID, ctx context.Context) error {
	sess := db.NewUserCollectionSession()
	defer sess.Close()
	coll := sess.Collection
	_, err := coll.UpdateOne(context.TODO(), bson.M{"_id": uid}, bson.M{"$set": bson.M{"verified": true}})
	if err == mongo.ErrNoDocuments {
		return UserNotFoundError
	} else if err != nil {
		return err
	}
	go func() {
		for _, value := range ValidSites {
			_ = AddSubmissions(uid, value, ctx)
			_ = AddOrUpdateProfile(uid, value, ctx)
		}
	}()
	return nil
}

func GetPicture(uid primitive.ObjectID) string {
	var user types.User
	coll := db.NewUserCollectionSession()
	defer coll.Close()
	err := coll.Collection.FindOne(context.TODO(), bson.M{"_id": uid}, options.FindOne().SetProjection(bson.M{"picture": 1})).Decode(&user)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	return user.Picture
}

func CheckUsernameExists(username string) (bool, error) {
	collection := db.NewUserCollectionSession()
	defer collection.Close()
	c, err := collection.Collection.CountDocuments(context.TODO(), bson.M{"username": username})
	if err != nil {
		log.Println(err.Error())
		return false, err
	}
	if c > 0 {
		return true, nil
	}
	return false, nil
}

func CheckEmailExists(email string) (bool, error) {
	collection := db.NewUserCollectionSession()
	defer collection.Close()
	c, err := collection.Collection.CountDocuments(context.TODO(), bson.M{"email": email})
	if err != nil {
		log.Println(err.Error())
		return false, err
	}
	if c > 0 {
		return true, nil
	}
	return false, nil
}
func UidExists(uid primitive.ObjectID) (bool, error) {
	collection := db.NewUserCollectionSession()
	defer collection.Close()
	c, err := collection.Collection.CountDocuments(context.TODO(), bson.M{"_id": uid})
	if err != nil {
		log.Println(err.Error())
		return false, err
	}
	if c > 0 {
		return true, nil
	}
	return false, nil
}

//checks if the user is verified, returns error if user doesn't exists
func IsUserVerified(uid primitive.ObjectID) (bool, error, string) {
	sess := db.NewUserCollectionSession()
	defer sess.Close()
	coll := sess.Collection
	var user types.User
	err := coll.FindOne(context.TODO(), bson.M{"_id": uid}, options.FindOne().SetProjection(bson.M{"email": 1, "verified": 1})).Decode(&user)
	if err != nil {
		return false, UserNotFoundError, ""
	}
	return user.Verified, nil, user.Email
}

func PasswordResetEmail(email string, hostName string, ctx context.Context) bool {
	collection := db.NewUserCollectionSession()
	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub()
	}
	defer collection.Close()
	var user types.User
	err := collection.Collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		hub.CaptureException(err)
		return false
	}
	client := redis.GetRedisClient()
	uniq_id := uuid.New().String()
	_, err = client.Set(user.ID.Hex(), uniq_id, time.Hour).Result()
	if err != nil {
		log.Println(err.Error())
		return false
	}
	link := hostName + "/v1/user/password-reset/" + uniq_id + "/" + user.ID.Hex()
	t := template.New("reset_email.html")
	var err1 error
	t, err1 = t.ParseFiles("views/reset_email.html")
	if err1 != nil {
		hub.CaptureException(err1)
		log.Println(err1.Error())
		return false
	}
	var tpl bytes.Buffer
	if err2 := t.Execute(&tpl, map[string]string{"link": link}); err2 != nil {
		hub.CaptureException(err2)
		log.Println(err2.Error())
		return false
	}
	body := tpl.String()
	go mail.SendMail(email, "Codephile Password Reset", body, ctx)
	return true
}

func SearchUser(query string, c int) ([]types.SearchDoc, error) {
	sess := db.NewUserCollectionSession()
	defer sess.Close()

	search := bson.M{
		"$search": bson.M{
			"index": "name_search",
			"text": bson.M{
				"query": query,
				"path":  bson.M{"wildcard": "*"},
				"fuzzy": bson.M{
					"maxEdits":      2,
					"maxExpansions": 50,
				},
			},
		},
	}
	limit := bson.M{"$limit": c}
	project := bson.M{
		"$project": bson.M{
			"_id":       1,
			"username":  1,
			"fullname":  1,
			"institute": 1,
			"picture":   1,
			"handle":    1,
		},
	}
	pipe, err := sess.Collection.Aggregate(context.TODO(), []bson.M{
		search,
		limit,
		project,
	})
	var result []types.SearchDoc
	err = pipe.All(context.TODO(), &result)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return result, nil
}

func ResetPassword(id primitive.ObjectID, newPassword string) error {
	sess := db.NewUserCollectionSession()
	defer sess.Close()
	coll := sess.Collection
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = coll.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": bson.M{"password": string(hash)}})
	return err
}

// Updates the password of a given uid
func UpdatePassword(uid primitive.ObjectID, updatePasswordRequest types.UpdatePassword) error {
	sess := db.NewUserCollectionSession()
	defer sess.Close()
	coll := sess.Collection
	var u types.User
	err := coll.FindOne(context.TODO(), bson.M{"_id": uid}, options.FindOne().SetProjection(bson.M{"password": 1})).Decode(&u)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(updatePasswordRequest.OldPassword))
	if err != nil || updatePasswordRequest.NewPassword == "" {
		return PasswordIncorrectError
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(updatePasswordRequest.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = coll.UpdateOne(context.TODO(), bson.M{"_id": uid}, bson.M{"$set": bson.M{"password": string(hash)}})
	return err
}

func FilterUsers(instituteName string) ([]types.SearchDoc, error) {
	sess := db.NewUserCollectionSession()
	defer sess.Close()
	coll := sess.Collection
	var result []types.SearchDoc
	cursor, err := coll.Find(context.TODO(), bson.M{"institute": instituteName}, options.Find().SetProjection(bson.M{"_id": 1, "username": 1, "email": 1,
		"handle": 1, "picture": 1, "fullname": 1, "institute": 1}))
	cursor.All(context.TODO(), &result)
	return result, err
}
