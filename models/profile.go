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
	user, err := GetUser(uid)
	if err != nil {
		//handle the error (Invalid user)
		return UserNotFoundError
	}
	var UserProfile types.ProfileInfo
	//runs code to fetch the particular script's getProfile function
	switch site {
	case CODECHEF:
		UserProfile = scripts.GetCodechefProfileInfo(user.Handle.Codechef)
		accuracy, err := GetAccuracy(user, CODECHEF)
		if err != nil {
			UserProfile.Accuracy = ""
		} else {
			UserProfile.Accuracy = accuracy
		}
	case CODEFORCES:
		UserProfile = scripts.GetCodeforcesProfileInfo(user.Handle.Codeforces)
		accuracy, err := GetAccuracy(user, CODEFORCES)
		if err != nil {
			UserProfile.Accuracy = ""
		} else {
			UserProfile.Accuracy = accuracy
		}
	case SPOJ:
		UserProfile = scripts.GetSpojProfileInfo(user.Handle.Spoj)
		accuracy, err := GetAccuracy(user, SPOJ)
		if err != nil {
			UserProfile.Accuracy = ""
		} else {
			UserProfile.Accuracy = accuracy
		}
	case HACKERRANK:
		UserProfile = scripts.GetHackerrankProfileInfo(user.Handle.Hackerrank)
		accuracy, err := GetAccuracy(user, HACKERRANK)
		if err != nil {
			UserProfile.Accuracy = ""
		} else {
			UserProfile.Accuracy = accuracy
		}
	} // add a default case for non-existent website
	//Profile fetched. Store in database
	var ProfileTobeInserted types.Profile
	ProfileTobeInserted.Website = site
	ProfileTobeInserted.Profileinfo = UserProfile
	// ProfileTobeInserted is all set to be put in the database
	collection := db.NewUserCollectionSession()
	defer collection.Close()
	// err2 := collection.Collection.Update(bson.D{{"_id" , user.ID}},bson.D{{"$set" , ProfileTobeInserted}})
	NewNode := site + "Profile"
	SelectedUser := bson.D{{Name: "_id", Value: user.ID}}
	Update := bson.D{{Name: "$set", Value: bson.D{{Name: NewNode, Value: ProfileTobeInserted}}}}
	_, err2 := collection.Collection.Upsert(SelectedUser, Update)
	//inserted into the document
	if err2 == nil {
		return nil
	}
	return err2
}

func GetProfiles(ID bson.ObjectId) (types.AllProfiles, error) {
	coll := db.NewUserCollectionSession()
	defer coll.Close()
	var profiles types.AllProfiles
	var profilesToBeReturned types.AllProfiles //appends the profile to this variable which will be returned
	err1 := coll.Collection.FindId(ID).Select(bson.M{"codechefProfile": 1}).One(&profiles)
	profilesToBeReturned.CodechefProfile = profiles.CodechefProfile
	err2 := coll.Collection.FindId(ID).Select(bson.M{"codeforcesProfile": 1}).One(&profiles)
	profilesToBeReturned.CodeforcesProfile = profiles.CodeforcesProfile
	err3 := coll.Collection.FindId(ID).Select(bson.M{"hackerrankProfile": 1}).One(&profiles)
	profilesToBeReturned.HackerrankProfile = profiles.HackerrankProfile
	err4 := coll.Collection.FindId(ID).Select(bson.M{"spojProfile": 1}).One(&profiles)
	profilesToBeReturned.SpojProfile = profiles.SpojProfile
	if err1 == nil && err2 == nil && err3 == nil && err4 == nil {
		return profilesToBeReturned, nil
	} else {
		if err1 != nil {
			return profilesToBeReturned, err1
		} else if err2 != nil {
			return profilesToBeReturned, err2
		} else if err3 != nil {
			return profilesToBeReturned, err3
		} else {
			return profilesToBeReturned, err4
		}
	}
}

func CompareUser(uid1 bson.ObjectId, uid2 bson.ObjectId) (types.AllWorldRanks, error) {
	collection := db.NewCollectionSession("coduser")
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
			WorldRank1: p1.CodechefProfile.Profileinfo.WorldRank,
			WorldRank2: p2.CodechefProfile.Profileinfo.WorldRank,
		},
		CodeforcesWorldRanks: types.WorldRankComparison{
			WorldRank1: p1.CodeforcesProfile.Profileinfo.WorldRank,
			WorldRank2: p2.CodeforcesProfile.Profileinfo.WorldRank,
		},
		HackerrankWorldRanks: types.WorldRankComparison{
			WorldRank1: p1.HackerrankProfile.Profileinfo.WorldRank,
			WorldRank2: p2.HackerrankProfile.Profileinfo.WorldRank,
		},
		SpojWorldRanks: types.WorldRankComparison{
			WorldRank1: p1.SpojProfile.Profileinfo.WorldRank,
			WorldRank2: p2.SpojProfile.Profileinfo.WorldRank,
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
	if err != nil || len(result) == 0 || len(result[0]["total"]) == 0 {
		return 0, 1, errors.New("could not get accuracy")
	}
	return result[0]["correct"][0]["correct"], result[0]["total"][0]["total"], nil
}

// GetAccuracy function calculates the accuracy of a particular site and returns it
func GetAccuracy(user *types.User, website string) (string, error) {
	switch website {
	case CODECHEF:
		correct, total, err := getCorrectIncorrectCount(user.ID, "http://www.codechef.com/", "AC")
		return fmt.Sprintf("%f", float64(correct)/float64(total)), err
	case CODEFORCES:
		correct, total, err := getCorrectIncorrectCount(user.ID, "http://codeforces.com/", "OK")
		return fmt.Sprintf("%f", float64(correct)/float64(total)), err
	case SPOJ:
		correct, total, err := getCorrectIncorrectCount(user.ID, "https://www.spoj.com", "accepted")
		return fmt.Sprintf("%f", float64(correct)/float64(total)), err
	case HACKERRANK:
		return "1", nil
	default:
		return "", errors.New("Invalid Website")
	}
}
