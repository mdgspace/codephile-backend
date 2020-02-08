package types

import "time"

type FeedObject struct {
	UserName     string    `bson:"username"`
	Name         string    `bson:"name"`
	URL          string    `bson:"url"`
	CreationDate time.Time `bson:"creation_date"`
	Status       string    `bson:"status"`
	Language     string    `bson:"language"`
	Points       string    `bson:"points"`
	Tags         []string  `bson:"tags"`
	Rating       int       `bson:"rating"`
}

