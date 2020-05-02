package models

import (
	"errors"
	"fmt"
	"github.com/globalsign/mgo/bson"
	. "github.com/mdg-iitr/Codephile/conf"
	. "github.com/mdg-iitr/Codephile/errors"
	"github.com/mdg-iitr/Codephile/models/db"
	"github.com/mdg-iitr/Codephile/models/types"
	"github.com/mdg-iitr/Codephile/scripts"
)

func AddOrUpdateProfile(uid bson.ObjectId, site string) error {
	sess := db.NewUserCollectionSession()
	defer sess.Close()
	coll := sess.Collection
	var user types.User
	err := coll.FindId(uid).Select(bson.M{"handle": 1}).One(&user)
	if err != nil {
		//handle the error (Invalid user)
		return UserNotFoundError
	}
	var userProfile types.ProfileInfo
	//runs code to fetch the particular script's getProfile function
	switch site {
	case CODECHEF:
		handle := user.Handle.Codechef
		if handle == "" {
			return HandleNotFoundError
		}
		userProfile = scripts.GetCodechefProfileInfo(handle)
		accuracy, err := GetAccuracy(uid, CODECHEF)
		if err != nil {
			userProfile.Accuracy = ""
		} else {
			userProfile.Accuracy = accuracy
		}
	case CODEFORCES:
		handle := user.Handle.Codeforces
		if handle == "" {
			return HandleNotFoundError
		}
		userProfile = scripts.GetCodeforcesProfileInfo(handle)
		accuracy, err := GetAccuracy(uid, CODEFORCES)
		if err != nil {
			userProfile.Accuracy = ""
		} else {
			userProfile.Accuracy = accuracy
		}
	case SPOJ:
		handle := user.Handle.Spoj
		if handle == "" {
			return HandleNotFoundError
		}
		userProfile = scripts.GetSpojProfileInfo(handle)
		accuracy, err := GetAccuracy(uid, SPOJ)
		if err != nil {
			userProfile.Accuracy = ""
		} else {
			userProfile.Accuracy = accuracy
		}
	case HACKERRANK:
		handle := user.Handle.Hackerrank
		if handle == "" {
			return HandleNotFoundError
		}
		userProfile = scripts.GetHackerrankProfileInfo(handle)
		accuracy, err := GetAccuracy(uid, HACKERRANK)
		if err != nil {
			userProfile.Accuracy = ""
		} else {
			userProfile.Accuracy = accuracy
		}
	} // add a default case for non-existent website
	//Profile fetched. Store in database
	newNode := "profiles." + site + "Profile"
	return coll.UpdateId(uid, bson.M{"$set": bson.M{newNode: userProfile}})
}

func GetProfiles(ID bson.ObjectId) (types.AllProfiles, error) {
	coll := db.NewUserCollectionSession()
	defer coll.Close()
	user := types.User{}
	err := coll.Collection.FindId(ID).Select(bson.M{"profiles": 1}).One(&user)
	return user.Profiles, err
}

func CompareUser(uid1 bson.ObjectId, uid2 bson.ObjectId) (types.AllWorldRanks, error) {
	collection := db.NewUserCollectionSession()
	defer collection.Close()
	//gets the different profiles to fetch world ranks
	p1, err1 := GetProfiles(uid1)
	p2, err2 := GetProfiles(uid2)
	if err1 != nil || err2 != nil {
		return types.AllWorldRanks{},
			fmt.Errorf("Could not get user: %s\n%s", err1, err2)
	}

	return types.AllWorldRanks{
		CodechefWorldRanks: types.WorldRankComparison{
			WorldRank1: p1.CodechefProfile.WorldRank,
			WorldRank2: p2.CodechefProfile.WorldRank,
		},
		CodeforcesWorldRanks: types.WorldRankComparison{
			WorldRank1: p1.CodeforcesProfile.WorldRank,
			WorldRank2: p2.CodeforcesProfile.WorldRank,
		},
		HackerrankWorldRanks: types.WorldRankComparison{
			WorldRank1: p1.HackerrankProfile.WorldRank,
			WorldRank2: p2.HackerrankProfile.WorldRank,
		},
		SpojWorldRanks: types.WorldRankComparison{
			WorldRank1: p1.SpojProfile.WorldRank,
			WorldRank2: p2.SpojProfile.WorldRank,
		},
	}, nil

}

func getCorrectIncorrectCount(uid bson.ObjectId, websiteUrl string, correctSubmissionIdentifier string) (int, int, error) {
	sess := db.NewUserCollectionSession()
	defer sess.Close()
	coll := sess.Collection
	match := bson.M{"$match": bson.M{
		"_id": uid,
	}}
	unwind := bson.M{
		"$unwind": "$submissions",
	}
	match2 := bson.M{"$match": bson.M{
		"submissions.url": bson.M{"$regex": bson.RegEx{
			Pattern: "^" + websiteUrl,
		}},
	}}
	pipe := coll.Pipe([]bson.M{
		match,
		unwind,
		match2,
		bson.M{
			"$facet": bson.M{
				"total": []bson.M{bson.M{"$count": "total"}},
				"correct": []bson.M{
					bson.M{"$match": bson.M{"submissions.status": correctSubmissionIdentifier}},
					bson.M{"$count": "correct"}},
			},
		},
	})
	var result []map[string][]map[string]int
	err := pipe.All(&result)
	if err != nil || len(result) == 0 || len(result[0]["total"]) == 0 || len(result[0]["correct"]) == 0 {
		return 0, 1, errors.New("could not get accuracy")
	}
	return result[0]["correct"][0]["correct"], result[0]["total"][0]["total"], nil
}

// GetAccuracy function calculates the accuracy of a particular site and returns it
func GetAccuracy(uid bson.ObjectId, website string) (string, error) {
	switch website {
	case CODECHEF:
		correct, total, err := getCorrectIncorrectCount(uid, "http://www.codechef.com/", "AC")
		return fmt.Sprintf("%f", float64(correct)/float64(total)), err
	case CODEFORCES:
		correct, total, err := getCorrectIncorrectCount(uid, "http://codeforces.com/", "OK")
		return fmt.Sprintf("%f", float64(correct)/float64(total)), err
	case SPOJ:
		correct, total, err := getCorrectIncorrectCount(uid, "https://www.spoj.com", "accepted")
		return fmt.Sprintf("%f", float64(correct)/float64(total)), err
	case HACKERRANK:
		return "1", nil
	default:
		return "", errors.New("Invalid Website")
	}
}
