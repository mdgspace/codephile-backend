package scripts

import (
	"fmt"
	"github.com/mdg-iitr/Codephile/models/types"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type SpojProblems struct {
	count     string
	Submitted string
}

type SolvedProblems struct {
	problem string
	link    string
}

func GetSpojProfileInfo(handle string) types.ProfileInfo {

	c := colly.NewCollector()
	var Profile types.ProfileInfo

	c.OnHTML("#user-profile-left", func(e *colly.HTMLElement) {
		Name := e.ChildText("h3")
		flag := 0
		var WorldRank string
		var School string
		UserName := handle
		for i := 4; i <= 6; i++ {
			cssSelector1 := fmt.Sprintf(":nth-child(%d)", i)
			if strings.Split(e.ChildText(cssSelector1), ":")[0] == "World Rank" {
				WorldRank = strings.Split(e.ChildText(cssSelector1), ":")[1]
				flag = i
				break
			}
		}
		cssSelector2 := fmt.Sprintf(":nth-child(%d)", flag+1)

		defer func() {
			if r := recover(); r != nil {
				//catching index out of range exception in fetching School
				School = ""
				Profile = types.ProfileInfo{Name: Name, UserName: UserName, School: School, WorldRank: WorldRank}
			}
		}()

		School = strings.Split(e.ChildText(cssSelector2), ":")[1]

		Profile = types.ProfileInfo{Name: Name, UserName: UserName, School: School, WorldRank: WorldRank}
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	err := c.Visit(fmt.Sprintf("https://www.spoj.com/users/%s/", handle))
	if err != nil {
		log.Println(err.Error())
	}
	return Profile
}

func GetSpojSubmissions(handle string, after time.Time) []types.Submission {
	var oldestSubIndex, current int
	var oldestSubFound = false
	subs := []types.Submission{{CreationDate: time.Now()}}
	//Fetch submission until oldest submission not found
	for !oldestSubFound {
		newSub := GetSpojSubmissionParts(handle, current)
		//Check for repetition of previous fetched submission
		if len(newSub) != 0 && newSub[0].CreationDate.Before(subs[len(subs)-1].CreationDate) {
			for i, sub := range newSub {
				subs = append(subs, sub)
				oldestSubIndex = current + i
				if sub.CreationDate.Equal(after) || sub.CreationDate.Before(after) {
					oldestSubFound = true
					break
				}
			}
			//20 submissions per page
			current += 20
		} else {
			oldestSubIndex++
			break
		}
	}
	subs = subs[1 : oldestSubIndex+1]
	return subs
}
func GetSpojSubmissionParts(handle string, afterIndex int) []types.Submission {

	c := colly.NewCollector()
	var submissions []types.Submission

	c.OnHTML("tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, elem *colly.HTMLElement) {
			Name := elem.ChildText(".sproblem a")
			URL := "https://www.spoj.com" + elem.ChildAttr(".sproblem a", "href")
			str_date := elem.ChildText(".status_sm span")
			CreationDate, err := time.Parse("2006-01-02 15:04:05", str_date)
			if err != nil {
				log.Println(err.Error())
			}
			status := elem.ChildText(".statusres")
			language := elem.ChildText(".slang span")
			points := 0
			if status == "accepted" {
				points = 100
			}
			tags := GetProbTags(URL)
			submissions = append(submissions, types.Submission{Name: Name, URL: URL, CreationDate: CreationDate, Status: status, Language: language, Points: points, Tags: tags})
		})
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnRequest(func(request *colly.Request) {
		fmt.Println(request.URL)
	})
	err := c.Visit(fmt.Sprintf("https://www.spoj.com/status/%s/all/start=%d", handle, afterIndex))
	if err != nil {
		log.Println(err.Error())
	}

	return submissions
}

func GetSpojProblems(handle string) SpojProblems {

	c := colly.NewCollector()
	var problems SpojProblems
	c.OnHTML(".dl-horizontal", func(e *colly.HTMLElement) {
		count := e.ChildText(":nth-child(2)")
		submitted := e.ChildText(":nth-child(4)")

		problems = SpojProblems{count, submitted}
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	err := c.Visit(fmt.Sprintf("https://www.spoj.com/users/%s/", handle))
	if err != nil {
		log.Println(err.Error())
	}
	return problems
}

func GetSpojSolvedProblems(handle string) []SolvedProblems {

	c := colly.NewCollector()
	var solved []SolvedProblems
	c.OnHTML("#user-profile-tables tr", func(e *colly.HTMLElement) {
		Name := e.ChildText("a")
		link := "https://www.spoj.com" + e.ChildAttr("a", "href")
		solved = append(solved, SolvedProblems{Name, link})
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	err := c.Visit(fmt.Sprintf("https://www.spoj.com/users/%s/", handle))
	if err != nil {
		log.Println(err.Error())
	}
	return solved

}
func CheckSpojHandle(handle string) bool {
	c := colly.NewCollector()
	var valid = false
	c.OnResponse(func(response *colly.Response) {
		valid, _ = regexp.Match("user-profile-left", response.Body)
	})
	err := c.Visit(fmt.Sprintf("https://www.spoj.com/users/%s/", handle))
	if err != nil {
		log.Println(err.Error())
	}
	return valid
}

func GetProbTags(url string) []string {
	var tags []string
	c := colly.NewCollector()
	c.OnHTML(".problem-tag", func(e *colly.HTMLElement) {
		tags = append(tags, e.Text)
	})
	err := c.Visit(url)
	if err != nil {
		log.Println("could not fetch tags")
	}
	return tags
}
