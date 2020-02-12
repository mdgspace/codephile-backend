package models

import (
	"encoding/json"
	"github.com/globalsign/mgo/bson"
	"github.com/mdg-iitr/Codephile/models/db"
	"github.com/mdg-iitr/Codephile/models/types"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func ReturnContests() (types.S, error) {
	data, err := fetchContests()
	if err != nil {
		return types.S{}, err
	}
	return data, nil
}

func ReturnSpecificContests(site string) (types.S, error) {
	InitialResult, err := fetchContests()
	if err != nil {
		// handle error 
		return types.S{}, err
	}
	//InitialResult stores all the contests
	var FinalResult types.S //FinalResult will store the website's contests only

	//looping over all the ongoing contests and selecting only those specific to the website
	for _, v := range InitialResult.Result.Ongoing {
		if strings.ToLower(v.Platform) == site {
			FinalResult.Result.Ongoing = append(FinalResult.Result.Ongoing, v)
		}
	}
	//looping over all the upcoming contests and selecting only those specific to the website
	for _, v := range InitialResult.Result.Upcoming {
		if strings.ToLower(v.Platform) == site {
			FinalResult.Result.Upcoming = append(FinalResult.Result.Upcoming, v)
		}
	}
	//equating the timestamp
	FinalResult.Result.Timestamp = InitialResult.Result.Timestamp
	return FinalResult, nil
}

func fetchContests() (types.S, error) {
	var Contests types.S
	var ContestsToBeReturned types.S

	collection := db.NewCollectionSession("contests")
	defer collection.Close()

	err := collection.Collection.Find(nil).Select(bson.M{"result": 1, "last_fetched": 1}).One(&Contests)
	if err != nil {
		//contests are not saved in the database
		ContestsWeb := FetchFromWeb()
		err := json.Unmarshal(ContestsWeb, &ContestsToBeReturned)
		if err != nil {
			//error in unmarshalling 
			return types.S{}, err
		}
		ContestsToBeReturned.LastFetched = time.Now()
		//save contests in database and return them
		_ = UpdateDatabase(ContestsToBeReturned)
		return ContestsToBeReturned, nil
	} else {
		//contests are saved in the database
		TimeLast := Contests.LastFetched
		Difference := time.Since(TimeLast)
		//Time difference since last call is greater than 1 hour
		if Difference.Hours() >= 1.0 {
			ContestsWeb := FetchFromWeb()
			err := json.Unmarshal(ContestsWeb, &ContestsToBeReturned)
			if err != nil {
				//error in unmarshalling 
				return types.S{}, err
			}
			//save contests in database and return
			ContestsToBeReturned.LastFetched = time.Now()
			_ = UpdateDatabase(ContestsToBeReturned)
			return ContestsToBeReturned, nil
		} else {
			//contests are to be returned from the database
			return Contests, nil
		}
	}
}

func FetchFromWeb() (data []byte) {
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

func UpdateDatabase(contests types.S) (error) {
	collection := db.NewCollectionSession("contests")
	defer collection.Close()

	Update := bson.D{{Name: "$set", Value: &contests}}
	_, err := collection.Collection.Upsert(nil, Update)
	return err
}
