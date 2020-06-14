package codeforces

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/mdg-iitr/Codephile/conf"
	"github.com/mdg-iitr/Codephile/models/types"
	"github.com/mdg-iitr/Codephile/scrappers/common"
	"log"
	"strconv"
	"time"
)

type Scrapper struct {
	Handle string
}

func (s Scrapper) GetProfileInfo() types.ProfileInfo {
	var profile types.ProfileInfo
	url := "http://codeforces.com/api/user.info?handles=" + s.Handle
	data := common.HitGetRequest(url)
	if data == nil {
		log.Println(errors.New("GetRequest failed. Please check connection status"))
		return types.ProfileInfo{}
	}
	err := json.Unmarshal(data, &profile)
	if err != nil {
		log.Println(err.Error())
		return types.ProfileInfo{}
	}
	return profile
}

// Calls the codeforces submission API and return the response in same format
func callCodeforcesAPI(handle string, afterIndex int) (types.CodeforcesSubmissions, error) {
	url := "http://codeforces.com/api/user.status?handle=" + handle + "&from=" + strconv.Itoa(afterIndex) + "&count=50"
	fmt.Println(url)
	data := common.HitGetRequest(url)
	if data == nil {
		return types.CodeforcesSubmissions{}, errors.New("GetRequest failed. Please check connection status")
	}
	var codeforcesSubmission types.CodeforcesSubmissions
	err := json.Unmarshal(data, &codeforcesSubmission)
	if err != nil {
		log.Println(err.Error())
		return types.CodeforcesSubmissions{}, err
	}
	return codeforcesSubmission, nil
}

//Get submissions of a user after an index.
//Returns an error if unsuccessful
//On receiving the error caller should return empty submission list
func getCodeforcesSubmissionParts(handle string, afterIndex int) ([]types.Submission, error) {
	codeforcesSubmission, err := callCodeforcesAPI(handle, afterIndex)
	if err != nil {
		return nil, err
	}
	if codeforcesSubmission.Status != "OK" {
		log.Println("Codeforces submission could not be retrieved. Retrying...")
		var newCodeforcesSub types.CodeforcesSubmissions
		for attempt := 1; attempt < 5; attempt++ {
			time.Sleep(time.Second * time.Duration(attempt))
			newCodeforcesSub, err = callCodeforcesAPI(handle, afterIndex)
			if err != nil {
				return nil, err
			}
			if newCodeforcesSub.Status == "OK" {
				codeforcesSubmission = newCodeforcesSub
				break
			}
		}
		if newCodeforcesSub.Status == "FAILED" {
			return nil, errors.New("codeforces API repeatedly returned FAILED")
		}
	}
	submissions := make([]types.Submission, len(codeforcesSubmission.Result))
	for i, result := range codeforcesSubmission.Result {
		problem := result["problem"].(map[string]interface{})
		var status string
		switch result["verdict"].(string) {
		case "FAILED":
			status = StatusWrongAnswer
		case "OK":
			status = StatusCorrect
		case "PARTIAL":
			status = StatusPartial
		case "COMPILATION_ERROR":
			status = StatusCompilationError
		case "RUNTIME_ERROR":
			status = StatusRuntimeError
		case "WRONG_ANSWER":
			status = StatusWrongAnswer
		case "TIME_LIMIT_EXCEEDED":
			status = StatusTimeLimitExceeded
		case "MEMORY_LIMIT_EXCEEDED":
			status = StatusMemoryLimitExceeded
		default:
			status = StatusWrongAnswer
		}
		submissions[i].Status = status
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
	return submissions, nil
}

func (s Scrapper) GetSubmissions(after time.Time) []types.Submission {
	var oldestSubIndex, current int
	var oldestSubFound = false
	var subs []types.Submission
	//Fetch submission until oldest submission not found
	for !oldestSubFound {
		newSub, err := getCodeforcesSubmissionParts(s.Handle, current+1)
		if err != nil {
			log.Println(err.Error())
			return nil
		}
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

func (s Scrapper) CheckHandle() bool {
	data := common.HitGetRequest("http://codeforces.com/api/user.info?handles=" + s.Handle)
	var i interface{}
	err := json.Unmarshal(data, &i)
	if err != nil {
		log.Println(err.Error())
	}
	return i.(map[string]interface{})["status"] != "FAILED"
}