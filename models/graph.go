package models

import (
	"github.com/globalsign/mgo/bson"
	"github.com/mdg-iitr/Codephile/conf"
	"github.com/mdg-iitr/Codephile/models/db"
	"github.com/mdg-iitr/Codephile/models/types"
)

func GetActivityGraph(uid bson.ObjectId) (types.ActivityGraph, error) {
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
	pipe := coll.Pipe([]bson.M{
		match,
		project,
		unwind,
		group,
	})
	var res types.ActivityGraph
	err := pipe.All(&res)
	return res, err
}
