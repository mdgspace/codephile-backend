package types

import (
	// "errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "github.com/mdg-iitr/Codephile/models/db"
	// "github.com/mdg-iitr/Codephile/models"
)

type Following struct {
	ID primitive.ObjectID `bson:"f_id" json:"f_id"`
}

type FollowingUser struct {
	ID       primitive.ObjectID 	`bson:"_id" json:"_id"`
	Username string        			`bson:"username" json:"username" schema:"username"`
	FullName string        			`bson:"fullname" json:"fullname" schema:"fullname"`
	Picture  string        			`bson:"picture" json:"picture"`
}

type WorldRankComparison struct {
	WorldRank1 string `bson:"rank1" json:"rank1"`
	WorldRank2 string `bson:"rank2" json:"rank2"`
}

type AllWorldRanks struct {
	CodechefWorldRanks   WorldRankComparison `bson:"codechef_ranks" json:"codechef_ranks"`
	CodeforcesWorldRanks WorldRankComparison `bson:"codeforces_ranks" json:"codeforces_ranks"`
	HackerrankWorldRanks WorldRankComparison `bson:"hackerrank_ranks" json:"hackerrank_ranks"`
	SpojWorldRanks       WorldRankComparison `bson:"spoj_ranks" json:"spoj_ranks"`
}
