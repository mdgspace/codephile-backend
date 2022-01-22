package leetcode

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/mdg-iitr/Codephile/models/types"
)

type Scrapper struct {
	Handle  string
	Context context.Context
}

func (s Scrapper) GetLeetcodesubmissions() []types.Submission {
	hub := sentry.GetHubFromContext(s.Context)
	if hub == nil {
		hub = sentry.CurrentHub()
	}
	jsonData := map[string]string{
		"query": `
            { 
				recentSubmissionList(username: "` + s.Handle + `"){
					title
					titleSlug
				    timestamp
				}	
            }`,
	}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		hub.CaptureException(err)
		log.Println(err.Error())
		return nil
	}
	request, err := http.NewRequest("POST", "https://leetcode.com/graphql", bytes.NewBuffer(jsonValue))
	if err != nil {
		hub.CaptureException(err)
		log.Println(err.Error())
		return nil
	}
	client := &http.Client{Timeout: time.Second * 10}
	response, err := client.Do(request)
	if err != nil {
		hub.CaptureException(err)
		log.Println(err.Error())
		return nil
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
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
