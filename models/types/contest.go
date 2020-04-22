package types

import (
	"encoding/json"
)

type Ongoing struct {
	EndTime       string `json:"EndTime" bson:"EndTime"`
	Name          string `json:"Name" bson:"Name"`
	Platform      string `json:"Platform" bson:"Platform"`
	ChallengeType string `json:"challenge_type,omitempty" bson:"challenge_type,omitempty"`
	URL           string `json:"url" bson:"url"`
}

type Upcoming struct {
	Duration      string `json:"Duration" bson:"Duration"`
	EndTime       string `json:"EndTime" bson:"EndTime"`
	Name          string `json:"Name" bson:"Name"`
	Platform      string `json:"Platform" bson:"Platform"`
	StartTime     string `json:"StartTime" bson:"StartTime"`
	URL           string `json:"url" bson:"url"`
	ChallengeType string `json:"challenge_type,omitempty" bson:"challenge_type,omitempty"`
}

type Result struct {
	Ongoing   []Ongoing  `json:"ongoing" bson:"ongoing"`
	Timestamp string     `json:"timestamp" bson:"timestamp"`
	Upcoming  []Upcoming `json:"upcoming" bson:"upcoming"`
}

func (res Result) MarshalBinary() ([]byte, error) {
	return json.Marshal(res)
}

func (res *Result) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, res)
}
