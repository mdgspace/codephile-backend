package scripts

import (
	"encoding/json"
	"fmt"
	"github.com/mdg-iitr/Codephile/models/submission"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"github.com/mdg-iitr/Codephile/models/profile"
	"github.com/gocolly/colly"
)

type CodechefGraphPoint struct {
	ContestName string
	Date        time.Time
	Rating      float64
}

func zero_pad(year *string) {
	fmt.Println(len(*year))
	if len(*year) == 1 {
		*year = "0" + *year
	}
}

func CheckCodechefHandle(handle string) bool {
	path := fmt.Sprintf("https://www.codechef.com/users/%s", handle)
	resp, err := http.Get(path)
	if err != nil {
		log.Fatal(err)
		return false
	}

	if resp.StatusCode == 200 {
		return true
	}

	return false
}

func GetCodechefGraphData(handle string) []CodechefGraphPoint {
	r, _ := regexp.Compile("\\[.*?\\]")

	c := colly.NewCollector()
	var long_contest_data []CodechefGraphPoint

	c.OnHTML("script", func(e *colly.HTMLElement) {

		if e.Index == 29 {

			data := r.FindAllString(e.Text, -1)
			// fmt.Println(data)
			// long_ratings := data[0]
			long_challenge := data[2]
			// short_ratings := data[4]
			lunch_time := data[6]

			var data1 []map[string]string
			err := json.Unmarshal([]byte(long_challenge), &data1)
			if err != nil {
				fmt.Println(err)
			}

			for _, event := range data1 {
				// fmt.Println(event)

				year := event["getyear"]
				month := event["getmonth"]
				ranking, err := strconv.ParseFloat(event["rank"], 64)
				time, err := time.Parse("2006-1", fmt.Sprintf("%v-%s", year, month))
				if err != nil {
					panic(err)
				}
				contest_name := fmt.Sprintf("%s Long Challenge 20%s", time.Format("January"), time.Format("06"))
				// contest_url := fmt.Sprintf("https://www.codechef.com/%s%s", time.Format("JAN"), time.Format("06"))

				long_contest_data = append(long_contest_data, CodechefGraphPoint{contest_name, time, ranking})
			}

			// fmt.Println(long_contest_data)

			var data2 []map[string]string
			err = json.Unmarshal([]byte(lunch_time), &data2)
			if err != nil {
				fmt.Println(err)
			}

			var lunch_contest_data []CodechefGraphPoint

			for _, event := range data1 {

				// fmt.Println(event)
				year := event["getyear"]
				month := event["getmonth"]
				// code := event["code"]
				ranking, err := strconv.ParseFloat(event["rank"], 64)
				time, err := time.Parse("2006-1", fmt.Sprintf("%v-%s", year, month))
				if err != nil {
					panic(err)
				}
				contest_name := fmt.Sprintf("%s Lunch Time 20%s", time.Format("January"), time.Format("06"))
				// contest_url := fmt.Sprintf("https://www.codechef.com/" + code)

				lunch_contest_data = append(lunch_contest_data, CodechefGraphPoint{contest_name, time, ranking})
			}

			// fmt.Println(lunch_contest_data)

			long_contest_data = append(long_contest_data, lunch_contest_data...)
		}
	})

	c.OnError(func(_ *colly.Response, err error) {
		// return nil
		fmt.Println("Something went wrong:", err)
	})
	c.Visit(fmt.Sprintf("https://www.codechef.com/users/%s", handle))

	return long_contest_data

}

