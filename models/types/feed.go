package types

import "github.com/globalsign/mgo/bson"

type FeedObject struct {
	UserName   string        `json:"username"`
	ID         bson.ObjectId `json:"user_id" bson:"_id"`
	FullName   string        `json:"fullname"`
	Picture    string        `json:"picture"`
	Submission Submission    `json:"submission"`
}
