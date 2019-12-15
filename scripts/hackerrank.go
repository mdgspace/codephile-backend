package scripts

import (
	"encoding/json"
	"github.com/mdg-iitr/Codephile/models/profile"
	"github.com/mdg-iitr/Codephile/models/submission"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type HackerrankGraphPoint struct {
	ContestName string
	Date        string
	Rating      float64
}

type Contests struct {
	Data  []HackerrankContest `json:"models"`
	Count int                 `json:"total"`
}

type HackerrankContest struct {
	ContestName string `json:"name"`
	Rated       bool   `json:"rated"`
	EpochStart  int64  `json:"epoch_starttime"`
	EpochEnd    int64  `json:"epoch_endtime"`
	Archived    bool   `json:"archived"`
}

func GetRequest(path string) []byte {
	client := http.Client{Timeout: time.Second * 10}
	resp, err := client.Get(path)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer resp.Body.Close()
	byteValue, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	return byteValue
}

func GetHackerrankProfileInfo(handle string) profile.ProfileInfo {
	path := "https://www.hackerrank.com/rest/contests/master/hackers/" + handle + "/profile";
	byteValue := GetRequest(path)
	var JsonInterFace interface{}
	json.Unmarshal(byteValue, &JsonInterFace)
	Profile := JsonInterFace.(map[string]interface{})["model"].(map[string]interface{})
	Name := Profile["name"].(string)
	// Date := Profile["created_at"].(string)
	UserName := Profile["username"].(string)
	School := Profile["school"].(string)
	return profile.ProfileInfo{Name, UserName, School, "", ""}
}

func GetHackerrankSubmissions(handle string, after time.Time) submission.HackerrankSubmissions {
	path := "https://www.hackerrank.com/rest/hackers/" + handle + "/recent_challenges?limit=1000&response_version=v1"
	byteValue := GetRequest(path)
	var submissions submission.HackerrankSubmissions
	json.Unmarshal(byteValue, &submissions)
	var oldestSubIndex int;
	if after.IsZero() {
		oldestSubIndex = submissions.Count
	} else {
		for i, sub := range submissions.Data {
			if sub.CreationDate.Equal(after) || sub.CreationDate.Before(after) {
				oldestSubIndex = i
				break
			}
		}
	}
	submissions.Data = submissions.Data[0:oldestSubIndex]
	submissions.Count = oldestSubIndex
	for i := 0; i < len(submissions.Data); i++ {
		submissions.Data[i].URL = "https://www.hackerrank.com" + submissions.Data[i].URL
	}
	return submissions
}

func GetHackerrankContests() Contests {
	path := "https://www.hackerrank.com/rest/contests/upcoming?offset=0&limit=20&contest_slug=active"
	byteValue := GetRequest(path)
	var ContestsArray Contests
	json.Unmarshal(byteValue, &ContestsArray)
	return ContestsArray
}

func GetHackerrankGraphData(handle string) []HackerrankGraphPoint {
	path := "https://www.hackerrank.com/rest/hackers/" + handle + "/rating_histories_elo"
	byteValue := GetRequest(path)
	var JsonInterFace interface{}
	json.Unmarshal(byteValue, &JsonInterFace)

	m := JsonInterFace.(map[string]interface{})

	models := m["models"].([]interface{})
	events := models[0].(map[string]interface{})["events"].([]interface{})
	var Graph []HackerrankGraphPoint
	for i := 0; i < len(events); i++ {
		contest := events[i].(map[string]interface{})
		name := contest["contest_name"].(string)
		date := contest["date"].(string)
		rating := contest["rating"].(float64)
		Graph = append(Graph, HackerrankGraphPoint{name, date, rating})
	}
	return Graph
}
func CheckHackerrankHandle(handle string) bool {
	resp, err := http.Get("https://www.hackerrank.com/rest/contests/master/hackers/" + handle + "/profile")
	if err != nil {
		log.Println(err.Error())
	}
	return resp.StatusCode != http.StatusNotFound
}
