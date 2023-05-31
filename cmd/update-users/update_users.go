package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/mdg-iitr/Codephile/conf"
	"github.com/mdg-iitr/Codephile/models"
	"github.com/mdg-iitr/Codephile/models/db"
	"github.com/mdg-iitr/Codephile/models/types"
	"log"
	"time"
)

// updates the submission and profile of all the users

func main() {
	sess := db.NewUserCollectionSession()
	defer sess.Close()
	coll := sess.Collection
	cursor, _ := coll.Find(context.TODO(), bson.M{}, options.Find().SetProjection(bson.M{"_id": 1, "verified": 1}))
	var user types.User
	for cursor.Next(context.TODO()) {
		cursor.Decode(&user)
		if !user.Verified {
			continue
		}
		for _, site := range conf.ValidSites {
			err := models.AddSubmissions(user.ID, site, context.Background())
			if err != nil {
				log.Println(err.Error())
			}
			err = models.AddOrUpdateProfile(user.ID, site, context.Background())
			if err != nil {
				log.Println(err.Error())
			}
		}
		time.Sleep(time.Second)
	}
}
