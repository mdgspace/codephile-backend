package models

import(
	  "log"
	  "encoding/json"
)

var FinalResult S

type S struct {
	Result Result `json:"result"`
}

type Ongoing struct {
	EndTime       string `json:"EndTime"`
	Name          string `json:"Name"`
	Platform      string `json:"Platform"`
	ChallengeType string `json:"challenge_type,omitempty"`
	URL           string `json:"url"`
}

type Upcoming struct {
	Duration      string `json:"Duration"`
	EndTime       string `json:"EndTime"`
	Name          string `json:"Name"`
	Platform      string `json:"Platform"`
	StartTime     string `json:"StartTime"`
	URL           string `json:"url"`
	ChallengeType string `json:"challenge_type,omitempty"`
}

type Result struct {
	Ongoing   []Ongoing  `json:"ongoing"`
	Timestamp string     `json:"timestamp"`
	Upcoming  []Upcoming `json:"upcoming"`
}

func ParseContests(data []byte) {
	error := json.Unmarshal(data,&FinalResult)

	if error!=nil {
		log.Println("Error")
		log.Fatal(error)
	}
	
}

func ReturnContests() S{
	 return FinalResult
}


