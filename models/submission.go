package models

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	. "github.com/mdg-iitr/Codephile/conf"
	. "github.com/mdg-iitr/Codephile/errors"
	"github.com/mdg-iitr/Codephile/models/db"
	"github.com/mdg-iitr/Codephile/models/types"
	"github.com/mdg-iitr/Codephile/scripts"
	"log"
	"time"
)

//Returns HandleNotFoundError/UserNotFoundError/error
func AddSubmissions(uid bson.ObjectId, site string) error {
	sess := db.NewUserCollectionSession()
	defer sess.Close()
	coll := sess.Collection
	var user types.User
	err := coll.FindId(uid).Select(bson.M{"handle": 1, "lastfetched": 1}).One(&user)
	if err != nil {
		//handle the error (Invalid user)
		return UserNotFoundError
	}
	var handle string
	var addSubmissions []types.Submission
	switch site {
	case CODECHEF:
		handle = user.Handle.Codechef
		if handle == "" {
			return HandleNotFoundError
		}
		//TODO: Return errors from scripts
		addSubmissions = scripts.GetCodechefSubmissions(handle, user.Last.Codechef)
		if len(addSubmissions) != 0 {
			user.Last.Codechef = addSubmissions[0].CreationDate
		}
	case CODEFORCES:
		handle = user.Handle.Codeforces
		if handle == "" {
			return HandleNotFoundError
		}
		//TODO: Return errors from scripts
		addSubmissions = scripts.GetCodeforcesSubmissions(handle, user.Last.Codeforces)
		if len(addSubmissions) != 0 {
			user.Last.Codeforces = addSubmissions[0].CreationDate
		}
	case SPOJ:
		handle = user.Handle.Spoj
		if handle == "" {
			return HandleNotFoundError
		}
		//TODO: Return errors from scripts
		addSubmissions = scripts.GetSpojSubmissions(handle, user.Last.Spoj)
		if len(addSubmissions) != 0 {
			user.Last.Spoj = addSubmissions[0].CreationDate
		}
	case HACKERRANK:
		handle = user.Handle.Hackerrank
		if handle == "" {
			return HandleNotFoundError
		}
		//TODO: Return errors from scripts
		addSubmissions = scripts.GetHackerrankSubmissions(handle, user.Last.Hackerrank)
		if len(addSubmissions) != 0 {
			user.Last.Hackerrank = addSubmissions[0].CreationDate
		}
	}
	change := bson.M{
		"$push": bson.M{
			"submissions": bson.M{
				"$each": addSubmissions,
				"$sort": bson.M{"created_at": -1},
			}},
		"$set": bson.M{"lastfetched": user.Last}}
	err = coll.UpdateId(user.ID, change)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func GetSubmissions(ID bson.ObjectId, before time.Time) ([]types.Submission, error) {
	sess := db.NewUserCollectionSession()
	defer sess.Close()
	coll := sess.Collection
	match := bson.M{
		"$match": bson.M{
			"_id": ID,
		},
	}
	project := bson.M{
		"$project": bson.M{
			"_id": 0,
			"submission": bson.M{"$filter": bson.M{
				"input": "$submissions",
				"as":    "sub",
				"cond":  bson.M{"$lt": []interface{}{"$$sub.created_at", before}},
			},
			},
		},
	}
	unwind := bson.M{
		"$unwind": "$submission",
	}
	limit := bson.M{
		"$limit": 100,
	}
	group := bson.M{"$group": bson.M{"_id": ID, "submissions": bson.M{"$push": "$submission"}}}
	pipe := coll.Pipe([]bson.M{
		match,
		project,
		unwind,
		limit,
		group,
	}, )
	var res types.User
	err := pipe.One(&res)
	return res.Submissions, err
}
func GetAllSubmissions(ID bson.ObjectId) ([]types.Submission, error) {
	coll := db.NewUserCollectionSession()
	defer coll.Close()
	var user types.User
	err := coll.Collection.FindId(ID).Select(bson.M{"submissions": 1}).One(&user)
	if err != nil {
		return nil, err
	}
	return user.Submissions, nil
}

//TODO: Return proper errors in FilterSubmission
func FilterSubmission(uid bson.ObjectId, status string, tag string, site string) ([]map[string]interface{}, error) {
	c := db.NewUserCollectionSession()
	defer c.Close()
	fmt.Println(status)
	match1 := bson.M{
		"$match": bson.M{
			"_id": uid,
		},
	}
	unwind := bson.M{
		"$unwind": "$submission." + site,
	}
	match2 := bson.M{
		"$match": bson.M{"submission." + site + ".status": status},
	}
	project := bson.M{
		"$project": bson.M{
			"_id":                0,
			"submission." + site: 1,
		},
	}
	all := []bson.M{match1, unwind, match2, project}
	pipe := c.Collection.Pipe(all)

	var result map[string]interface{}
	iter := pipe.Iter()
	var final []map[string]interface{}
	for iter.Next(&result) {
		final = append(final, result["submission"].(map[string]interface{})[site].(map[string]interface{}))
	}
	return final, nil
}
