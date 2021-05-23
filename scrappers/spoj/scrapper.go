package spoj

import (
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/gocolly/colly"
	. "github.com/mdg-iitr/Codephile/conf"
	"github.com/mdg-iitr/Codephile/models/types"
	"log"
	"regexp"
	"strings"
	"time"
)

type Scrapper struct {
	Handle  string
	Context context.Context
}

func (s Scrapper) GetProfileInfo() types.ProfileInfo {
	hub := sentry.GetHubFromContext(s.Context)
	if hub == nil {
		hub = sentry.CurrentHub()
	}
	c := colly.NewCollector()
	var Profile types.ProfileInfo

	c.OnHTML("#user-profile-left", func(e *colly.HTMLElement) {
		Name := e.ChildText("h3")
		flag := 0
		var WorldRank string
		var School string
		UserName := s.Handle
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
		hub.CaptureException(err)
	})

	err := c.Visit(fmt.Sprintf("https://www.spoj.com/users/%s/", s.Handle))
	if err != nil {
		hub.CaptureException(err)
		log.Println(err.Error())
	}
	return Profile
}

func (s Scrapper) GetSubmissions(after time.Time) []types.Submission {
	hub := sentry.GetHubFromContext(s.Context)
	if hub == nil {
		hub = sentry.CurrentHub()
	}
	var oldestSubIndex, current int
	var oldestSubFound = false
	subs := []types.Submission{{CreationDate: time.Now()}}
	//Fetch submission until oldest submission not found
	for !oldestSubFound {
		newSub := getSubmissionParts(s.Handle, current, hub)
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
	if len(subs) < 2 {
		return nil
	}
	subs = subs[1 : oldestSubIndex+1]
	return subs
}

func getSubmissionParts(handle string, afterIndex int, hub *sentry.Hub) []types.Submission {

	c := colly.NewCollector()
	var submissions []types.Submission

	c.OnHTML("tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, elem *colly.HTMLElement) {
			Name := elem.ChildText(".sproblem a")
			URL := "https://www.spoj.com" + elem.ChildAttr(".sproblem a", "href")
			str_date := elem.ChildText(".status_sm span")
			CreationDate, err := time.Parse("2006-01-02 15:04:05", str_date)
			if err != nil {
				sentry.CaptureException(err)
				log.Println(err.Error())
			}
			status := elem.ChildText(".statusres")
			language := elem.ChildText(".slang span")
			points := 0
			switch status {
			case "accepted":
				status = StatusCorrect
			case "wrong answer":
				status = StatusWrongAnswer
			case "compilation error":
				status = StatusCompilationError
			case "runtime error":
				status = StatusRuntimeError
			default:
				status = StatusWrongAnswer
			}
			if status == StatusCorrect {
				points = 100
			}
			tags := getProbTags(URL)
			submissions = append(submissions, types.Submission{Name: Name, URL: URL, CreationDate: CreationDate, Status: status, Language: language, Points: points, Tags: tags})
		})
	})

	c.OnError(func(_ *colly.Response, err error) {
		hub.CaptureException(err)
		fmt.Println("Something went wrong:", err)
	})

	c.OnRequest(func(request *colly.Request) {
		fmt.Println(request.URL)
	})
	err := c.Visit(fmt.Sprintf("https://www.spoj.com/status/%s/all/start=%d", handle, afterIndex))
	if err != nil {
		hub.CaptureException(err)
		log.Println(err.Error())
	}

	return submissions
}

func (s Scrapper) CheckHandle() (bool, error) {
	hub := sentry.GetHubFromContext(s.Context)
	if hub == nil {
		hub = sentry.CurrentHub()
	}
	c := colly.NewCollector()
	var valid = false
	var err error
	c.OnResponse(func(response *colly.Response) {
		valid, err = regexp.Match("user-profile-left", response.Body)
		if err != nil {
			hub.CaptureException(err)
		}
	})
	err = c.Visit(fmt.Sprintf("https://www.spoj.com/users/%s/", s.Handle))
	if err != nil {
		log.Println(err.Error())
		hub.CaptureException(err)
	}
	return valid, nil
}

func getProbTags(url string) []string {
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
