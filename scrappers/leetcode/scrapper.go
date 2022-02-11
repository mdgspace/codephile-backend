package leetcode

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/mdg-iitr/Codephile/models/types"
)

type Scrapper struct {
	Handle  string
	Context context.Context
}

func leetcodeGraphQLRequest(query string) ([]byte, error) {
	jsonData := map[string]string{
		"query": query,
	}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post("https://leetcode.com/graphql", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	responseValue, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return responseValue, err
}

func (s Scrapper) GetProfileInfo() types.ProfileInfo {
	hub := sentry.GetHubFromContext(s.Context)
	if hub == nil {
		hub = sentry.CurrentHub()
	}

	query := `
		{
			matchedUser(username: "` + s.Handle + `") {
				username
				profile {
					realName
					school
					ranking
				}
				submitStats {
					acSubmissionNum {
						submissions
					}
					totalSubmissionNum {
						submissions
					}
				}
			}
		}
	`
	responseData, err := leetcodeGraphQLRequest(query)
	if err != nil {
		log.Println(err.Error())
		hub.CaptureException(err)
		return types.ProfileInfo{}
	}
	var responseValue types.GraphQLResponse
	err = json.Unmarshal(responseData, &responseValue)
	if err != nil {
		log.Println(err.Error())
		hub.CaptureException(err)
		return types.ProfileInfo{}
	}
	matchedUser := responseValue.Data.MatchedUser
	profile := matchedUser.Profile
	submitStats := matchedUser.SubmitStats
	accuracy := submitStats.AcSubmissionNum[0].Submissions / math.Max(1, submitStats.TotalSubmissionNum[0].Submissions) * 100
	return types.ProfileInfo{
		Name:      profile.RealName,
		UserName:  matchedUser.Username,
		School:    profile.School,
		WorldRank: fmt.Sprintf("%.0f", profile.Ranking),
		Accuracy:  fmt.Sprintf("%.2f", accuracy),
	}
}

func (s Scrapper) CheckHandle() (bool, error) {
	hub := sentry.GetHubFromContext(s.Context)
	if hub == nil {
		hub = sentry.CurrentHub()
	}

	query := `
		{
			matchedUser(username: "` + s.Handle + `") {
				username
			}
		}
	`
	responseData, err := leetcodeGraphQLRequest(query)
	if err != nil {
		log.Println(err.Error())
		hub.CaptureException(err)
		return false, err
	}
	var responseValue map[string]interface{}
	err = json.Unmarshal(responseData, &responseValue)
	if err != nil {
		log.Println(err.Error())
		hub.CaptureException(err)
		return false, err
	}
	matchedUser := responseValue["data"].(map[string]interface{})["matchedUser"]
	return matchedUser != nil, err
}

func (s Scrapper) GetLeetcodesubmissions() []types.Submission {
	hub := sentry.GetHubFromContext(s.Context)
	if hub == nil {
		hub = sentry.CurrentHub()
	}

	query := `
            { 
				recentSubmissionList(username: "` + s.Handle + `"){
					title
					titleSlug
				    timestamp
				}	
            }`

	body, err := leetcodeGraphQLRequest(query)
	if err != nil {
		hub.CaptureException(err)
		log.Println(err.Error())
		return nil
	}
	var Leetcodesubmissions []types.LeetcodeSubmissions
	err1 := json.Unmarshal(body, &Leetcodesubmissions)
	if err1 != nil {
		hub.CaptureException(err1)
		log.Println(err1.Error())
		return nil
	}
	fmt.Println(Leetcodesubmissions)

	submissions := make([]types.Submission, len(Leetcodesubmissions))
	for i, result := range Leetcodesubmissions {
		submissions[i].Name = result.Title
		submissions[i].URL = "https://leetcode.com/problems/" + result.URL
		t, err := time.Parse("2006-01-02 15:04:05", result.TimeSTamp)
		if err != nil {
			hub.CaptureException(err)
		}
		submissions[i].CreationDate = t

	}
	return submissions
}
