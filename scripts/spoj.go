package scripts

import (
	"fmt"
	"strings"
	"github.com/gocolly/colly"
)

type ProfileInfo struct {
	Name      string
	UserName  string
	School    string
	WorldRank string 
}

type Submission struct {
	Name         string
	URL          string
	CreationDate string
	status       string
	language     string
}

type Problems struct {
	count       string
	Submitted   string
}

type SolvedProblems struct {
	problem string
	link    string
}

func main(){
	Profile := GetProfileInfo("boemogensen")
	fmt.Println(Profile)
	submissions := GetSubmissions("boemogensen")
	fmt.Println(submissions)
	problems := GetProblems("boemogensen")
	fmt.Println(problems)
	solved := GetSolvedProblems("boemogensen")
	fmt.Println(solved)
}

func GetProfileInfo(handle string) ProfileInfo {

	c := colly.NewCollector()
	var Profile ProfileInfo

	c.OnHTML("#user-profile-left", func(e *colly.HTMLElement) {
		Name := e.ChildText("h3")		
		UserName := e.ChildText("h4")		
		List := strings.Split(e.ChildText(":nth-child(7)"), ":")
		School := List[1]
		List = strings.Split(e.ChildText(":nth-child(6)"), ":")
		WorldRank := List[1]
		Profile = ProfileInfo{Name, UserName, School, WorldRank}
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.Visit(fmt.Sprintf("https://www.spoj.com/users/%s/", handle))

	return Profile
}


func GetSubmissions(handle string) []Submission {

	c := colly.NewCollector()
	var submissions []Submission
   
	c.OnHTML("tbody", func(e *colly.HTMLElement){
		e.ForEach("tr", func(_ int, elem *colly.HTMLElement) {
			Name := elem.ChildText(".sproblem a")
			URL := "https://www.spoj.com"+elem.ChildAttr(".sproblem a","href")
			CreationDate := elem.ChildText(".status_sm span")
			status := elem.ChildText(".statusres")
			language := elem.ChildText(".slang span")	
			submissions = append(submissions, Submission{Name, URL, CreationDate, status, language})
		})
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})
	
	c.Visit(fmt.Sprintf("https://www.spoj.com/status/%s/", handle))

	return submissions

}

func GetProblems(handle string) Problems {

	c := colly.NewCollector()
	var problems Problems
	c.OnHTML(".dl-horizontal", func(e *colly.HTMLElement){
		count := e.ChildText(":nth-child(2)")
		submitted := e.ChildText(":nth-child(4)")
		
		problems = Problems{count, submitted}
	})
	
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})
	
	c.Visit(fmt.Sprintf("https://www.spoj.com/users/%s/", handle))

	return problems
}

func GetSolvedProblems(handle string) []SolvedProblems{

	c := colly.NewCollector()
	var solved []SolvedProblems
	c.OnHTML("#user-profile-tables tr", func(e *colly.HTMLElement){
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

