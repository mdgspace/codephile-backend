package models

import (
	"errors"
	"github.com/mdg-iitr/Codephile/models/db"
	"github.com/mdg-iitr/Codephile/models/types"
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
	sess := db.NewUserCollectionSession()
	defer sess.Close()
	coll := sess.Collection
	var u types.User
	err := coll.FindId(uid).Select(bson.M{"followingUsers.f_id": 1}).One(&u)
	if err != nil {
		return nil, err
	}
	followingUID := make([]bson.ObjectId, 0, len(u.FollowingUsers))
	for _, f := range u.FollowingUsers {
		followingUID = append(followingUID, f.ID)
	}
	filter := bson.M{
		"$match": bson.M{
			"_id": bson.M{
				"$in": followingUID,
			},
		},
	}
	project := bson.M{
		"$project": bson.M{
			"_id":        1,
			"username":   1,
			"submission": "$submissions",
		},
	}
	unwind := bson.M{
		"$unwind": "$submission",
	}
	sort := bson.M{
		"$sort": bson.M{
			"submission.created_at": -1,
		},
	}
	pipe := coll.Pipe([]bson.M{
		filter,
		project,
		unwind,
		sort,
	}, )
	var res []types.FeedObject
	err = pipe.All(&res)
	//fmt.Println(res)
	return res, err
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
