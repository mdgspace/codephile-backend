package profile

import(
	  "encoding/json"
	  "errors"
)

type ProfileInfo struct{
	Name      string   `bson:"name" json:"name" schema:"name"`
	UserName  string   `bson:"userName" json:"userName" schema:"userName"`
	School    string   `bson:"school" json:"school" schema:"school"`
	WorldRank string   `bson:"rank" json:"rank" schema:"rank"`
	Accuracy  string   `bson:"accuracy" json:"accuracy" schema:"accuracy"`
}

type Profile struct {
	Website     string          `bson:"website" json:"website" schema:"website"`
	Profileinfo ProfileInfo     `bson:"profile" json:"profile" schema:"profile"`
}

//create an allProfilesStruct
type AllProfiles struct {
	CodechefProfile    Profile   `bson:"codechefProfile" json:"codechefProfile"`
	CodeforcesProfile  Profile   `bson:"codeforcesProfile" json:"codeforcesProfile"`
	HackerrankProfile  Profile   `bson:"hackerrankProfile" json:"hackerrankProfile"`
	SpojProfile        Profile   `bson:"spojProfile" json:"spojProfile"`
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