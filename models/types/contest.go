package types

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/mdg-iitr/Codephile/conf"
)

type Contest struct {
	ID       int         `json:"id" bson:"id"`
	Host     string      `json:"host" bson:"host"`
	Event    string      `json:"event" bson:"event"`
	Href     string      `json:"href" bson:"href"`
	Duration int         `json:"duration" bson:"duration"`
	Start    ContestTime `json:"start" bson:"start"`
	End      ContestTime `json:"end" bson:"end"`
}

type CListResult struct {
	Meta     map[string]interface{} `json:"meta" bson:"meta"`
	Contests []Contest              `json:"objects" bson:"objects"`
}

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
	ts := string(b)
	if ts == "null" {
		return nil
	}
	ts = strings.Trim(ts, "\"")
	timeToAssign, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		ts += "Z"
		timeToAssign, err = time.Parse(time.RFC3339, ts)
	}
	*c = ContestTime{timeToAssign}
	return err
}

func (clistRes CListResult) ToResult() (Result, error) {
	var result Result
	currTime := time.Now()
	result.Timestamp = currTime.Format(time.RFC3339)
	for _, c := range clistRes.Contests {
		site, err := conf.GetSiteFromURL(c.Host)
		if err != nil {
			return Result{}, err
		}
		if diff := c.Start.Time.Sub(currTime).Seconds(); diff > 0.0 {
			upcoming := Upcoming{
				Duration:      fmt.Sprint(c.Duration),
				EndTime:       c.End,
				StartTime:     c.Start,
				Name:          c.Event,
				Platform:      site,
				URL:           c.Href,
				ChallengeType: "",
			}
			result.Upcoming = append(result.Upcoming, upcoming)
		} else {
			ongoing := Ongoing{
				EndTime:       c.End,
				Name:          c.Event,
				Platform:      site,
				URL:           c.Href,
				ChallengeType: "",
			}
			result.Ongoing = append(result.Ongoing, ongoing)
		}
	}
	return result, nil
}
