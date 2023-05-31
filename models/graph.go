package models

import (
	"context"

	"github.com/mdg-iitr/Codephile/conf"
	"github.com/mdg-iitr/Codephile/models/db"
	"github.com/mdg-iitr/Codephile/models/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetActivityGraph(uid primitive.ObjectID) (types.ActivityGraph, error) {
	sess := db.NewUserCollectionSession()
	defer sess.Close()
	coll := sess.Collection

	match := bson.M{"$match": bson.M{"_id": uid}}
	project := bson.M{"$project": bson.M{"submission": "$submissions"}}
	unwind := bson.M{"$unwind": "$submission"}
	group := bson.M{
		"$group": bson.M{
			"_id": bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$submission.created_at"}},
			//date: {$min: "$submission.created_at"},
			"correct": bson.M{
				"$sum": bson.M{
					"$cond": []interface{}{bson.M{"$eq": []string{"$submission.status", conf.StatusCorrect}}, 1, 0},
				},
			},
			"total": bson.M{"$sum": 1},
		}}
	pipe, err := coll.Aggregate(context.TODO(), []bson.M{
		match,
		project,
		unwind,
		group,
	})
	var res types.ActivityGraph
	err = pipe.All(context.TODO(), &res)
	return res, err
}

func GetStatusCounts(uid primitive.ObjectID) (types.StatusCounts, error) {
	sess := db.NewUserCollectionSession()
	defer sess.Close()
	coll := sess.Collection

	match := bson.M{"$match": bson.M{"_id": uid}}
	project := bson.M{
		"$project": bson.M{
			"ac_count":  GetStatusQuery(conf.StatusCorrect),
			"wa_count":  GetStatusQuery(conf.StatusWrongAnswer),
			"ce_count":  GetStatusQuery(conf.StatusCompilationError),
			"re_count":  GetStatusQuery(conf.StatusRuntimeError),
			"tle_count": GetStatusQuery(conf.StatusTimeLimitExceeded),
			"mle_count": GetStatusQuery(conf.StatusMemoryLimitExceeded),
			"ptl_count": GetStatusQuery(conf.StatusPartial),
		}}
	pipe, err := coll.Aggregate(context.TODO(), []bson.M{
		match,
		project,
	})

	var statusCounts types.StatusCounts
	err = pipe.Decode(&statusCounts)
	return statusCounts, err
}

func GetStatusQuery(status string) bson.M {
	return bson.M{
		"$size": bson.M{
			"$filter": bson.M{
				"input": "$submissions",
				"as":    "sub",
				"cond":  bson.M{"$eq": []string{"$$sub.status", status}},
			},
		},
	}
}
