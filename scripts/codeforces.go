package scripts

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

//CodeforcesProfileInfo represents the user profile for codeforces
type CodeforcesProfileInfo struct {
	Name     string
	UserName string
	School   string
	JoinDate time.Time
}

// CodeforcesSubmissions represents the submissions for codeforces
type CodeforcesSubmissions struct {
	Data  []CodeforcesSubmission
	Count int
}

// CodeforcesSubmission represents the single submissions for codeforces
type CodeforcesSubmission struct {
	URL          string
	CreationDate time.Time
	Name         string
}

// CodeforcesGraphPoint represents a single point for codeforces
type CodeforcesGraphPoint struct {
	ContestName string
	Date        time.Time
	Rating      float64
}

// CodeforcesGraphPoints represents the graph points for codeforces
type CodeforcesGraphPoints struct {
	Count  int
	Points []CodeforcesGraphPoint
}

//CodeforcesContests represents the codeforces contest
type CodeforcesContests struct {
	Data  []CodeforcesContest
	Count int
}
type CodeforcesContest struct {
	ContestName string `json:"name"`
	Rated       bool   `json:"rated"`
	EpochStart  int64  `json:"epoch_starttime"`
	EpochEnd    int64  `json:"epoch_endtime"`
	Archived    bool   `json:"archived"`
}

func Get_Request(path string) []byte {
	resp, err := http.Get(path)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	byteValue, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return byteValue
}

//UnmarshalJSON implements the unmarshaler interface for CodeforcesProfileInfo
func (data *CodeforcesProfileInfo) UnmarshalJSON(b []byte) error {
	var profile map[string]interface{}
	err := json.Unmarshal(b, &profile)
	if profile["status"] != "OK" {
		return errors.New("Bad Request")
	}
	result := profile["result"].([]interface{})[0].(map[string]interface{})
	if result["firstName"] != nil && result["lastName"] != nil {
		data.Name = result["firstName"].(string) + result["lastName"].(string)
	}
	data.UserName = result["handle"].(string)
	data.JoinDate = time.Unix(int64(result["registrationTimeSeconds"].(float64)), 0)
	if result["organization"] != nil {
		data.School = result["organization"].(string)
	}
	return err
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

// UnmarshalJSON implements the unmarshaler interface for CodeforcesGraphPoint
func (points *CodeforcesGraphPoints) UnmarshalJSON(b []byte) error {
	var data map[string]interface{}
	err := json.Unmarshal(b, &data)
	if data["status"] != "OK" {
		return errors.New("Bad Request")
	}
	results := data["result"].([]interface{})
	points.Count = len(results)
	for _, result := range results {
		point := CodeforcesGraphPoint{
			ContestName: result.(map[string]interface{})["contestName"].(string),
			Date:        time.Unix(int64(result.(map[string]interface{})["ratingUpdateTimeSeconds"].(float64)), 0),
			Rating:      result.(map[string]interface{})["newRating"].(float64),
		}
		points.Points = append(points.Points, point)
	}
	return err
}

// UnmarshalJSON implements the unmarshaler interface for CodeforcesContests
func (contests *CodeforcesContests) UnmarshalJSON(b []byte) error {
	var data map[string]interface{}
	err := json.Unmarshal(b, &data)
	if data["status"] != "OK" {
		return errors.New("Bad Request")
	}
	results := data["result"].([]interface{})[0:20]
	contests.Count = 20
	for _, result := range results {
		resultMap := result.(map[string]interface{})
		Contest := CodeforcesContest{
			ContestName: resultMap["name"].(string),
			Rated:       true,
			EpochStart:  int64(resultMap["startTimeSeconds"].(float64)),
		}
		Contest.EpochEnd = int64(resultMap["durationSeconds"].(float64)) + Contest.EpochStart
		phase := resultMap["phase"].(string)
		if phase == "FINISHED" {
			Contest.Archived = true
		}
		contests.Data = append(contests.Data, Contest)
	}
	return err
}

func GetCodeforcesProfileInfo(handle string) CodeforcesProfileInfo {
	var profile CodeforcesProfileInfo
	url := "http://codeforces.com/api/user.info?handles=" + handle
	data := Get_Request(url)
	json.Unmarshal(data, &profile)
	return profile
}

func GetCodeforcesGraphData(handle string) CodeforcesGraphPoints {
	var points CodeforcesGraphPoints
	url := "http://codeforces.com/api/user.rating?handle=" + handle
	data := Get_Request(url)
	json.Unmarshal(data, &points)
	fmt.Println(points.Count)
	return points
}
func GetCodeforcesSubmissions(handle string) CodeforcesSubmissions {
	url := "http://codeforces.com/api/user.status?handle=" + handle + "&from=1&count=10"
	data := Get_Request(url)
	var submissions CodeforcesSubmissions
	json.Unmarshal(data, &submissions)
	return submissions
}
func GetCodeforcesContests() CodeforcesContests {
	data := Get_Request("https://codeforces.com/api/contest.list?gym=false")
	var contests CodeforcesContests
	json.Unmarshal(data, &contests)
	return contests
}
