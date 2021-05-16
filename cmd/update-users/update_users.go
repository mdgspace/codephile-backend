package main

import (
	"context"
	"github.com/globalsign/mgo/bson"
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
	iter := coll.Find(nil).Select(bson.M{"_id": 1}).Iter()
	var user types.User
	for iter.Next(&user) {
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