func GetCodechefProfileInfo(handle string) profile.ProfileInfo {

	/* path := fmt.Sprintf("https://www.codechef.com/api/ratings/all?sortBy=global_rank&order=asc&search=%s&page=1&itemsPerPage=40", handle)
		byteValue := GetRequest(path)
		var JsonInterFace interface{}
		// var Profile map[string]interface{}
		json.Unmarshal(byteValue, &JsonInterFace)
		pagesToTraverse := int(JsonInterFace.(map[string]interface{})["availablePages"].(float64))
		Profile := JsonInterFace.(map[string]interface{})["list"].([]interface{})[0].(map[string]interface{})
	    for i:=1 ; i <= pagesToTraverse ; i++ {
			newPath := fmt.Sprintf("https://www.codechef.com/api/ratings/all?sortBy=global_rank&order=asc&search=%s&page=%d&itemsPerPage=40", handle , i)
			newbyteValue := GetRequest(newPath)
			var newJsonInterFace interface{}
			json.Unmarshal(newbyteValue , &newJsonInterFace)
			log.Println(newJsonInterFace.(map[string]interface{})["list"].([]interface{})[0])
			for j:=0 ; j<=39 ; j++ {
				log.Println(newJsonInterFace.(map[string]interface{})["list"].([]interface{})[j])//.(map[string]interface{})["username"].(string))
				if newJsonInterFace.(map[string]interface{})["list"].([]interface{})[j].(map[string]interface{})["username"].(string) == handle {
				  Profile = newJsonInterFace.(map[string]interface{})["list"].([]interface{})[j].(map[string]interface{})
				  break
				}
			}
		}
		// Profile := JsonInterFace.(map[string]interface{})["list"].([]interface{})[0].(map[string]interface{})

		// all_rating := Profile["all_rating"]
		// country := Profile["country"]
		// country_code := Profile["country_code"]
		// country_rank := Profile["country_rank"]
		// diff := Profile["diff"]
		global_rank := Profile["global_rank"].(float64)
		global_rank_string := strconv.FormatFloat(global_rank,'f',0,64)
		Institution := Profile["institution"].(string)
		// institution_type := Profile["institution_type"]
		Name := Profile["name"].(string)
		// rating := Profile["rating"]
		UserName := Profile["username"].(string)
		return profile.ProfileInfo{Name, UserName, Institution,global_rank_string} */

	//scrapping the profile
	c := colly.NewCollector()
	var Profile profile.ProfileInfo
	var School string
	c.OnHTML(".user-profile-container", func(e *colly.HTMLElement) {
		Name := e.ChildText("h2")
		UserName := handle
		for i := 2; i <= 10; i++ {
			if e.ChildText(fmt.Sprintf(".user-details .side-nav li:nth-child(%d) label", i)) == "Institution:" {
				School = e.ChildText(fmt.Sprintf(".user-details .side-nav li:nth-child(%d) span", i))
			}
		}
		WorldRank := e.ChildText(".rating-ranks .inline-list li:nth-child(1) a")
		Profile = profile.ProfileInfo{Name, UserName, School, WorldRank}
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.Visit(fmt.Sprintf("https://www.codechef.com/users/%s", handle))
	return Profile
}

func GetCodechefSubmissions(handle string, after time.Time) []submission.CodechefSubmission {
	var oldestSubIndex, current int;
	var oldestSubFound = false
	subs := []submission.CodechefSubmission{{CreationDate: time.Now()}}
	//Fetch submission until oldest submission not found
	for !oldestSubFound {
		newSub := GetCodechefSubmissionParts(handle, current);
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
func GetCodechefSubmissionParts(handle string, pageNo int) []submission.CodechefSubmission {
	var JsonInterFace interface{}
	user_url := fmt.Sprintf("http://www.codechef.com/recent/user?user_handle=%s&page=%d", handle, pageNo)
	fmt.Println(user_url)
	byteValue := GetRequest(user_url)
	json.Unmarshal(byteValue, &JsonInterFace)
	data := JsonInterFace.(map[string]interface{})
	content := data["content"].(string)

	return GetSubmissionsFromString(content)

}

func GetSubmissionsFromString(content string) []submission.CodechefSubmission {

	var submissions []submission.CodechefSubmission

	data := strings.Split(content, "<tr >")
	for i := 1; i <= len(data)-1; i++ {
		info := strings.Split(data[i], "</tr>")[0]

		contents := strings.Split(info, "<td >")

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
		data := GetRequest(fmt.Sprintf("https://www.codechef.com/api/contests/PRACTICE/problems/%s", prob))
		var JsonInterface map[string]interface{}
		json.Unmarshal(data, &JsonInterface)
		htmlTag := JsonInterface["tags"].(string)
		htmlTag = regexp.MustCompile("<a class='(.*?)'>").ReplaceAllString(htmlTag, "")
		tags := strings.Split(htmlTag, "</a>, ")
		tags[len(tags)-1] = strings.Replace(tags[len(tags)-1], "</a>", "", 1)
		// SpojSubmission status
		stat := strings.Split(strings.Split(contents[3], "/misc/")[1], ".gif")[0]
		st := "AC"
		if stat == "tick-icon" {
			st = "AC"
		} else if stat == "cross-icon" {
			st = "WA"
		} else if stat == "alert-icon" {
			st = "CE"
		} else if stat == "runtime-error" {
			st = "RE"
		} else if stat == "clock_error" {
			st = "TLE"
		} else {
			st = "OTH"
		}

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
		submissionTime, err := time.Parse("03:04 PM 02/01/06", tos)
		if err != nil {
			log.Println(err.Error())
		}
		submissions = append(submissions, submission.CodechefSubmission{prob, url, submissionTime, st, points, tags})

	}

	return submissions
}
