package submission

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

type Submissions struct {
	Codechef   []CodechefSubmission   `bson:"codechef"`
	Codeforces []CodeforcesSubmission `bson:"codeforces"`
	Hackerrank []HackerrankSubmission `bson:"hackerrank"`
	Spoj       []SpojSubmission       `bson:"spoj"`
}

type CodechefSubmission struct {
	Name         string `bson:"name"`
	URL          string `bson:"url"`
	CreationDate string `bson:"creation_date"`
	Status       string `bson:"status"`
	Points       string `bson:"points"`
}

type SpojSubmission struct {
	Name         string `bson:"name"`
	URL          string `bson:"url"`
	CreationDate string `bson:"creation_date"`
	Status       string `bson:"status"`
	Language     string `bson:"language"`
}

type HackerrankSubmissions struct {
	Data  []HackerrankSubmission `json:"models" bson:"data"`
	Count int                    `json:"total" bson:"count"`
}

type HackerrankSubmission struct {
	URL          string `json:"url" bson:"url"`
	CreationDate string `json:"created_at" bson:"created_at"`
	Name         string `json:"name" bson:"name"`
}

// CodeforcesSubmission represents the single submission for codeforces
type CodeforcesSubmission struct {
	URL          string    `bson:"url"`
	CreationDate time.Time `bson:"created_at"`
	Name         string    `bson:"name"`
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
		problem := result.(map[string]interface{})["problem"].(map[string]interface{})
		submission := CodeforcesSubmission{}
		submission.URL = "http://codeforces.com/problemset/problem/" + strconv.Itoa(int(problem["contestId"].(float64))) + "/" + problem["index"].(string)
		submission.Name = problem["name"].(string)
		submission.CreationDate = time.Unix(int64(result.(map[string]interface{})["creationTimeSeconds"].(float64)), 0)
		sub.Data = append(sub.Data, submission)
	}
	return err
}
