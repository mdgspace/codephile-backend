package hackerrank

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/getsentry/sentry-go"
	. "github.com/mdg-iitr/Codephile/conf"
	"github.com/mdg-iitr/Codephile/models/types"
	"github.com/mdg-iitr/Codephile/scrappers/common"
	"log"
	"net/http"
	"time"
)

type Scrapper struct {
	Handle string
	Context context.Context
}

func (s Scrapper) GetProfileInfo() types.ProfileInfo {
	hub := sentry.GetHubFromContext(s.Context)
	if hub == nil {
		hub = sentry.CurrentHub()
	}
	path := "https://www.hackerrank.com/rest/contests/master/hackers/" + s.Handle + "/profile";
	byteValue, _ := common.HitGetRequest(path)
	if byteValue == nil {
		err := errors.New("GetRequest failed. Please check connection status")
		log.Println(err)
		hub.CaptureException(err)
		return types.ProfileInfo{}
	}
	var JsonInterFace interface{}
	err := json.Unmarshal(byteValue, &JsonInterFace)
	if err != nil {
		hub.AddBreadcrumb(&sentry.Breadcrumb{
			Category:  "JSON parse error",
			Message:   string(byteValue),
		}, nil)
		hub.CaptureException(err)
		log.Println(err.Error())
		return types.ProfileInfo{}
	}
	Profile := JsonInterFace.(map[string]interface{})["model"].(map[string]interface{})
	Name := Profile["name"].(string)
	// Date := Profile["created_at"].(string)
	UserName := Profile["username"].(string)
	School := Profile["school"].(string)
	return types.ProfileInfo{Name: Name, UserName: UserName, School: School}
}

func (s Scrapper) GetSubmissions(after time.Time) []types.Submission {
	hub := sentry.GetHubFromContext(s.Context)
	if hub == nil {
		hub = sentry.CurrentHub()
	}
	path := "https://www.hackerrank.com/rest/hackers/" + s.Handle + "/recent_challenges?limit=1000&response_version=v1"
	byteValue, _ := common.HitGetRequest(path)
	var data types.HackerrankSubmisson
	err := json.Unmarshal(byteValue, &data)
	submissions := data.Models
	if err != nil {
		hub.AddBreadcrumb(&sentry.Breadcrumb{
			Category:  "JSON parse error",
			Message:   string(byteValue),
		}, nil)
		hub.CaptureException(err)
		log.Println(err.Error())
		return nil
	}
	var oldestSubIndex int;
	if after.IsZero() {
		oldestSubIndex = len(submissions)
	} else {
		for i, sub := range submissions {
			if sub.CreationDate.Equal(after) || sub.CreationDate.Before(after) {
				oldestSubIndex = i
				break
			}
		}
	}
	submissions = submissions[0:oldestSubIndex]
	for i := 0; i < len(submissions); i++ {
		submissions[i].Points = 100
		submissions[i].Status = StatusCorrect
		submissions[i].URL = "https://www.hackerrank.com" + submissions[i].URL
	}
	return submissions
}

func (s Scrapper) CheckHandle() bool {
	hub := sentry.GetHubFromContext(s.Context)
	if hub == nil {
		hub = sentry.CurrentHub()
	}
	resp, err := http.Get("https://www.hackerrank.com/rest/contests/master/hackers/" + s.Handle + "/profile")
	if err != nil {
		log.Println(err.Error())
		hub.CaptureException(err)
		return false
	}
	return resp.StatusCode != http.StatusNotFound
}
