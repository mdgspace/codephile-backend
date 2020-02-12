package types

import "time"

type Submission struct {
	Name         string    `json:"name" bson:"name"`
	URL          string    `json:"url" bson:"url"`
	CreationDate time.Time `json:"creation_date" bson:"creation_date"`
	Status       string    `json:"status" bson:"status"`
	Language     string    `json:"language" bson:"language"`
	Points       string    `json:"points" bson:"points"`
	Tags         []string  `json:"tags" bson:"tags"`
	Rating       int       `json:"rating" bson:"rating"`
}

type FeedObject struct {
	Username   string     `json:"username" bson:"username"`
	Submission Submission `json:"submission" bson:"submission"`
}
