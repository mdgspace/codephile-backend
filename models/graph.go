package models

import (
	"github.com/globalsign/mgo/bson"
	"github.com/mdg-iitr/Codephile/models/types"
)

func GetActivityGraph(uid bson.ObjectId) (types.ActivityGraph, error) {
	return types.ActivityGraph{}, nil

}
