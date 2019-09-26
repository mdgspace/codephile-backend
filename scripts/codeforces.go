package scripts

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mdg-iitr/Codephile/models/submission"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/mdg-iitr/Codephile/models/profile"
	"time"
)

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

func GetCodeforcesProfileInfo(handle string) profile.ProfileInfo {
	var profile profile.ProfileInfo
	url := "http://codeforces.com/api/user.info?handles=" + handle
	data := GetRequest(url)
	json.Unmarshal(data, &profile)
	return profile
}

func GetCodeforcesGraphData(handle string) CodeforcesGraphPoints {
	var points CodeforcesGraphPoints
	url := "http://codeforces.com/api/user.rating?handle=" + handle
	data := GetRequest(url)
	json.Unmarshal(data, &points)
	fmt.Println(points.Count)
	return points
}
func GetCodeforcesSubmissions(handle string) submission.CodeforcesSubmissions {
	url := "http://codeforces.com/api/user.status?handle=" + handle + "&from=1&count=10"
	data := GetRequest(url)
	var submissions submission.CodeforcesSubmissions
	json.Unmarshal(data, &submissions)
	return submissions
}
func GetCodeforcesContests() CodeforcesContests {
	data := GetRequest("https://codeforces.com/api/contest.list?gym=false")
	var contests CodeforcesContests
	json.Unmarshal(data, &contests)
	return contests
}
func CheckCodeforcesHandle(handle string) bool {
	resp, err := http.Get("http://codeforces.com/api/user.info?handles=" + handle)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err.Error())
	}
	data, _ := ioutil.ReadAll(resp.Body)
	var i interface{}
	err = json.Unmarshal(data, &i)
	if err != nil {
		log.Println(err.Error())
	}
	return i.(map[string]interface{})["status"] != "FAILED"
}
