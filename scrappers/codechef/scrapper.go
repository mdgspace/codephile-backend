package codechef

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/mdg-iitr/Codephile/conf"
	"github.com/mdg-iitr/Codephile/models/types"
	// "github.com/mdg-iitr/Codephile/scrappers/common"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	// "regexp"
	// "strconv"
	// "strings"
	"time"
)

type Scrapper struct {
	Handle string
}

func GetBearerToken() string {
	tokenURL := "https://api.codechef.com/oauth/token"
	resp, err := http.PostForm(tokenURL, map[string][]string{
		"client_id": {os.Getenv("CLIENT_ID")},
		"client_secret": {os.Getenv("CLIENT_SECRET")},
		"grant_type":    {"client_credentials"},
		"scope":         {"public"},
	})
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	byteValue, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	var respStruct map[string]interface{}
	_ = json.Unmarshal(byteValue, &respStruct)
	result := respStruct["result"].(map[string]interface{})
	accessToken := result["data"].(map[string]interface{})["access_token"].(string)
	fmt.Println(accessToken)
	return accessToken
}

func (s Scrapper) CheckHandle() bool {
	path := fmt.Sprintf("https://www.codechef.com/users/%s", s.Handle)
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Get(path)
	if err != nil {
		log.Fatal(err)
		return false
	}

	if resp.StatusCode == 200 {
		return true
	}

	return false
}

func (s Scrapper) GetProfileInfo() types.ProfileInfo {
	fields := "username,fullname,organization,rankings"
	profileURL := fmt.Sprintf("https://api.codechef.com/users/%s?fields=%s",
		s.Handle, url.QueryEscape(fields))
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, profileURL, nil)
	token := GetBearerToken()
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, _ := client.Do(req)
	data, _ := ioutil.ReadAll(resp.Body)
	var profileInfo types.CodechefProfileInfo
	_ = json.Unmarshal(data, &profileInfo)
	resultData := profileInfo.Result["data"].Content
	return types.ProfileInfo{
		Name:      resultData.Fullname,
		UserName:  resultData.Username,
		School:    resultData.Organization,
		WorldRank: fmt.Sprint(resultData.Rankings["allContestRanking"].(map[string]interface{})["global"].(float64)),
	}
}

func callCodechefAPI(handle string, afterIndex int) (types.CodechefSubmissions, error) {
	fields := "id, date, username, problemCode, language, result"
	submissionURL := fmt.Sprintf("https://api.codechef.com/submissions/?&username=%s&after=%d&limit=20&fields=%s",
		handle, afterIndex, url.QueryEscape(fields))
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, submissionURL, nil)
	token := GetBearerToken()
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, _ := client.Do(req)
	data, _ := ioutil.ReadAll(resp.Body)
	if data == nil {
		return types.CodechefSubmissions{}, errors.New("GetRequest failed. Please check connection status")
	}
	var codechefSubmissions types.CodechefSubmissions
	err := json.Unmarshal(data, &codechefSubmissions)
	if err != nil {
		log.Println(err.Error())
		return types.CodechefSubmissions{}, err
	}
	return codechefSubmissions, nil
}

func getCodechefSubmissionParts(handle string, afterIndex int) ([]types.Submission, error, int) {
	codechefSubmission, err := callCodechefAPI(handle, afterIndex)
	if err != nil {
		return nil, err, afterIndex
	}
	if codechefSubmission.Status != "OK" {
		log.Println("Codechef submission could not be retrieved. Retrying...")
		var newCodechefSub types.CodechefSubmissions
		for attempt := 1; attempt < 5; attempt++ {
			time.Sleep(time.Second * time.Duration(attempt))
			newCodechefSub, err = callCodechefAPI(handle, afterIndex)
			if err != nil {
				return nil, err, afterIndex
			}
			if newCodechefSub.Status == "OK" {
				codechefSubmission = newCodechefSub
				break
			}
		}
		if newCodechefSub.Status == "FAILED" {
			return nil, errors.New("codechef API repeatedly returned FAILED"), 0
		}
	}
	submissions := make([]types.Submission, len(codechefSubmission.Result.Data.Content))
	var lastID int
	for i, result := range codechefSubmission.Result.Data.Content {
		var status string
		switch result.Result {
		case "AC":
			status = StatusCorrect
		case "CTE":
			status = StatusCompilationError
		case "RTE":
			status = StatusRuntimeError
		case "WA":
			status = StatusWrongAnswer
		default:
			status = StatusWrongAnswer
		}
		submissions[i].Status = status
		submissions[i].Language = result.Language
		submissions[i].URL = "https://www.codechef.com/problems/" + result.ProblemCode
		t, _ := time.Parse("2006-01-02 15:04:05", result.Date)
		submissions[i].CreationDate = t
		lastID = result.ID
	}
	return submissions, nil, lastID
}

func (s Scrapper) GetSubmissions(after time.Time) []types.Submission {
	log.Println("Hello")
	oldestSubIndex := 0
	var lastID int
	var oldestSubFound = false
	var subs []types.Submission
	var newSub []types.Submission
	var err error
	//Fetch submission until oldest submission not found
	for !oldestSubFound {
		newSub, err, lastID = getCodechefSubmissionParts(s.Handle, lastID)
		if err != nil {
			log.Println(err.Error())
			return nil
		}
		//Check for repetition of previous fetched submission
		if len(newSub) != 0 {
			for _, sub := range newSub {
				subs = append(subs, sub)
				oldestSubIndex += 1
				if sub.CreationDate.Equal(after) || sub.CreationDate.Before(after) {
					oldestSubFound = true
					break
				}
			}
		} else {
			break
		}
	}
	subs = subs[0:oldestSubIndex]
	return subs
}
