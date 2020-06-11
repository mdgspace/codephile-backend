package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	r "github.com/go-redis/redis"
	"github.com/mdg-iitr/Codephile/models/types"
	"github.com/mdg-iitr/Codephile/services/redis"
)

func ReturnContests() (types.Result, error) {
	return contestsFromCache()
}

func ReturnSpecificContests(site string) (types.Result, error) {
	initialResult, err := contestsFromCache()
	if err != nil {
		// handle error
		return types.Result{}, err
	}
	//initialResult stores all the contests
	var finalResult types.Result //finalResult will store the website's contests only

	//looping over all the ongoing contests and selecting only those specific to the website
	for _, v := range initialResult.Ongoing {
		if strings.ToLower(v.Platform) == site {
			finalResult.Ongoing = append(finalResult.Ongoing, v)
		}
	}
	//looping over all the upcoming contests and selecting only those specific to the website
	for _, v := range initialResult.Upcoming {
		if strings.ToLower(v.Platform) == site {
			finalResult.Upcoming = append(finalResult.Upcoming, v)
		}
	}
	//equating the timestamp
	finalResult.Timestamp = initialResult.Timestamp
	return finalResult, nil
}

func contestsFromCache() (types.Result, error) {
	var result types.Result
	client := redis.GetRedisClient()
	err := client.Get("contest").Scan(&result)
	if err == r.Nil {
		log.Println("cache miss")
		result, err = updateCache()
		if err != nil {
			return types.Result{}, err
		}
	} else if err != nil {
		return types.Result{}, err
	}
	return result, nil
}

func updateCache() (types.Result, error) {
	data := fetchFromWeb()
	var s map[string]types.Result
	err := json.Unmarshal(data, &s)
	if err != nil {
		//error in unmarshalling
		return types.Result{}, err
	}
	result := s["result"]
	sort.Slice(result.Upcoming, func(i, j int) bool {
		time1 := result.Upcoming[i].StartTime.Time
		time2 := result.Upcoming[j].StartTime.Time
		diff := time2.Sub(time1).Seconds()
		return diff > 0.0
	})
	sort.Slice(result.Ongoing, func(i, j int) bool {
		time1 := result.Ongoing[i].EndTime.Time
		time2 := result.Ongoing[j].EndTime.Time
		diff := time2.Sub(time1).Seconds()
		return diff > 0.0
	})
	client := redis.GetRedisClient()
	_, err = client.Set("contest", result, time.Hour).Result()
	if err != nil {
		return types.Result{}, err
	}
	return result, nil
}

func fetchFromWeb() (data []byte) {
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
