package types

import (
	"time"
)

type Submission struct {
	Name         string    `json:"name" bson:"name"`
	URL          string    `json:"url" bson:"url"`
	CreationDate time.Time `json:"created_at" bson:"created_at"`
	Status       string    `json:"status" bson:"status"`
	Language     string    `json:"language" bson:"language"`
	Points       int       `json:"points" bson:"points"`
	Tags         []string  `json:"tags" bson:"tags"`
	Rating       int       `json:"rating" bson:"rating"`
}

type HackerrankSubmisson struct {
	Models []Submission `json:"models"`
}
type CodeforcesSubmissions struct {
	Status string                   `json:"status"`
	Result []map[string]interface{} `json:"result"`
}
