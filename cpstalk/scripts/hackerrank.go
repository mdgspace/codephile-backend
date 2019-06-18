package scripts

import "net/http"
import "log"
import "io/ioutil"
import "encoding/json"

type ProfileInfo struct {
	Name string
	UserName string
	JoinDate string
	School string
}

type Submissions struct {
	Data []Submission `json:"models"`
	Count int `json:"total"`
}

type Submission struct { 
	URL string `json:"url"`
	CreationDate string `json:"created_at"`
	Name string `json:"name"`
}

type GraphPoint struct {
	ContestName string 
	Date string
	Rating float64 
}

type Contests struct {
	Data []Contest `json:"models"`
	Count int `json:"total"`
}

type Contest struct {
	ContestName string `json:"name"`
	Rated bool `json:"rated"`
	EpochStart int64 `json:"epoch_starttime"`
	EpochEnd int64 `json:"epoch_endtime"`
	Archived bool `json:"archived"` 
}

func Get_Request(path string) []byte{
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

func GetProfileInfo(handle string) ProfileInfo {
	path := "https://www.hackerrank.com/rest/contests/master/hackers/" + handle + "/profile";
	byteValue := Get_Request(path)
	var JsonInterFace interface{}
    json.Unmarshal(byteValue, &JsonInterFace)
    Profile := JsonInterFace.(map[string]interface{})["model"].(map[string]interface{})
    Name := Profile["name"].(string)
    Date := Profile["created_at"].(string)
    UserName := Profile["username"].(string)
    School := Profile["school"].(string)
	return ProfileInfo {Name, UserName, Date, School}
}

func GetSubmissions(handle string) Submissions {
	path := "https://www.hackerrank.com/rest/hackers/" + handle + "/recent_challenges?limit=1000&response_version=v1"
	byteValue := Get_Request(path)
	var submissions Submissions
	json.Unmarshal(byteValue, &submissions)
	for i:=0; i<len(submissions.Data); i++ {
		submissions.Data[i].URL = "https://www.hackerrank.com" + submissions.Data[i].URL
	}
	return submissions
}

func GetContests() Contests {
	path := "https://www.hackerrank.com/rest/contests/upcoming?offset=0&limit=20&contest_slug=active"
	byteValue := Get_Request(path)
	var ContestsArray Contests
	json.Unmarshal(byteValue, &ContestsArray)
	return ContestsArray
}

func GetGraphData(handle string) []GraphPoint{
	path := "https://www.hackerrank.com/rest/hackers/" + handle +"/rating_histories_elo"
	byteValue := Get_Request(path)
	var JsonInterFace interface{}
    json.Unmarshal(byteValue, &JsonInterFace)

    m := JsonInterFace.(map[string]interface{})

    models := m["models"].([]interface{})
    events := models[0].(map[string]interface{})["events"].([]interface{})
    var Graph []GraphPoint
    for i := 0; i<len(events); i++ {
    	contest := events[i].(map[string]interface{})
    	name := contest["contest_name"].(string)
    	date := contest["date"].(string)
    	rating := contest["rating"].(float64)
    	Graph = append(Graph, GraphPoint{name, date, rating})
    }
    return Graph
}