package models

import (
	"errors"
	"github.com/globalsign/mgo/bson"
	"github.com/mdg-iitr/Codephile/models/db"
	"github.com/mdg-iitr/Codephile/models/types"
)

func GetFollowingUsers(ID bson.ObjectId) ([]types.Following, error) {
	coll := db.NewUserCollectionSession()
	defer coll.Close()
	var user types.User
	err := coll.Collection.FindId(ID).Select(bson.M{"followingUsers": 1}).One(&user)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user.FollowingUsers, nil
}

func FollowUser(uid1 bson.ObjectId, uid2 string) error {
	//uid1 is of the person who wants to follow
	//uid2 is the person being followed
	if uid2 != "" && bson.IsObjectIdHex(uid2) {
		user1, err1 := GetUser(uid1)
		user2, err2 := GetUser(bson.ObjectIdHex(uid2))
		if err1 == nil && err2 == nil {
			//add the uid2 in the database of uid1
			var following types.Following
			following.ID = user2.ID
			following.CodephileHandle = user2.Username
			update := bson.M{"$addToSet": bson.M{"followingUsers": following}}
			collection := db.NewUserCollectionSession()
			defer collection.Close()
			err := collection.Collection.UpdateId(user1.ID, update)
			return err
		} else {
			//unable to get the user from database
			return errors.New("Unable to fetch the user from the database")
		}
	} else {
		//uid is not valid
		return errors.New("UID Invalid")
	}
}
