package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FeedObject struct {
	UserName   string        		`json:"username"`
	ID         primitive.ObjectID 	`json:"user_id" bson:"_id"`
	FullName   string        		`json:"fullname"`
	Picture    string        		`json:"picture"`
	Submission Submission    		`json:"submission"`
}
