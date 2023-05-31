package models

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	// "github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	. "github.com/mdg-iitr/Codephile/conf"
	. "github.com/mdg-iitr/Codephile/errors"
	"github.com/mdg-iitr/Codephile/models/db"
	"github.com/mdg-iitr/Codephile/models/types"
	"github.com/mdg-iitr/Codephile/scrappers"
)

// Fetches Submissions which are made after the lastFetched time, and
// adds that to the database.
//Returns HandleNotFoundError/UserNotFoundError/error
func AddSubmissions(uid primitive.ObjectID, site string, ctx context.Context) error {
	if !IsSiteValid(site) {
		return errors.New("site invalid")
	}
	sess := db.NewUserCollectionSession()
	defer sess.Close()
	coll := sess.Collection
	var result map[string]interface{}
	err := coll.FindOne(context.TODO(), bson.M{"_id": uid}, options.FindOne().SetProjection(bson.M{"handle": 1, "lastfetched": 1})).Decode(&result)
	if err != nil {
		//handle the error (Invalid user)
		return UserNotFoundError
	}
	var addSubmissions []types.Submission
	lastFetched := result["lastfetched"].(map[string]interface{})[site].(time.Time)
	handle := result["handle"].(map[string]interface{})[site].(string)
	scrapper, err := scrappers.NewScrapper(site, handle, ctx)
	if err != nil {
		return err
	}
	addSubmissions = scrapper.GetSubmissions(lastFetched)
	if len(addSubmissions) != 0 {
		lastFetched = addSubmissions[0].CreationDate
	}

	change := bson.M{
		"$push": bson.M{
			"submissions": bson.M{
				"$each": addSubmissions,
				"$sort": bson.M{"created_at": -1},
			}},
		"$set": bson.M{"lastfetched." + site: lastFetched}}
	_, err = coll.UpdateByID(context.TODO(), uid, change)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func DeleteSubmissions(uid primitive.ObjectID, site string) error {
	sess := db.NewUserCollectionSession()
	defer sess.Close()
	coll := sess.Collection

	var resetTime time.Time
	_, err := coll.UpdateByID(context.TODO(), uid, bson.M{
		"$pull": bson.M{
			"submissions": bson.M{
				"url": bson.M{
					"$regex": primitive.Regex{
						Pattern: "^" + GetRegexSite(site)},
				}},
		},
		"$set": bson.M{"lastfetched." + site: resetTime},
	})

	return err
}

func GetSubmissions(ID primitive.ObjectID, before time.Time) ([]types.Submission, error) {
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
	pipe, err := coll.Aggregate(context.TODO(), []bson.M{
		match,
		project,
		unwind,
		limit,
		group,
	})
	var res types.User
	err = pipe.Decode(&res)
	return res.Submissions, err
}

func GetAllSubmissions(ID primitive.ObjectID) ([]types.Submission, error) {
	coll := db.NewUserCollectionSession()
	defer coll.Close()
	var user types.User
	err := coll.Collection.FindOne(context.TODO(), bson.M{"_id": ID}, options.FindOne().SetProjection(bson.M{"submissions": 1})).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user.Submissions, nil
}

//TODO: Return proper errors in FilterSubmission
func FilterSubmission(uid primitive.ObjectID, status string, tag string, site string) ([]map[string]interface{}, error) {
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
	cursor, _ := c.Collection.Aggregate(context.TODO(), all)

	var result map[string]interface{}
	var final []map[string]interface{}
	for cursor.Next(context.TODO()) {
		if err := cursor.Decode(&result); err != nil {
			panic(err)
		}
		final = append(final, result["submission"].(map[string]interface{})[site].(map[string]interface{}))
	}
	return final, nil
}
