package models  

import (
	// "errors"
	"github.com/globalsign/mgo/bson"
	// "github.com/mdg-iitr/Codephile/models/db"
	// "github.com/mdg-iitr/Codephile/models"
)

type Following struct{
	CodephileHandle string          `bson:"codephile_handle" json:"codephile_handle"`
	ID              bson.ObjectId   `bson:"f_id" json:"f_id"`
}

type WorldRankComparison struct{
	WorldRank1   string             `bson:"rank1" json:"rank1"`
	WorldRank2   string             `bson:"rank2" json:"rank2"`
} 

type AllWorldRanks struct {
	CodechefWorldRanks      WorldRankComparison    `bson:"codechef_ranks" json:"codechef_ranks"`
	CodeforcesWorldRanks    WorldRankComparison    `bson:"codeforces_ranks" json:"codeforces_ranks"`
	HackerrankWorldRanks    WorldRankComparison    `bson:"hackerrank_ranks" json:"hackerrank_ranks"`
	SpojWorldRanks          WorldRankComparison    `bson:"spoj_ranks" json:"spoj_ranks"`
}


