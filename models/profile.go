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

func AddorUpdateProfile(uid bson.ObjectId, site string) (*types.User, error) {
	user, err := GetUser(uid)
	if err != nil {
		//handle the error (Invalid user)
		return nil, UserNotFoundError
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
		return user, nil
	}
	return nil, err2
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

// GetAccuracy function calculates the accuracy of a particular site and returns it
func GetAccuracy(user *types.User, website string) (string, error) {
	submissions, err := GetSubmissions(user.ID)

	var accuracy string

	if err != nil {
		return accuracy, err
	}

	var correctSubmissions float32
	var totalSubmissions float32

	switch website {
	case CODECHEF:
		{
			for _, value := range submissions.Codechef {
				totalSubmissions += 1.0
				if value.Status == "AC" {
					if value.Points == "100" {
						correctSubmissions += 1.0
					}
				}
			}
			accuracy = fmt.Sprintf("%f", correctSubmissions/totalSubmissions)
			return accuracy, nil
		}
	case CODEFORCES:
		{
			for _, value := range submissions.Codeforces {
				totalSubmissions += 1.0
				if value.Status == "OK" {
					correctSubmissions += 1.0
				}
			}
			accuracy = fmt.Sprintf("%f", correctSubmissions/totalSubmissions)
			return accuracy, nil
		}
	case SPOJ:
		{
			for _, value := range submissions.Spoj {
				totalSubmissions += 1.0
				if value.Status == "accepted" {
					correctSubmissions += 1.0
				}
			}
			accuracy = fmt.Sprintf("%f", correctSubmissions/totalSubmissions)
			return accuracy, nil
		}
	case HACKERRANK:
		{
			//accuracy would be 100%
			return "100", nil
		}
	default:
		return "", errors.New("Invalid Website")
	}
}
