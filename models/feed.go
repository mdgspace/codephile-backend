package models

import (
	"context"
	"github.com/mdg-iitr/Codephile/models/db"
	"github.com/mdg-iitr/Codephile/models/types"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetContestsFeed() (types.Result, error) {
	//contests stored are already sorted
	return contestsFromCache()
}

func GetAllFeed(uid primitive.ObjectID) ([]types.FeedObject, error) {
	sess := db.NewUserCollectionSession()
	defer sess.Close()
	coll := sess.Collection
	var u types.User
	err := coll.FindOne(context.TODO(), bson.M{"_id": uid}, options.FindOne().SetProjection(bson.M{"followingUsers.f_id": 1})).Decode(&u)
	if err != nil {
		return nil, err
	}
	followingUID := make([]primitive.ObjectID, 0, len(u.FollowingUsers))
	for _, f := range u.FollowingUsers {
		followingUID = append(followingUID, f.ID)
	}
	filter := bson.M{
		"$match": bson.M{
			"_id": bson.M{
				"$in": followingUID,
			},
		},
	}
	project := bson.M{
		"$project": bson.M{
			"_id":        1,
			"username":   1,
			"picture":    1,
			"fullname":   1,
			"submission": "$submissions",
		},
	}
	unwind := bson.M{
		"$unwind": "$submission",
	}
	sort := bson.M{
		"$sort": bson.M{
			"submission.created_at": -1,
		},
	}
	pipe, err := coll.Aggregate(context.TODO(), []bson.M{
		filter,
		project,
		unwind,
		sort,
	}, )
	var res []types.FeedObject
	err = pipe.All(context.TODO(), &res)
	return res, err
}

func GetFeed(uid primitive.ObjectID, before time.Time) ([]types.FeedObject, error) {
	sess := db.NewUserCollectionSession()
	defer sess.Close()
	coll := sess.Collection
	var u types.User
	err := coll.FindOne(context.TODO(), bson.M{"_id": uid}, options.FindOne().SetProjection(bson.M{"followingUsers.f_id": 1})).Decode(&u)
	if err != nil {
		return nil, err
	}
	followingUID := make([]primitive.ObjectID, 0, len(u.FollowingUsers))
	for _, f := range u.FollowingUsers {
		followingUID = append(followingUID, f.ID)
	}
	filter := bson.M{
		"$match": bson.M{
			"_id": bson.M{
				"$in": followingUID,
			},
		},
	}
	project := bson.M{
		"$project": bson.M{
			"_id":      1,
			"username": 1,
			"picture":  1,
			"fullname": 1,
			"submission": bson.M{"$filter": bson.M{"input": "$submissions",
				"as":   "sub",
				"cond": bson.M{"$lt": []interface{}{"$$sub.created_at", before}},
			}},
		},
	}
	unwind := bson.M{
		"$unwind": "$submission",
	}
	sort := bson.M{
		"$sort": bson.M{
			"submission.created_at": -1,
		},
	}
	limit := bson.M{
		"$limit": 100,
	}

	pipe, err := coll.Aggregate(context.TODO(), []bson.M{
		filter,
		project,
		unwind,
		sort,
		limit,
	}, )

	var res []types.FeedObject
	err = pipe.All(context.TODO(), &res)
	//fmt.Println(res)
	return res, err
}
