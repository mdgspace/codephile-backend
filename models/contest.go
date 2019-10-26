package models

import(
	  "log"
	  "encoding/json"
	  "net/http"
	  "io/ioutil"
	  "strings"
)

type S struct {
	Result Result `json:"result"`
}

type Ongoing struct {
	EndTime       string `json:"EndTime"`
	Name          string `json:"Name"`
	Platform      string `json:"Platform"`
	ChallengeType string `json:"challenge_type,omitempty"`
	URL           string `json:"url"`
}

type Upcoming struct {
	Duration      string `json:"Duration"`
	EndTime       string `json:"EndTime"`
	Name          string `json:"Name"`
	Platform      string `json:"Platform"`
	StartTime     string `json:"StartTime"`
	URL           string `json:"url"`
	ChallengeType string `json:"challenge_type,omitempty"`
}

type Result struct {
	Ongoing   []Ongoing  `json:"ongoing"`
	Timestamp string     `json:"timestamp"`
	Upcoming  []Upcoming `json:"upcoming"`
}

func ReturnContests() S{
	data := fetchContests()
	var FinalResult S
	err := json.Unmarshal(data, &FinalResult)
	
	if err != nil {
		log.Println("Error")
		log.Fatal(err)
	}

    // log.Println(FinalResult.Result.Ongoing[0].Platform)
	// log.Println(len(FinalResult.Result.Ongoing))
	return FinalResult
}

func ReturnSpecificContests(site string) S {
	body := fetchContests()
	var InitialResult S  //InitialResult stores all the contests
	var FinalResult S    //FinalResult will store the website's contests only
	err := json.Unmarshal(body, &InitialResult)
	
	if err != nil {
		log.Println("Error")
		log.Fatal(err)
	}   
	//looping over all the ongoing contests and selecting only those specific to the website
		for _,v := range InitialResult.Result.Ongoing{
            if strings.ToLower(v.Platform) == site {
                FinalResult.Result.Ongoing = append(FinalResult.Result.Ongoing, v) 
			}
		}
	//looping over all the upcoming contests and selecting only those specific to the website
		for _,v := range InitialResult.Result.Upcoming{
            if strings.ToLower(v.Platform) == site {
				FinalResult.Result.Upcoming = append(FinalResult.Result.Upcoming, v)
			}
		}
		//equating the timestamp
		FinalResult.Result.Timestamp = InitialResult.Result.Timestamp
		return FinalResult
}

func fetchContests()(data []byte) {
	resp, err := http.Get("https://contesttrackerapi.herokuapp.com/")

	if err != nil {
		log.Println("Error")
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return
	}	
	return body
}


