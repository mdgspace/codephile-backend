package leetcode

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/mdg-iitr/Codephile/models/types"
)

type Scrapper struct {
	Handle  string
	Context context.Context
}

func leetcodeGraphQLRequest(query string) (map[string]interface{}, error) {
	jsonData := map[string]string{
		"query": query,
	}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post("https://leetcode.com/graphql", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	responseValue, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var responseData map[string]interface{}
	err = json.Unmarshal(responseValue, &responseData)
	if err != nil {
		return nil, err
	}
	data := responseData["data"].(map[string]interface{})
	return data, err
}

func (s Scrapper) GetProfileInfo() types.ProfileInfo {
	hub := sentry.GetHubFromContext(s.Context)
	if hub == nil {
		hub = sentry.CurrentHub()
	}

	query := `
		{
			matchedUser(username: "` + s.Handle + `") {
				username
				profile {
					realName
					school
					ranking
				}
				submitStats {
					acSubmissionNum {
						submissions
					}
					totalSubmissionNum {
						submissions
					}
				}
			}
		}
	`
	responseData, err := leetcodeGraphQLRequest(query)
	if err != nil {
		log.Println(err.Error())
		hub.CaptureException(err)
		return types.ProfileInfo{}
	}
	matchedUser := responseData["matchedUser"].(map[string]interface{})
	UserName := matchedUser["username"].(string)
	profile := matchedUser["profile"].(map[string]interface{})
	Name := profile["realName"].(string)
	var School string
	if profile["school"] != nil {
		School = profile["school"].(string)
	} else {
		School = ""
	}
	WorldRank := profile["ranking"].(float64)
	submitStats := matchedUser["submitStats"].(map[string]interface{})
	acSubmissionNumAll := submitStats["acSubmissionNum"].([]interface{})[0].(map[string]interface{})["submissions"].(float64)
	totalSubmissionNumAll := submitStats["totalSubmissionNum"].([]interface{})[0].(map[string]interface{})["submissions"].(float64)
	var Accuracy float64
	if totalSubmissionNumAll == 0 {
		Accuracy = 0
	} else {
		Accuracy = (acSubmissionNumAll / totalSubmissionNumAll) * 100
	}
	return types.ProfileInfo{Name: Name, UserName: UserName, School: School, WorldRank: fmt.Sprintf("%.0f", WorldRank), Accuracy: fmt.Sprintf("%.2f", Accuracy)}
}

func (s Scrapper) CheckHandle() (bool, error) {
	hub := sentry.GetHubFromContext(s.Context)
	if hub == nil {
		hub = sentry.CurrentHub()
	}

	query := `
		{
			matchedUser(username: "` + s.Handle + `") {
				username
			}
		}
	`
	responseData, err := leetcodeGraphQLRequest(query)
	if err != nil {
		log.Println(err.Error())
		hub.CaptureException(err)
		return false, err
	}
	matchedUser := responseData["matchedUser"]
	return matchedUser != nil, err
}
