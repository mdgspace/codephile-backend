package hackerrank

import (
	"encoding/json"
	"errors"
	"github.com/mdg-iitr/Codephile/models/types"
	"github.com/mdg-iitr/Codephile/scrappers/common"
	"log"
	"net/http"
	"time"
)

type Scrapper struct {
	Handle string
}

func (s Scrapper) GetProfileInfo() types.ProfileInfo {
	path := "https://www.hackerrank.com/rest/contests/master/hackers/" + s.Handle + "/profile";
	byteValue := common.HitGetRequest(path)
	if byteValue == nil {
		log.Println(errors.New("GetRequest failed. Please check connection status"))
		return types.ProfileInfo{}
	}
	var JsonInterFace interface{}
	err := json.Unmarshal(byteValue, &JsonInterFace)
	if err != nil {
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
	path := "https://www.hackerrank.com/rest/hackers/" + s.Handle + "/recent_challenges?limit=1000&response_version=v1"
	byteValue := common.HitGetRequest(path)
	var data types.HackerrankSubmisson
	err := json.Unmarshal(byteValue, &data)
	submissions := data.Models
	if err != nil {
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
		submissions[i].URL = "https://www.hackerrank.com" + submissions[i].URL
	}
	return submissions
}

func (s Scrapper) CheckHandle() bool {
	resp, err := http.Get("https://www.hackerrank.com/rest/contests/master/hackers/" + s.Handle + "/profile")
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return resp.StatusCode != http.StatusNotFound
}
