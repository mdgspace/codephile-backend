package codechef

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
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
	Context context.Context
}

var token string

func GetBearerToken(hub *sentry.Hub) string {
	tokenURL := "https://api.codechef.com/oauth/token"
	resp, err := http.PostForm(tokenURL, map[string][]string{
		"client_id":     {os.Getenv("CLIENT_ID")},
		"client_secret": {os.Getenv("CLIENT_SECRET")},
		"grant_type":    {"client_credentials"},
		"scope":         {"public"},
	})
	if err != nil {
		log.Println(err.Error())
		hub.CaptureException(err)
		return ""
	}
	byteValue, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		hub.CaptureException(err)
		return ""
	}
	var respStruct map[string]interface{}
	err = json.Unmarshal(byteValue, &respStruct)
	if err != nil {
		hub.AddBreadcrumb(&sentry.Breadcrumb{
			Category:  "JSON parse error",
			Message:   string(byteValue),
		}, nil)
		hub.CaptureException(err)
	}
	result := respStruct["result"].(map[string]interface{})
	accessToken := result["data"].(map[string]interface{})["access_token"].(string)
	return accessToken
}

func fetchAndParseProfileData(handle string, fields string, hub *sentry.Hub) (types.CodechefProfileInfo, int) {
	profileURL := fmt.Sprintf("https://api.codechef.com/users/%s?fields=%s",
		handle, url.QueryEscape(fields))
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, profileURL, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, _ := client.Do(req)
	if resp.StatusCode != http.StatusOK {
		return types.CodechefProfileInfo{}, resp.StatusCode
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		hub.CaptureException(err)
	}
	var profileInfo types.CodechefProfileInfo
	err = json.Unmarshal(data, &profileInfo)
	if err != nil {
		hub.AddBreadcrumb(&sentry.Breadcrumb{
			Category:  "JSON parse error",
			Message:   string(data),
		}, nil)
		hub.CaptureException(err)
	}
	return profileInfo, resp.StatusCode
}

func (s Scrapper) CheckHandle() bool {
	hub := sentry.GetHubFromContext(s.Context)
	if hub == nil {
		hub = sentry.CurrentHub()
	}
	fields := "username"
	var (
		profileInfo types.CodechefProfileInfo
		status      int
	)
	for attempt := 0; attempt < 5; attempt++ {
		time.Sleep(time.Second * time.Duration(attempt))
		profileInfo, status = fetchAndParseProfileData(s.Handle, fields, hub)
		if status == http.StatusUnauthorized {
			token = GetBearerToken(hub)
			profileInfo, _ = fetchAndParseProfileData(s.Handle, fields, hub) // nolint: ineffassign
		}
		// 9002 implies rate limit exceeded
		if profileInfo.Result["data"].Code != 9002 {
			break
		}
	}
	code := profileInfo.Result["data"].Code
	return code == 9001
}

func (s Scrapper) GetProfileInfo() types.ProfileInfo {
	hub := sentry.GetHubFromContext(s.Context)
	if hub == nil {
		hub = sentry.CurrentHub()
	}
	fields := "username,fullname,organization,rankings"
	var (
		profileInfo types.CodechefProfileInfo
		status      int
	)
	for attempt := 0; attempt < 5; attempt++ {
		time.Sleep(time.Second * time.Duration(attempt))
		profileInfo, status = fetchAndParseProfileData(s.Handle, fields, hub)
		if status == http.StatusUnauthorized {
			token = GetBearerToken(hub)
			profileInfo, _ = fetchAndParseProfileData(s.Handle, fields, hub) // nolint: ineffassign
		}
		// 9002 implies rate limit exceeded
		if profileInfo.Result["data"].Code != 9002 {
			break
		}
	}
	resultData := profileInfo.Result["data"].Content
	return types.ProfileInfo{
		Name:      resultData.Fullname,
		UserName:  resultData.Username,
		School:    resultData.Organization,
		WorldRank: fmt.Sprint(resultData.Rankings["allContestRanking"].(map[string]interface{})["global"].(float64)),
	}
}

func callCodechefAPI(handle string, afterIndex int, hub *sentry.Hub) (types.CodechefSubmissions, error) {
	fields := "id, date, username, problemCode, language, result"
	submissionURL := fmt.Sprintf("https://api.codechef.com/submissions/?&username=%s&after=%d&limit=20&fields=%s",
		handle, afterIndex, url.QueryEscape(fields))
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, submissionURL, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, _ := client.Do(req)
	if resp.StatusCode == http.StatusUnauthorized {
		token = GetBearerToken(hub)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		resp, _ = client.Do(req)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		hub.CaptureException(err)
	}
	if data == nil {
		return types.CodechefSubmissions{}, errors.New("GetRequest failed. Please check connection status")
	}
	var codechefSubmissions types.CodechefSubmissions
	err = json.Unmarshal(data, &codechefSubmissions)
	if err != nil {
		hub.AddBreadcrumb(&sentry.Breadcrumb{
			Category:  "JSON parse error",
			Message:   string(data),
		}, nil)
		hub.CaptureException(err)
		log.Println(err.Error())
		return types.CodechefSubmissions{}, err
	}
	return codechefSubmissions, nil
}

func getCodechefSubmissionParts(handle string, afterIndex int, hub *sentry.Hub) ([]types.Submission, error, int) {
	codechefSubmission, err := callCodechefAPI(handle, afterIndex, hub)
	if err != nil {
		return nil, err, afterIndex
	}
	if codechefSubmission.Status != "OK" {
		log.Println("Codechef submission could not be retrieved. Retrying...")
		var newCodechefSub types.CodechefSubmissions
		for attempt := 1; attempt < 5; attempt++ {
			time.Sleep(time.Second * time.Duration(attempt))
			newCodechefSub, err = callCodechefAPI(handle, afterIndex, hub)
			if err != nil {
				return nil, err, afterIndex
			}
			if newCodechefSub.Status == "OK" {
				codechefSubmission = newCodechefSub
				break
			}
		}
		if newCodechefSub.Status == "FAILED" {
			hub.CaptureException(errors.New("codechef API repeatedly returned FAILED"))
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
		submissions[i].Name = result.ProblemCode
		submissions[i].Status = status
		submissions[i].Language = result.Language
		submissions[i].URL = "https://www.codechef.com/problems/" + result.ProblemCode
		t, err := time.Parse("2006-01-02 15:04:05", result.Date)
		if err != nil {
			hub.CaptureException(err)
		}
		submissions[i].CreationDate = t
		lastID = result.ID
	}
	return submissions, nil, lastID
}

func (s Scrapper) GetSubmissions(after time.Time) []types.Submission {
	hub := sentry.GetHubFromContext(s.Context)
	if hub == nil {
		hub = sentry.CurrentHub()
	}
	oldestSubIndex := 0
	var lastID int
	var oldestSubFound = false
	var subs []types.Submission
	var newSub []types.Submission
	var err error
	//Fetch submission until oldest submission not found
	for !oldestSubFound {
		newSub, err, lastID = getCodechefSubmissionParts(s.Handle, lastID, hub)
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
