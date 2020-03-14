package scripts

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mdg-iitr/Codephile/models/types"
	"log"
	"strconv"
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

func GetCodeforcesProfileInfo(handle string) types.ProfileInfo {
	var profile types.ProfileInfo
	url := "http://codeforces.com/api/user.info?handles=" + handle
	data := GetRequest(url)
	err := json.Unmarshal(data, &profile)
	if err != nil {
		log.Println(err.Error())
	}
	return profile
}

func GetCodeforcesGraphData(handle string) CodeforcesGraphPoints {
	var points CodeforcesGraphPoints
	url := "http://codeforces.com/api/user.rating?handle=" + handle
	data := GetRequest(url)
	err := json.Unmarshal(data, &points)
	if err != nil {
		log.Println(err.Error())
	}
	return points
}
func getCodeforcesSubmissionParts(handle string, afterIndex int) []types.Submission {
	url := "http://codeforces.com/api/user.status?handle=" + handle + "&from=" + strconv.Itoa(afterIndex) + "&count=50"
	fmt.Println(url)
	data := GetRequest(url)
	var codeforcesSubmission types.CodeforcesSubmissions
	err := json.Unmarshal(data, &codeforcesSubmission)
	if err != nil {
		log.Println(err.Error())
	}
	if codeforcesSubmission.Status != "OK" {
		log.Println("Codeforces submission could not be retrieved\n", string(data))
		return nil
	}
	submissions := make([]types.Submission, len(codeforcesSubmission.Result))
	for i, result := range codeforcesSubmission.Result {
		problem := result["problem"].(map[string]interface{})
		submissions[i].Status = result["verdict"].(string)
		submissions[i].Language = result["programmingLanguage"].(string)
		submissions[i].Name = problem["name"].(string)
		submissions[i].URL = "http://codeforces.com/problemset/problem/" + strconv.Itoa(int(problem["contestId"].(float64))) + "/" + problem["index"].(string)
		submissions[i].CreationDate = time.Unix(int64(result["creationTimeSeconds"].(float64)), 0)
		if problem["points"] != nil {
			submissions[i].Points = int(problem["points"].(float64))
		}
		if problem["rating"] != nil {
			submissions[i].Rating = int(problem["rating"].(float64))
		}
		for _, x := range problem["tags"].([]interface{}) {
			submissions[i].Tags = append(submissions[i].Tags, x.(string))
		}
	}
	return submissions
}

func GetCodeforcesSubmissions(handle string, after time.Time) []types.Submission {
	var oldestSubIndex, current int
	var oldestSubFound = false
	var subs []types.Submission
	//Fetch submission until oldest submission not found
	for !oldestSubFound {
		newSub := getCodeforcesSubmissionParts(handle, current+1)
		//Check for repetition of previous fetched submission
		if len(newSub) != 0 {
			for i, sub := range newSub {
				subs = append(subs, sub)
				oldestSubIndex = current + i + 1
				if sub.CreationDate.Equal(after) || sub.CreationDate.Before(after) {
					oldestSubFound = true
					break
				}
			}
			//50 submissions per page
			current += 50
		} else {
			break
		}
	}
	subs = subs[0:oldestSubIndex]
	return subs
}

func GetCodeforcesContests() CodeforcesContests {
	data := GetRequest("https://codeforces.com/api/contest.list?gym=false")
	var contests CodeforcesContests
	err := json.Unmarshal(data, &contests)
	if err != nil {
		log.Println(err.Error())
	}
	return contests
}
func CheckCodeforcesHandle(handle string) bool {
	data := GetRequest("http://codeforces.com/api/user.info?handles=" + handle)
	var i interface{}
	err := json.Unmarshal(data, &i)
	if err != nil {
		log.Println(err.Error())
	}
	return i.(map[string]interface{})["status"] != "FAILED"
}
