package models

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/mdg-iitr/Codephile/services/mail"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/google/uuid"
	. "github.com/mdg-iitr/Codephile/conf"
	. "github.com/mdg-iitr/Codephile/errors"
	"github.com/mdg-iitr/Codephile/models/db"
	"github.com/mdg-iitr/Codephile/models/types"
	search "github.com/mdg-iitr/Codephile/services/elastic"
	"github.com/mdg-iitr/Codephile/services/redis"
	"github.com/olivere/elastic/v7"
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
								"regex": bson.RegEx{Pattern: "^" + "http://www.codechef.com"},
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
								"regex": bson.RegEx{Pattern: "^" + "http://codeforces.com"},
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
								"regex": bson.RegEx{Pattern: "^" + "https://www.hackerrank.com"},
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
								"regex": bson.RegEx{Pattern: "^" + "https://www.spoj.com"},
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
	u.ID = bson.NewObjectId()
	u.Verified = false
	collection := db.NewUserCollectionSession()
	defer collection.Close()
	//hashing the password
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	//data type of hash is []byte
	if err != nil {
		return "", err
	}
	u.Password = string(hash)
	err = collection.Collection.Insert(u)
	if err != nil {
		return "", UserAlreadyExistError
	}
	client := search.GetElasticClient()
	_, err = client.Index().Index("codephile").BodyJson(types.SearchDoc{
		ID:        u.ID,
		Username:  u.Username,
		FullName:  u.FullName,
		Institute: u.Institute,
		Picture:   u.Picture,
		Handle:    u.Handle,
	}).Id(u.ID.Hex()).Refresh("true").Do(context.Background())

	if err != nil {
		return "", err
	}

	go func() {
		for _, value := range ValidSites {
			_ = AddSubmissions(u.ID, value)
			_ = AddOrUpdateProfile(u.ID, value)
		}
	}()

	return u.ID.Hex(), nil
}

