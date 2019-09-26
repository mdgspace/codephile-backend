package scripts

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/mdg-iitr/Codephile/models/submission"
	"log"
	"regexp"
	"strings"
	"github.com/mdg-iitr/Codephile/models/profile"
)

type SpojProblems struct {
	count     string
	Submitted string
}

type SolvedProblems struct {
	problem string
	link    string
}

func GetSpojProfileInfo(handle string) profile.ProfileInfo {

	c := colly.NewCollector()
	var Profile profile.ProfileInfo

	c.OnHTML("#user-profile-left", func(e *colly.HTMLElement) {
		Name := e.ChildText("h3")
		UserName := e.ChildText("h4")
		List := strings.Split(e.ChildText(":nth-child(7)"), ":")
		School := List[1]
		List = strings.Split(e.ChildText(":nth-child(6)"), ":")
		WorldRank := List[1]
		Profile = profile.ProfileInfo{Name, UserName, School, WorldRank}
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.Visit(fmt.Sprintf("https://www.spoj.com/users/%s/", handle))

	return Profile
}

func GetSpojSubmissions(handle string) []submission.SpojSubmission {

	c := colly.NewCollector()
	var submissions []submission.SpojSubmission

	c.OnHTML("tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, elem *colly.HTMLElement) {
			Name := elem.ChildText(".sproblem a")
			URL := "https://www.spoj.com" + elem.ChildAttr(".sproblem a", "href")
			CreationDate := elem.ChildText(".status_sm span")
			status := elem.ChildText(".statusres")
			language := elem.ChildText(".slang span")
			submissions = append(submissions, submission.SpojSubmission{Name, URL, CreationDate, status, language})
		})
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.Visit(fmt.Sprintf("https://www.spoj.com/status/%s/", handle))

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

	c.Visit(fmt.Sprintf("https://www.spoj.com/users/%s/", handle))

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

	c.Visit(fmt.Sprintf("https://www.spoj.com/users/%s/", handle))
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
