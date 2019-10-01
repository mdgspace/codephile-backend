package models  //try package models if import cycle error

import (
	// "errors"
	"github.com/globalsign/mgo/bson"
	"github.com/mdg-iitr/Codephile/models/db"
	"github.com/mdg-iitr/Codephile/models"
)

type Following struct{
	CodephileHandle string          `bson:"codephile_handle"`
	ID              bson.ObjectId   `bson:"id"`
}

func FollowUser(uid1 string, uid2 string) error{
	//uid1 is of the person who wants to follow
	//uid2 is the person being followed
     if uid1 != "" && uid2 != "" && bson.IsObjectIdHex(uid1) && bson.IsObjectIdHex(uid2) {
			user1 , err1 := models.GetUser(bson.ObjectIdHex(uid1))
			user2 , err2 := models.GetUser(bson.ObjectIdHex(uid2))
            if err1 == nil && err2 == nil {
				//add the uid2 in the database of uid1
				collection := db.NewCollectionSession("coduser")
				defer collection.Close()
				var following Following
				following.ID = user2.ID
				following.CodephileHandle = user2.Username
				SelectedUser := bson.D{{"_id", user1.ID}}
				Update := bson.D{{"$addToSet" , bson.D{{"followingUsers" , following}}}}
				_, err := collection.Session.Upsert(SelectedUser,Update)
				return err
			} else {
				//unable to get the user from database
                return err1
			}
	 } else {
		 //uid is not valid
		 return nil	
	 }
}

