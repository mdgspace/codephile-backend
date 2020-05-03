package models

import (
	"github.com/globalsign/mgo/bson"
	. "github.com/mdg-iitr/Codephile/errors"
	"github.com/mdg-iitr/Codephile/models/db"
	"github.com/mdg-iitr/Codephile/models/types"
)

func GetFollowingUsers(ID bson.ObjectId) ([]types.FollowingUser, error) {
	sess := db.NewUserCollectionSession()
	defer sess.Close()
	coll := sess.Collection
	var user types.User
	err := coll.FindId(ID).Select(bson.M{"followingUsers.f_id": 1}).One(&user)
	if err != nil {
		return nil, UserNotFoundError
	}
	followingUIDs := make([]bson.ObjectId, 0, len(user.FollowingUsers))
	for _, f := range user.FollowingUsers {
		followingUIDs = append(followingUIDs, f.ID)
	}
	var followingUsers []types.FollowingUser
	err2 := coll.Find(bson.M{"_id": bson.M{"$in": followingUIDs}}).Select(
		bson.M{"_id": 1, "username": 1, "picture": 1,
			"fullname": 1}).All(&followingUsers)

	if err2 != nil {
		return nil, err2
	}
	return followingUsers, nil
}

func UnFollowUser(uid1 bson.ObjectId, uid2 bson.ObjectId) error {
	sess := db.NewUserCollectionSession()
	defer sess.Close()
	coll := sess.Collection
	update := bson.M{
		"$pull": bson.M{
			"followingUsers": bson.M{
				"f_id": uid2,
			},
		},
	}
	return coll.UpdateId(uid1, update)
}

func FollowUser(uid1 bson.ObjectId, uid2 bson.ObjectId) error {
	//uid1 is of the person who wants to follow
	//uid2 is the person being followed
	user1, err1 := GetUser(uid1)
	user2, err2 := GetUser(uid2)
	if err1 != nil || err2 != nil {
		return UserNotFoundError
	}
	//add the uid2 in the database of uid1
	var following types.Following
	following.ID = user2.ID
	update := bson.M{"$addToSet": bson.M{"followingUsers": following}}
	collection := db.NewUserCollectionSession()
	defer collection.Close()
	return collection.Collection.UpdateId(user1.ID, update)
}
