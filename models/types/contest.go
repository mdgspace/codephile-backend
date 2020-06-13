package types

import (
	"encoding/json"
	"time"
)

type Ongoing struct {
	EndTime       ContestTime `json:"EndTime" bson:"EndTime"`
	Name          string      `json:"Name" bson:"Name"`
	Platform      string      `json:"Platform" bson:"Platform"`
	ChallengeType string      `json:"challenge_type,omitempty" bson:"challenge_type,omitempty"`
	URL           string      `json:"url" bson:"url"`
}

type Upcoming struct {
	Duration      string      `json:"Duration" bson:"Duration"`
	EndTime       ContestTime `json:"EndTime" bson:"EndTime"`
	Name          string      `json:"Name" bson:"Name"`
	Platform      string      `json:"Platform" bson:"Platform"`
	StartTime     ContestTime `json:"StartTime" bson:"StartTime"`
	URL           string      `json:"url" bson:"url"`
	ChallengeType string      `json:"challenge_type,omitempty" bson:"challenge_type,omitempty"`
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

type ContestTime struct {
	time.Time
}

func (c *ContestTime) UnmarshalJSON(b []byte) error {
	timeLayout := "\"Mon, 2 Jan 2006 15:04\""
	timeLayout2 := "\"Mon Jan 2 2006 00:00\""
	ts := string(b)
	if ts == "null" {
		return nil
	}
	timeToAssign, err := time.Parse(`"`+time.RFC3339+`"`, ts)
	if err != nil {
		timeToAssign, err = time.Parse(timeLayout, ts)
		if err != nil {
			timeToAssign, err = time.Parse(timeLayout2, ts)
		}
	}
	*c = ContestTime{timeToAssign}
	return err
}
