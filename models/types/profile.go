package types

import (
	"encoding/json"
	"errors"
)

type ProfileInfo struct {
	Name      string `bson:"name" json:"name" schema:"name"`
	UserName  string `bson:"userName" json:"userName" schema:"userName"`
	School    string `bson:"school" json:"school" schema:"school"`
	WorldRank string `bson:"rank" json:"rank" schema:"rank"`
	Accuracy  string `bson:"accuracy" json:"accuracy" schema:"accuracy"`
}

//create an allProfilesStruct
type AllProfiles struct {
	CodechefProfile   ProfileInfo `bson:"codechefProfile" json:"codechefProfile"`
	CodeforcesProfile ProfileInfo `bson:"codeforcesProfile" json:"codeforcesProfile"`
	HackerrankProfile ProfileInfo `bson:"hackerrankProfile" json:"hackerrankProfile"`
	SpojProfile       ProfileInfo `bson:"spojProfile" json:"spojProfile"`
	LeetcodeProfile   ProfileInfo `bson:"leetcodeProfile" json:"leetcodeProfile"`
}

//UnmarshalJSON implements the unmarshaler interface for CodeforcesProfileInfo
func (data *ProfileInfo) UnmarshalJSON(b []byte) error {
	var profile map[string]interface{}
	err := json.Unmarshal(b, &profile)
	if profile["status"] != "OK" {
		return errors.New("Bad Request")
	}
	result := profile["result"].([]interface{})[0].(map[string]interface{})
	if result["firstName"] != nil && result["lastName"] != nil {
		data.Name = result["firstName"].(string) + result["lastName"].(string)
	}
	data.UserName = result["handle"].(string)
	// data.JoinDate = time.Unix(int64(result["registrationTimeSeconds"].(float64)), 0)
	if result["organization"] != nil {
		data.School = result["organization"].(string)
	}
	data.WorldRank = ""
	return err
}

type SolvedProblemsCount struct {
	Codechef   int `json:"codechef"`
	Codeforces int `json:"codeforces"`
	Hackerrank int `json:"hackerrank"`
	Spoj       int `json:"spoj"`
	Leetcode   int `json:"leetcode"`
}

type CodechefProfileInfo struct {
	Status string                 `json:"status"`
	Result map[string]ProfileData `json:"result"`
}

type ProfileData struct {
	Content ProfileInfoContent `json:"content"`
	Code    int64              `json:"code"`
	Message string             `json:"message"`
}
type ProfileInfoContent struct {
	Username     string                 `json:"username"`
	Fullname     string                 `json:"fullname"`
	Rankings     map[string]interface{} `json:"rankings"`
	Organization string                 `json:"organization"`
}

type GraphQLResponse struct {
	Data LeetcodeData
}

type LeetcodeData struct {
	MatchedUser LeetcodeMatchedUser
}

type LeetcodeMatchedUser struct {
	Username    string
	Profile     LeetcodeProfile
	SubmitStats LeetcodeSubmitStats
}

type LeetcodeProfile struct {
	RealName string
	School   string
	Ranking  float64
}

type LeetcodeSubmitStats struct {
	AcSubmissionNum    []LeetcodeSubmission
	TotalSubmissionNum []LeetcodeSubmission
}

type LeetcodeSubmission struct {
	Submissions float64
}
