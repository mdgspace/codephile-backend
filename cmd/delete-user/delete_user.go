package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/globalsign/mgo/bson"
	_ "github.com/mdg-iitr/Codephile/conf"

	"os"

	"github.com/mdg-iitr/Codephile/models"
	"github.com/mdg-iitr/Codephile/models/db"
	search "github.com/mdg-iitr/Codephile/services/elastic"
	"github.com/mdg-iitr/Codephile/services/firebase"

	"github.com/mdg-iitr/Codephile/services/auth"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run ./delete_user <uid>")
		os.Exit(1)
	}
	user, err := models.GetUser(bson.ObjectIdHex(os.Args[1]))
	if err != nil {
		panic(err)
	}
	if user.Picture != "" {
		// Delete profile pic
		pic_name := strings.Split(user.Picture, "/profile/")[1]
		err = firebase.DeleteObject("profile/" + pic_name)
		if err != nil {
			panic(err)
		}
	}
	// Delete user from elasticsearch index
	client := search.GetElasticClient()
	_, err = client.Delete().Index("codephile").Id(user.ID.Hex()).Do(context.TODO())
	if err != nil {
		panic(err)
	}
	// Delete user from our database
	sess := db.NewUserCollectionSession()
	defer sess.Close()
	coll := sess.Collection
	err = coll.RemoveId(user.ID)
	if err != nil {
		panic(err)
	}
	// Lastly block all the issued tokens
	err = auth.BlacklistUser(user.ID)
	if err != nil {
		panic(err)
	}
	fmt.Println("Success")
}
