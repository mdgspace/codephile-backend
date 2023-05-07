package models

import (
	"context"

	. "github.com/mdg-iitr/Codephile/errors"
	"github.com/mdg-iitr/Codephile/models/db"
	"github.com/mdg-iitr/Codephile/models/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetFollowingUsers(ID primitive.ObjectID) ([]types.FollowingUser, error) {
	sess := db.NewUserCollectionSession()
	defer sess.Close()
	coll := sess.Collection
	var user types.User
	err := coll.FindOne(context.TODO(), bson.M{"_id": ID}, options.FindOne().SetProjection(bson.M{"followingUsers.f_id": 1})).Decode(&user)
	if err != nil {
		return nil, UserNotFoundError
	}
	followingUIDs := make([]primitive.ObjectID, 0, len(user.FollowingUsers))
	for _, f := range user.FollowingUsers {
		followingUIDs = append(followingUIDs, f.ID)
	}
	var followingUsers []types.FollowingUser
	cursor, err2 := coll.Find(context.TODO(), bson.M{"_id": bson.M{"$in": followingUIDs}}, options.Find().SetProjection(
		bson.M{"_id": 1, "username": 1, "picture": 1,
			"fullname": 1}))
	err2 = cursor.All(context.TODO(), &followingUsers)
	if err2 != nil {
		return nil, err2
	}
	return followingUsers, nil
}

func UnFollowUser(uid1 primitive.ObjectID, uid2 primitive.ObjectID) error {
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
	_, err := coll.UpdateByID(context.TODO(), uid1, update)
	return err
}

func FollowUser(uid1 primitive.ObjectID, uid2 primitive.ObjectID) error {
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
	_, err := collection.Collection.UpdateByID(context.TODO(), user1.ID, update)
	return err
}
