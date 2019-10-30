package submission

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

type Submissions struct {
	Codechef   []CodechefSubmission   `bson:"codechef" json:"codechef"`
	Codeforces []CodeforcesSubmission `bson:"codeforces" json:"codeforces"`
	Hackerrank []HackerrankSubmission `bson:"hackerrank" json:"hackerrank"`
	Spoj       []SpojSubmission       `bson:"spoj" json:"spoj"`
}

type CodechefSubmission struct {
	Name         string    `bson:"name" json:"name"`
	URL          string    `bson:"url" json:"url"`
	CreationDate time.Time `bson:"creation_date" json:"creation_date"`
	Status       string    `bson:"status" json:"status"`
	Points       string    `bson:"points" json:"points"`
	Tags         []string  `bson:"tags" json:"tags"`
	LanguageUsed string    `bson:"language" json:"language"`
}

type SpojSubmission struct {
	Name         string    `bson:"name" json:"name"`
	URL          string    `bson:"url" json:"url"`
	CreationDate time.Time `bson:"creation_date" json:"creation_date"`
	Status       string    `bson:"status" json:"status"`
	Language     string    `bson:"language" json:"language"`
	Points       int       `bson:"points" json:"points"`
	Tags         []string  `bson:"tags" json:"tags"`
}

type HackerrankSubmissions struct {
	Data  []HackerrankSubmission `json:"models" bson:"data"`
	Count int                    `json:"total" bson:"count"`
}

type HackerrankSubmission struct {
	URL          string    `json:"url" bson:"url"`
	CreationDate time.Time `json:"created_at" bson:"created_at"`
	Name         string    `json:"name" bson:"name"`
}

// CodeforcesSubmission represents the single submission for codeforces
type CodeforcesSubmission struct {
	URL          string    `bson:"url" json:"url"`
	CreationDate time.Time `bson:"created_at" json:"creation_date"`
	Name         string    `bson:"name" json:"name"`
	Status       string    `bson:"status" json:"status"`
	Points       int       `bson:"points" json:"points"`
	Rating       int       `bson:"rating" json:"rating"`
	Tags         []string  `bson:"tags" json:"tags"`
}

// CodeforcesSubmissions represents the submission for codeforces
type CodeforcesSubmissions struct {
	Data  []CodeforcesSubmission `bson:"data"`
	Count int                    `bson:"count"`
}

// UnmarshalJSON implements the unmarshaler interface for CodeforcesSubmissions
func (sub *CodeforcesSubmissions) UnmarshalJSON(b []byte) error {
	var data map[string]interface{}
	err := json.Unmarshal(b, &data)
	if data["status"] != "OK" {
		return errors.New("Bad Request")
	}
	results := data["result"].([]interface{})
	sub.Count = len(results)
	for _, result := range results {
		r := result.(map[string]interface{})
		problem := result.(map[string]interface{})["problem"].(map[string]interface{})
		submission := CodeforcesSubmission{}
		submission.URL = "http://codeforces.com/problemset/problem/" + strconv.Itoa(int(problem["contestId"].(float64))) + "/" + problem["index"].(string)
		submission.Name = problem["name"].(string)
		for _, x := range problem["tags"].([]interface{}) {
			submission.Tags = append(submission.Tags, x.(string))
		}
		if (problem["points"] != nil) {
			submission.Points = int(problem["points"].(float64))
		}
		if problem["rating"] != nil {
			submission.Rating = int(problem["rating"].(float64))
		}
		submission.Status = r["verdict"].(string)
		submission.CreationDate = time.Unix(int64(r["creationTimeSeconds"].(float64)), 0)
		sub.Data = append(sub.Data, submission)
	}
	return err
}
