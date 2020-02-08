package models

import (
	"errors"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/mdg-iitr/Codephile/conf"
	"github.com/mdg-iitr/Codephile/models/db"
	"github.com/mdg-iitr/Codephile/models/types"
	"github.com/mdg-iitr/Codephile/scripts"
)

func AddorUpdateProfile(uid bson.ObjectId, site string) (*types.User, error) {
	user, err := GetUser(uid)
	if err != nil {
		//handle the error (Invalid user)
		return nil, err
	}
	var UserProfile types.ProfileInfo
	//runs code to fetch the particular script's getProfile function
	switch site {
	case conf.CODECHEF:
		UserProfile = scripts.GetCodechefProfileInfo(user.Handle.Codechef)
		accuracy, err := GetAccuracy(user, conf.CODECHEF)
		if err != nil {
			UserProfile.Accuracy = ""
		} else {
			UserProfile.Accuracy = accuracy
		}
	case conf.CODEFORCES:
		UserProfile = scripts.GetCodeforcesProfileInfo(user.Handle.Codeforces)
		accuracy, err := GetAccuracy(user, conf.CODEFORCES)
		if err != nil {
			UserProfile.Accuracy = ""
		} else {
			UserProfile.Accuracy = accuracy
		}
	case conf.SPOJ:
		UserProfile = scripts.GetSpojProfileInfo(user.Handle.Spoj)
		accuracy, err := GetAccuracy(user, conf.SPOJ)
		if err != nil {
			UserProfile.Accuracy = ""
		} else {
			UserProfile.Accuracy = accuracy
		}
	case conf.HACKERRANK:
		UserProfile = scripts.GetHackerrankProfileInfo(user.Handle.Hackerrank)
		accuracy, err := GetAccuracy(user, conf.HACKERRANK)
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
	SelectedUser := bson.D{{"_id", user.ID}}
	Update := bson.D{{"$set", bson.D{{NewNode, ProfileTobeInserted}}}}
	_, err2 := collection.Collection.Upsert(SelectedUser, Update)
	//inserted into the document
	if err2 == nil {
		return user, nil
	} else {
		return nil, err2
	}
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

func CompareUser(uid1 bson.ObjectId, uid2 string) (types.AllWorldRanks, error) {
	var worldRanksComparison types.AllWorldRanks
	if uid2 != "" && bson.IsObjectIdHex(uid2) {
		//add the uid2 in the database of uid1
		collection := db.NewCollectionSession("coduser")
		defer collection.Close()
		//gets the different profiles to fetch world ranks
		profiles1, err1 := GetProfiles(uid1)
		profiles2, err2 := GetProfiles(bson.ObjectIdHex(uid2))

		//puts the world ranks in the struct fields
		worldRanksComparison.CodechefWorldRanks.WorldRank1 = profiles1.CodechefProfile.Profileinfo.WorldRank
		worldRanksComparison.CodechefWorldRanks.WorldRank2 = profiles2.CodechefProfile.Profileinfo.WorldRank

		worldRanksComparison.CodeforcesWorldRanks.WorldRank1 = profiles1.CodeforcesProfile.Profileinfo.WorldRank
		worldRanksComparison.CodeforcesWorldRanks.WorldRank2 = profiles2.CodeforcesProfile.Profileinfo.WorldRank

		worldRanksComparison.HackerrankWorldRanks.WorldRank1 = profiles1.HackerrankProfile.Profileinfo.WorldRank
		worldRanksComparison.HackerrankWorldRanks.WorldRank2 = profiles2.HackerrankProfile.Profileinfo.WorldRank

		worldRanksComparison.SpojWorldRanks.WorldRank1 = profiles1.SpojProfile.Profileinfo.WorldRank
		worldRanksComparison.SpojWorldRanks.WorldRank2 = profiles2.SpojProfile.Profileinfo.WorldRank

		//handle the errors
		if err1 != nil || err2 != nil {
			return worldRanksComparison, errors.New("Unable to fetch user from database")
		} else {
			return worldRanksComparison, nil
		}
	} else {
		//uid is not valid
		return worldRanksComparison, errors.New("UID Invalid")
	}
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
	case conf.CODECHEF:
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
	case conf.CODEFORCES:
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
	case conf.SPOJ:
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
	case conf.HACKERRANK:
		{
			//accuracy would be 100%
			return "100", nil
		}
	default:
		return "", errors.New("Invalid Website")
	}
}
