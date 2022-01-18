package leetcode

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/mdg-iitr/Codephile/models/types"
)

type Scrapper struct {
	Handle  string
	Context context.Context
}

func (s Scrapper) GetSubmissions(after time.Time) []types.LeetcodeSubmissions {
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
	jsonValue, _ := json.Marshal(jsonData)
	request, err := http.NewRequest("POST", "https://leetcode.com/graphql", bytes.NewBuffer(jsonValue))
	client := &http.Client{Timeout: time.Second * 10}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	var submissions []types.LeetcodeSubmissions
	json.Unmarshal(body, &submissions)
	fmt.Println(submissions)
	for i := 0; i < len(submissions); i++ {
		submissions[i].URL = "https://leetcode.com/problems/" + submissions[i].URL
	}
	return submissions
}
