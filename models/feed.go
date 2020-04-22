package models

import (
	"errors"
	"github.com/mdg-iitr/Codephile/models/db"
	"github.com/mdg-iitr/Codephile/models/types"
	"time"

	"github.com/globalsign/mgo/bson"
)

var ErrGeneric = errors.New("Feed is not absolutely correct")

func GetContestsFeed() (types.Result, error) {
	//contests stored are already sorted
	return contestsFromCache()
}

func GetAllFeed(uid bson.ObjectId) ([]types.FeedObject, error) {
	sess := db.NewUserCollectionSession()
	defer sess.Close()
	coll := sess.Collection
	var u types.User
	err := coll.FindId(uid).Select(bson.M{"followingUsers.f_id": 1}).One(&u)
	if err != nil {
		return nil, err
	}
	followingUID := make([]bson.ObjectId, 0, len(u.FollowingUsers))
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
	pipe := coll.Pipe([]bson.M{
		filter,
		project,
		unwind,
		sort,
	}, )
	var res []types.FeedObject
	err = pipe.All(&res)
	return res, err
}

func GetFeed(uid bson.ObjectId, before time.Time) ([]types.FeedObject, error) {
	sess := db.NewUserCollectionSession()
	defer sess.Close()
	coll := sess.Collection
	var u types.User
	err := coll.FindId(uid).Select(bson.M{"followingUsers.f_id": 1}).One(&u)
	if err != nil {
		return nil, err
	}
	followingUID := make([]bson.ObjectId, 0, len(u.FollowingUsers))
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

	pipe := coll.Pipe([]bson.M{
		filter,
		project,
		unwind,
		sort,
		limit,
	}, )

	var res []types.FeedObject
	err = pipe.All(&res)
	//fmt.Println(res)
	return res, err
}
