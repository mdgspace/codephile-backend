package models

import (
	"errors"
	"github.com/mdg-iitr/Codephile/models/types"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/globalsign/mgo/bson"
)

var ErrGeneric = errors.New("Feed is not absolutely correct")

func ReturnFeedContests() (types.S, error) {
	UnsortedContests, err := ReturnContests()
	if err != nil {
		return types.S{}, err
	}
	SortedContests := SortContests(UnsortedContests)
	return SortedContests, nil
}

func ReturnFeedFriends(uid bson.ObjectId) ([]types.FeedObject, error) {
	user, err := GetUser(uid)
	if err != nil {
		//user invalid return error
		log.Println("Invalid user")
		return nil, err
	}
	UserMissing := false
	UsernameError := false
	var Feed []types.FeedObject
	followingUsers, err2 := GetFollowingUsers(user.ID)
	if err2 != nil {
		return Feed, err2
	}
	for _, value := range followingUsers {
		UserSubmissions, err1 := GetSubmissions(value.ID)
		feedUser, err3 := GetUser(value.ID)
		user_name := feedUser.Username
		if err3 != nil {
			//some alteration in feed
			//this error will rarely occur
			UsernameError = true
		}
		if err1 != nil {
			//unable to fetch this user (feed will not be consisting this user's activity)
			log.Println("Unable to fetch a user for feed obtain")
			//handle error
			UserMissing = true
			continue
		} else {
			//user fetched
			for _, submission := range UserSubmissions.Codechef {
				var feedObject types.FeedObject
				feedObject.UserName = user_name
				feedObject.Name = submission.Name
				feedObject.URL = submission.URL
				feedObject.CreationDate = submission.CreationDate
				feedObject.Status = submission.Status
				feedObject.Points = submission.Points
				feedObject.Tags = submission.Tags
				Feed = append(Feed, feedObject)
			}
			for _, submission := range UserSubmissions.Codeforces {
				var feedObject types.FeedObject
				feedObject.UserName = user_name
				feedObject.Name = submission.Name
				feedObject.URL = submission.URL
				feedObject.CreationDate = submission.CreationDate
				feedObject.Status = submission.Status
				feedObject.Points = strconv.Itoa(submission.Points)
				feedObject.Tags = submission.Tags
				feedObject.Rating = submission.Rating
				Feed = append(Feed, feedObject)
			}
			for _, submission := range UserSubmissions.Spoj {
				var feedObject types.FeedObject
				feedObject.UserName = user_name
				feedObject.Name = submission.Name
				feedObject.URL = submission.URL
				feedObject.CreationDate = submission.CreationDate
				feedObject.Status = submission.Status
				feedObject.Points = strconv.Itoa(submission.Points)
				feedObject.Language = submission.Language
				feedObject.Tags = submission.Tags
				Feed = append(Feed, feedObject)
			}
			for _, submission := range UserSubmissions.Hackerrank {
				var feedObject types.FeedObject
				feedObject.UserName = user_name
				feedObject.Name = submission.Name
				feedObject.URL = submission.URL
				feedObject.CreationDate = submission.CreationDate
				Feed = append(Feed, feedObject)
			}
		}
	}
	sort.Slice(Feed, func(i, j int) bool {
		return Feed[i].CreationDate.After(Feed[j].CreationDate)
	})
	if UserMissing == true {
		return Feed, ErrGeneric
	} else if UsernameError == true {
		return Feed, ErrGeneric
	}
	return Feed, nil
}

//SortContests to sort contests according to StartTime and EndTime
func SortContests(contests types.S) types.S {
	// sorting the ongoing contests
	var (
		n          = len(contests.Result.Ongoing)
		sorted     = false
		timeLayout = "Mon, 2 Jan 2006 15:04"
	)
	for !sorted {
		swapped := false
		for i := 0; i < n-1; i++ {
			time1, _ := time.Parse(timeLayout, contests.Result.Ongoing[i].EndTime)
			time2, _ := time.Parse(timeLayout, contests.Result.Ongoing[i+1].EndTime)
			diff := time2.Sub(time1).Seconds()
			if diff < 0.0 {
				//swap objects
				contests.Result.Ongoing[i+1], contests.Result.Ongoing[i] = contests.Result.Ongoing[i], contests.Result.Ongoing[i+1]
				swapped = true
			}
		}
		if !swapped {
			sorted = true
		}
		n = n - 1
	}
	//sorting the upcoming contests

	n = len(contests.Result.Upcoming)
	sorted = false

	for !sorted {
		swapped := false
		for i := 0; i < n-1; i++ {
			time1, _ := time.Parse(timeLayout, contests.Result.Upcoming[i].StartTime)
			time2, _ := time.Parse(timeLayout, contests.Result.Upcoming[i+1].StartTime)
			diff := time2.Sub(time1).Seconds()
			if diff < 0.0 {
				//swap objects
				contests.Result.Ongoing[i+1], contests.Result.Ongoing[i] = contests.Result.Ongoing[i], contests.Result.Ongoing[i+1]
				swapped = true
			}
		}
		if !swapped {
			sorted = true
		}
		n = n - 1
	}
	return contests
}
