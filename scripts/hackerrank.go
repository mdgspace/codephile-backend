package scripts

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type HackerrankProfileInfo struct {
	Name     string
	UserName string
	JoinDate string
	School   string
}

type Submissions struct {
	Data  []HackerrankSubmission `json:"models"`
	Count int                    `json:"total"`
}

type HackerrankSubmission struct {
	URL          string `json:"url"`
	CreationDate string `json:"created_at"`
	Name         string `json:"name"`
}

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
	resp, err := http.Get(path)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	byteValue, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return byteValue
}

func GetHackerrankProfileInfo(handle string) HackerrankProfileInfo {
	path := "https://www.hackerrank.com/rest/contests/master/hackers/" + handle + "/profile";
	byteValue := GetRequest(path)
	var JsonInterFace interface{}
	json.Unmarshal(byteValue, &JsonInterFace)
	Profile := JsonInterFace.(map[string]interface{})["model"].(map[string]interface{})
	Name := Profile["name"].(string)
	Date := Profile["created_at"].(string)
	UserName := Profile["username"].(string)
	School := Profile["school"].(string)
	return HackerrankProfileInfo{Name, UserName, Date, School}
}

func GetHackerrankSubmissions(handle string) Submissions {
	path := "https://www.hackerrank.com/rest/hackers/" + handle + "/recent_challenges?limit=1000&response_version=v1"
	byteValue := GetRequest(path)
	var submissions Submissions
	json.Unmarshal(byteValue, &submissions)
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