func GetUser(uid bson.ObjectId) (*types.User, error) {
	var user types.User
	collection := db.NewUserCollectionSession()
	defer collection.Close()
	err := collection.Collection.FindId(uid).Select(bson.M{"_id": 1, "username": 1, "email": 1,
		"handle": 1, "lastfetched": 1, "profiles": 1,
		"picture": 1, "fullname": 1, "institute": 1, "submissions": bson.M{"$slice": 5}}).One(&user)
	//fmt.Println(err.Error())
	if err != nil {
		return nil, err
	}
	pipe := collection.Collection.Pipe([]bson.M{
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
	})
	var res map[string]int
	err = pipe.One(&res)
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
	err := collection.Collection.Find(nil).Select(bson.M{"_id": 1, "username": 1, "email": 1,
		"handle": 1, "lastfetched": 1, "profiles": 1,
		"picture": 1, "fullname": 1, "institute": 1, "submissions": bson.M{"$slice": 5}}).All(&users)
	if err != nil {
		return nil, err
	}
	pipe := collection.Collection.Pipe([]bson.M{
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
	err = pipe.All(&res)
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

func GetHandle(uid bson.ObjectId) (types.Handle, error) {
	var user types.User
	collection := db.NewUserCollectionSession()
	defer collection.Close()
	err := collection.Collection.FindId(uid).Select(bson.M{"handle": 1}).One(&user)
	if err != nil {
		return types.Handle{}, err
	}
	return user.Handle, nil
}

func UpdateUser(uid bson.ObjectId, uu *types.User) (a *types.User, err error) {
	var updateDoc = bson.M{}
	var elasticDoc = map[string]interface{}{}
	newHandle, err := GetHandle(uid)
	var UpdatedSites []string
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	if uu.Username != "" {
		updateDoc["username"] = uu.Username
		elasticDoc["username"] = uu.Username
	}
	if uu.Institute != "" {
		updateDoc["institute"] = uu.Institute
		elasticDoc["institute"] = uu.Institute
	}
	if uu.FullName != "" {
		updateDoc["fullname"] = uu.FullName
		elasticDoc["fullname"] = uu.FullName
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
	elasticDoc["handle"] = newHandle
	if len(updateDoc) != 0 {
		collection := db.NewUserCollectionSession()
		defer collection.Close()
		err = collection.Collection.UpdateId(uid, bson.M{"$set": updateDoc})
		if err == mgo.ErrNotFound {
			return nil, UserNotFoundError
		} else if err != nil {
			log.Println(err.Error())
			return nil, UserAlreadyExistError
		}
		client := search.GetElasticClient()
		_, err = client.Update().Index("codephile").Id(uid.Hex()).Doc(elasticDoc).Do(context.Background())
		if err != nil {
			log.Println(err.Error())
		}
	}

	go func() {
		for _, value := range UpdatedSites {
			_ = DeleteSubmissions(uid, value)
			_ = AddSubmissions(uid, value)
		}
	}()

	u, err := GetUser(uid)
	if err != nil {
		return nil, err
	}
	return u, err
}

func AuthenticateUser(username string, password string) (*types.User, bool) {
	var user types.User
	collection := db.NewUserCollectionSession()
	defer collection.Close()
	err := collection.Collection.Find(bson.M{"username": username}).Select(bson.M{"password": 1, "verified": 1}).One(&user)
	//fmt.Println(err.Error())
	if err != nil {
		//log.Println(err)
		return nil, false
	}

	err2 := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err2 != nil {
		//log.Println(err2)
		return nil, false
	} else {
		return &user, true
	}

}

func UpdatePicture(uid bson.ObjectId, url string) error {
	client := search.GetElasticClient()
	_, err := client.Update().Index("codephile").Id(uid.Hex()).Doc(map[string]interface{}{"picture": url}).Do(context.Background())
	if err != nil {
		log.Println(err.Error())
	}
	coll := db.NewUserCollectionSession()
	defer coll.Close()
	return coll.Collection.UpdateId(uid, bson.M{"$set": bson.M{"picture": url}})
}

func GetPicture(uid bson.ObjectId) string {
	var user types.User
	coll := db.NewUserCollectionSession()
	defer coll.Close()
	err := coll.Collection.FindId(uid).Select(bson.M{"picture": 1}).One(&user)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	return user.Picture
}

func UserExists(username string) (bool, error) {
	collection := db.NewUserCollectionSession()
	defer collection.Close()
	c, err := collection.Collection.Find(bson.M{"username": username}).Count()
	if err != nil {
		log.Println(err.Error())
		return false, err
	}
	if c > 0 {
		return true, nil
	}
	return false, nil
}

func PasswordResetEmail(email string) bool {
	collection := db.NewUserCollectionSession()
	defer collection.Close()
	var user types.User
	err := collection.Collection.Find(bson.M{"email": email}).One(&user)
	if err != nil {
		return false
	}
	client := redis.GetRedisClient()
	uniq_id := uuid.New().String()
	_, err = client.Set(user.ID.Hex(), uniq_id, 0).Result()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	link := "http://localhost:8080/v1/user/password-reset/" + uniq_id + "/" + user.ID.Hex()
	body := "Please reset your password by clicking on the following link: \n" + link
	go mail.SendMail(email, "Codephile Password Reset", body)
	return true
}

func SearchUser(query string, c int) ([]types.SearchDoc, error) {
	pq := elastic.NewQueryStringQuery("*" + query + "*").
		Field("username").Field("fullname").
		Field("handle.codechef").Field("handle.spoj").
		Field("handle.codeforces").Field("handle.hackerrank").
		Fuzziness("4")
	client := search.GetElasticClient()
	result, err := client.Search().Index("codephile").
		Pretty(false).Query(pq).Size(c).
		Do(context.Background())
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	results := make([]types.SearchDoc, 0, result.TotalHits())
	for _, hit := range result.Hits.Hits {
		var result types.SearchDoc
		err := json.Unmarshal(hit.Source, &result)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func ResetPassword(id bson.ObjectId, newPassword string) error {
	sess := db.NewUserCollectionSession()
	defer sess.Close()
	coll := sess.Collection
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return coll.UpdateId(id, bson.M{"$set": bson.M{"password": string(hash)}})
}

// Updates the password of a given uid
func UpdatePassword(uid bson.ObjectId, updatePasswordRequest types.UpdatePassword) error {
	sess := db.NewUserCollectionSession()
	defer sess.Close()
	coll := sess.Collection
	var u types.User
	err := coll.FindId(uid).Select(bson.M{"password": 1}).One(&u)
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
	return coll.UpdateId(uid, bson.M{"$set": bson.M{"password": string(hash)}})
}

func FilterUsers(instituteName string) ([]types.SearchDoc, error) {
	sess := db.NewUserCollectionSession()
	defer sess.Close()
	coll := sess.Collection
	var result []types.SearchDoc
	err := coll.Find(bson.M{"institute": instituteName}).Select(bson.M{"_id": 1, "username": 1, "email": 1,
		"handle": 1, "picture": 1, "fullname": 1, "institute": 1}).All(&result)
	return result, err
}
