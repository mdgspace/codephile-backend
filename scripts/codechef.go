package scripts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type GraphPoint struct {
	ContestName string
	Date        time.Time
	Rating      float64
}

type ProfileInfo struct {
	Name     string
	UserName string
	School   string
}

type Submission struct {
	Name         string
	URL          string
	CreationDate string
	status       string
	points       string
}

func zero_pad(year *string) {
	fmt.Println(len(*year))
	if len(*year) == 1 {
		*year = "0" + *year
	}
}

// func main() {
// 	graph := GetGraphData("ritikkne")
// 	profile := GetProfileInfo("ritikkne")
// 	submittions := GetSubmissions("ritikkne")

// 	fmt.Println(submittions)
// 	fmt.Println(profile)
// 	fmt.Println(graph)
// }

func Get_Request(path string) []byte {
	resp, err := http.Get(path)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	byteValue, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return byteValue
}

func CheckHandle(handle string) bool {
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

func GetGraphData(handle string) []GraphPoint {
	r, _ := regexp.Compile("\\[.*?\\]")

	c := colly.NewCollector()
	var long_contest_data []GraphPoint

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

				long_contest_data = append(long_contest_data, GraphPoint{contest_name, time, ranking})
			}

			// fmt.Println(long_contest_data)

			var data2 []map[string]string
			err = json.Unmarshal([]byte(lunch_time), &data2)
			if err != nil {
				fmt.Println(err)
			}

			var lunch_contest_data []GraphPoint

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

				lunch_contest_data = append(lunch_contest_data, GraphPoint{contest_name, time, ranking})
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

func GetProfileInfo(handle string) ProfileInfo {
	path := fmt.Sprintf("https://www.codechef.com/api/ratings/all?sortBy=global_rank&order=asc&search=%s&page=1&itemsPerPage=20", handle)
	byteValue := Get_Request(path)
	var JsonInterFace interface{}
	json.Unmarshal(byteValue, &JsonInterFace)
	Profile := JsonInterFace.(map[string]interface{})["list"].([]interface{})[0].(map[string]interface{})

	// all_rating := Profile["all_rating"]
	// country := Profile["country"]
	// country_code := Profile["country_code"]
	// country_rank := Profile["country_rank"]
	// diff := Profile["diff"]
	// global_rank := Profile["global_rank"]
	Institution := Profile["institution"].(string)
	// institution_type := Profile["institution_type"]
	Name := Profile["name"].(string)
	// rating := Profile["rating"]
	UserName := Profile["username"].(string)
	return ProfileInfo{Name, UserName, Institution}
}

func GetSubmissions(handle string) []Submission {

	user_url := "http://www.codechef.com/recent/user?user_handle=" + handle
	byteValue := Get_Request(user_url)
	var JsonInterFace interface{}
	json.Unmarshal(byteValue, &JsonInterFace)
	data := JsonInterFace.(map[string]interface{})

	var submissions []Submission

	max_page := int(data["max_page"].(float64))
	content := data["content"].(string)

	submissions = append(submissions, GetSubmissionsFromString(content)...)

	for i := 1; i < max_page; i++ {
		user_url = fmt.Sprintf("http://www.codechef.com/recent/user?user_handle=%s&page=%d", handle, i)

		byteValue = Get_Request(user_url)
		json.Unmarshal(byteValue, &JsonInterFace)
		data = JsonInterFace.(map[string]interface{})

		content := data["content"].(string)

		submissions = append(submissions, GetSubmissionsFromString(content)...)

	}

	return submissions

}

func GetSubmissionsFromString(content string) []Submission {

	var submissions []Submission

	data := strings.Split(content, "<tr >")
	for i := 1; i < 4; i++ {
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
		url := "http://www.codechef.com" + prob

		// Submission status
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
		submissions = append(submissions, Submission{prob, url, tos, st, points})

	}

	return submissions
}
