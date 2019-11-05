package models 

import (
	"errors"
	"github.com/globalsign/mgo/bson"
	"github.com/mdg-iitr/Codephile/models"
	"strings"
	"strconv"
	"log"
	"time"
	"sort"
)

type FeedObject struct {
	UserName     string    `bson:"username"`
	Name         string    `bson:"name"`
	URL          string    `bson:"url"`
	CreationDate time.Time `bson:"creation_date"`
	Status       string    `bson:"status"`
	Language     string    `bson:"language"`
	Points       string       `bson:"points"`
	Tags         []string  `bson:"tags"`
	Rating       int       `bson:"rating"`
}

var ErrGeneric = errors.New("Feed is not absolutely correct")

func ReturnFeedContests() (models.S, error) {
	 UnsortedContests, err := models.ReturnContests()
	 if err != nil {
		 return models.S{}, err
	 } 
	 SortedContests := SortContests(UnsortedContests)
	 return SortedContests, nil
} 

func ReturnFeedFriends(uid bson.ObjectId) ([]FeedObject , error) {
	 user , err := models.GetUser(uid)
	 if err != nil{ 
		//user invalid return error
		log.Println("Invalid user")
		return nil, err
	 }
	 UserMissing := false
	 UsernameError := false
	 var Feed []FeedObject
	 followingUsers , err2 := models.GetFollowingUsers(user.ID)
	 if err2 != nil {
		 return Feed , err2
	 }
	 for _ , value := range followingUsers {
		UserSubmissions , err1 := models.GetSubmissions(value.ID)
		feedUser ,err3 := models.GetUser(value.ID)
		user_name := feedUser.Username
		if err3 != nil{
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
			for _ , submission := range UserSubmissions.Codechef {
				var feedObject FeedObject
				feedObject.UserName     = user_name 
				feedObject.Name         = submission.Name
				feedObject.URL          = submission.URL
				feedObject.CreationDate = submission.CreationDate
				feedObject.Status       = submission.Status
				feedObject.Points       = submission.Points
				feedObject.Tags         = submission.Tags
				Feed = append(Feed,feedObject)
			}
			for _ , submission := range UserSubmissions.Codeforces {
				var feedObject FeedObject
				feedObject.UserName     = user_name
				feedObject.Name         = submission.Name
				feedObject.URL          = submission.URL
				feedObject.CreationDate = submission.CreationDate
				feedObject.Status       = submission.Status
				feedObject.Points       = strconv.Itoa(submission.Points)
				feedObject.Tags         = submission.Tags
				feedObject.Rating       = submission.Rating
				Feed = append(Feed,feedObject)
			}
			for _ , submission := range UserSubmissions.Spoj {
				var feedObject FeedObject
				feedObject.UserName     = user_name
				feedObject.Name         = submission.Name
				feedObject.URL          = submission.URL
				feedObject.CreationDate = submission.CreationDate
				feedObject.Status       = submission.Status
				feedObject.Points       = strconv.Itoa(submission.Points)
				feedObject.Language     = submission.Language
				feedObject.Tags         = submission.Tags
				Feed = append(Feed,feedObject)
			}
			for _ , submission := range UserSubmissions.Hackerrank {
				var feedObject FeedObject
				feedObject.UserName     = user_name
				feedObject.Name         = submission.Name
				feedObject.URL          = submission.URL
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
	 return Feed , nil	 			
}

func SortContests(contests models.S) (models.S) {
	// sorting the ongoing contests 
	var (
        n = len(contests.Result.Ongoing)
        sorted = false
    )
    for !sorted {
        swapped := false
        for i := 0; i < n-1; i++ {
			splitResult1 := strings.Split(contests.Result.Ongoing[i].EndTime," ")
			splitResult2 := strings.Split(contests.Result.Ongoing[i+1].EndTime," ")
			if !CheckContest(splitResult1, splitResult2){
				contests.Result.Ongoing[i+1], contests.Result.Ongoing[i] = contests.Result.Ongoing[i] , contests.Result.Ongoing[i+1]
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
			splitResult1 := strings.Split(contests.Result.Upcoming[i].StartTime," ")
			splitResult2 := strings.Split(contests.Result.Upcoming[i+1].StartTime," ")
			if !CheckContest(splitResult1, splitResult2){
				contests.Result.Upcoming[i+1], contests.Result.Upcoming[i] = contests.Result.Upcoming[i] , contests.Result.Upcoming[i+1]
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

func CheckContest(result1 []string, result2 []string) bool {
	var Num1, Num2 int

	Date1,_  := strconv.Atoi(result1[1])
	Month1 := result1[2]
	Year1,_  := strconv.Atoi(result1[3])
	Time1  := result1[4]

	Date2,_  := strconv.Atoi(result2[1])
	Month2 := result2[2]
	Year2,_  := strconv.Atoi(result2[3])
	Time2  := result2[4]
	
	switch Month1 {
	case "Jan":
			Num1 = 1
	case "Feb":
			Num1 = 2
	case "Mar":
			Num1 = 3
	case "Apr":
			Num1 = 4
	case "May":
			Num1 = 5
	case "Jun":
			Num1 = 6
	case "Jul":
			Num1 = 7
	case "Aug":
			Num1 = 8
	case "Sep":
			Num1 = 9
	case "Oct":
			Num1 = 10
	case "Nov":
			Num1 = 11
	case "Dec":
			Num1 = 12
	}
	switch Month2 {
	case "Jan":
			Num2 = 1
	case "Feb":
			Num2 = 2
	case "Mar":
			Num2 = 3
	case "Apr":
			Num2 = 4
	case "May":
			Num2 = 5
	case "Jun":
			Num2 = 6
	case "Jul":
			Num2 = 7
	case "Aug":
			Num2 = 8
	case "Sep":
			Num2 = 9
	case "Oct":
			Num2 = 10
	case "Nov":
			Num2 = 11
	case "Dec":
			Num2 = 12
	}
	if Year2 > Year1 {
	  return true
	} else if Num2 > Num1 {
      return true
	} else if Date2 > Date1 {
	  return true
	} else if (strings.Split(Time2,":"))[0] > (strings.Split(Time1,":"))[0] {
	  return true
	} else if (strings.Split(Time2,":"))[1] > (strings.Split(Time1,":"))[1] {
	  return true	
	} 
    return false
}