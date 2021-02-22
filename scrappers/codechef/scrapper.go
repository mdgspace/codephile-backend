package codechef

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/mdg-iitr/Codephile/conf"
	"github.com/mdg-iitr/Codephile/models/types"
	"github.com/mdg-iitr/Codephile/scrappers/common"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
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

func (s Scrapper) GetSubmissions(after time.Time) []types.Submission {
	var oldestSubIndex, current int
	var oldestSubFound = false
	subs := []types.Submission{{CreationDate: time.Now()}}
	//Fetch submission until oldest submission not found
	for !oldestSubFound {
		newSub, err := getSubmissionParts(s.Handle, current)
		if err != nil {
			return nil
		}
		//Check for repetition of previous fetched submission
		if newSub[0].CreationDate.Before(subs[len(subs)-1].CreationDate) {
			for i, sub := range newSub {
				subs = append(subs, sub)
				//12 submissions per page
				oldestSubIndex = 12*current + i
				if sub.CreationDate.Equal(after) || sub.CreationDate.Before(after) {
					oldestSubFound = true
					break
				}
			}
			current++
		} else {
			oldestSubIndex++
			break
		}
	}
	subs = subs[1 : oldestSubIndex+1]
	return subs
}

//Get submissions of a user after an index.
//Returns an error if unsuccessful
//On receiving the error caller should return empty submission list
func getSubmissionParts(handle string, pageNo int) ([]types.Submission, error) {
	var JsonInterFace interface{}
	user_url := fmt.Sprintf("http://www.codechef.com/recent/user?user_handle=%s&page=%d", handle, pageNo)
	fmt.Println(user_url)
	byteValue := common.HitGetRequest(user_url)
	if byteValue == nil {
		return nil, errors.New("GetRequest failed. Please check connection status")
	}
	err := json.Unmarshal(byteValue, &JsonInterFace)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	data := JsonInterFace.(map[string]interface{})
	content := data["content"].(string)

	return getSubmissionsFromString(content)

}

func getSubmissionsFromString(content string) ([]types.Submission, error) {

	var submissions []types.Submission

	data := strings.Split(content, "<tr >")
	for i := 1; i <= len(data)-1; i++ {
		info := strings.Split(data[i], "</tr>")[0]

		contents := strings.Split(info, "<td >")
		if len(contents) < 2 {
			return nil, errors.New("Invalid Handle")
		}
		// tos = time_of_submission
		tos := strings.Split(contents[1], "</td>")[0]
		tos = strings.Replace(tos, "\\", "", -1)

		// time, err := time.Parse("03:04 PM 02/01/06", tos)
		// if err != nil {
		// 	time = tos
		// }
		// fmt.Println(err)

		// Problem name/url
		prob := strings.TrimRight(strings.Split(contents[2], ">")[1], "</a")
		url := "http://www.codechef.com/problems/" + prob
		data := common.HitGetRequest(fmt.Sprintf("https://www.codechef.com/api/contests/PRACTICE/problems/%s", prob))
		var JsonInterface map[string]interface{}
		err := json.Unmarshal(data, &JsonInterface)
		if err != nil {
			log.Println(err.Error())
		}
		var tags []string
		if JsonInterface["tags"] != nil {
			htmlTag := JsonInterface["tags"].(string)
			htmlTag = regexp.MustCompile("<a class='(.*?)'>").ReplaceAllString(htmlTag, "")
			tags = strings.Split(htmlTag, "</a>, ")
			tags[len(tags)-1] = strings.Replace(tags[len(tags)-1], "</a>", "", 1)
		}
		// SpojSubmission status
		stat := strings.Split(strings.Split(contents[3], "/misc/")[1], ".gif")[0]
		var st string
		if stat == "tick-icon" {
			st = StatusCorrect
		} else if stat == "cross-icon" {
			st = StatusWrongAnswer
		} else if stat == "alert-icon" {
			st = StatusCompilationError
		} else if stat == "runtime-error" {
			st = StatusRuntimeError
		} else if stat == "clock_error" {
			st = StatusTimeLimitExceeded
		} else {
			st = StatusWrongAnswer
		}

		//Language used
		lang := strings.Split(contents[4], "</td>")[0]

		//  Question points
		pts := strings.Split(contents[3], "<br/>")
		var points string

		if len(pts) == 2 {
			points = strings.Split(pts[1], "<br />")[0]
		} else {
			if st == "AC" {
				points = "100"
			} else {
				points = "0"
			}
		}

		//  Language
		// lang := strings.TrimRight(contents[4], "</td>")

		var submissionTime time.Time
		//List[0] = number of hours or minutes to reduce
		//List[1] = hours or minutes
		//List[2] = "ago"
		List := strings.Split(tos, " ")
		if List[2] == "ago" {
			count, err := strconv.Atoi(List[0])
			if err != nil {
				log.Println(err.Error())
			}
			now := time.Now()
			if List[1] == "min" {
				submissionTime = now.Add(time.Duration(-count) * time.Minute)
			} else {
				submissionTime = now.Add(time.Duration(-count) * time.Hour)
			}
		} else {
			submissionTime, err = time.Parse("03:04 PM 02/01/06", tos)
			if err != nil {
				log.Println(err.Error())
			}
		}
		pt, err := strconv.Atoi(points)
		if err != nil {
			pt = 0
		}
		submissions = append(submissions, types.Submission{
			Name:         prob,
			URL:          url,
			CreationDate: submissionTime,
			Status:       st,
			Points:       pt,
			Tags:         tags,
			Language:     lang,
		})

	}

	return submissions, nil
}
