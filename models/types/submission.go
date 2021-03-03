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

type CodechefSubmissions struct {
	Status string `json:"status"`
	Result Result2 `json:"result"`
}
type Content struct {
	ID          int    `json:"id"`
	ProblemCode string `json:"problemCode"`
	Language    string `json:"language"`
	Result      string `json:"result"`
	Username    string `json:"username"`
	Date        string `json:"date"`
}
type Data struct {
	Content []Content `json:"content"`
	Code    int       `json:"code"`
	Message string    `json:"message"`
}
type Result2 struct {
	Data Data `json:"data"`
}
