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

func (s Scrapper) GetSubmissions() []types.LeetcodeSubmissions {
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
				
            }
        `,
	}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		hub.CaptureException(err)
		log.Println(err.Error())
		return nil
	}
	request, err := http.NewRequest("POST", "https://leetcode.com/graphql", bytes.NewBuffer(jsonValue))
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
	var submissions []types.LeetcodeSubmissions
	json.Unmarshal(body, &submissions)
	fmt.Println(submissions)
	for i := 0; i < len(submissions); i++ {
		submissions[i].URL = "https://leetcode.com/problems/" + submissions[i].URL
	}
	return submissions
}
