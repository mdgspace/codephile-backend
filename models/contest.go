package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
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
	result, err := fetchFromWeb()
	if err != nil {
		return types.Result{}, err
	}
	client := redis.GetRedisClient()
	_, err = client.Set("contest", result, time.Minute).Result()
	if err != nil {
		return types.Result{}, err
	}
	return result, nil
}

func fetchFromWeb() (types.Result, error) {

	clistURL, _ := url.Parse("https://clist.by/api/v2/contest/")

	values := clistURL.Query()
	values.Set("host__regex", "codeforces.com|codechef.com|spoj.com|hackerrank.com|leetcode.com")
	values.Set("end__gte", time.Now().Format(time.RFC3339))
	values.Set("order_by", "start")
	values.Set("total_count", "true")
	clistURL.RawQuery = values.Encode()
	req, err := http.NewRequest(http.MethodGet, clistURL.String(), nil)
	if err != nil {
		return types.Result{}, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("ApiKey %s", os.Getenv("CLIST_KEY")))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return types.Result{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return types.Result{}, err
	}
	var clistResult types.CListResult
	err = json.Unmarshal(body, &clistResult)
	if err != nil {
		return types.Result{}, err
	}
	result, err := clistResult.ToResult()
	if err != nil {
		return types.Result{}, err
	}
	return result, nil
}
