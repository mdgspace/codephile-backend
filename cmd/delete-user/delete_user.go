package delete_user

import (
	"context"
	"fmt"
	"github.com/globalsign/mgo/bson"
	_ "github.com/mdg-iitr/Codephile/conf"

	"github.com/mdg-iitr/Codephile/models"
	"github.com/mdg-iitr/Codephile/models/db"
	search "github.com/mdg-iitr/Codephile/services/elastic"
	"github.com/mdg-iitr/Codephile/services/firebase"
	"os"

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
	// Delete profile pic
	err = firebase.DeleteObject(user.Picture)
	if err != nil {
		panic(err)
	}
	// Delete user from elasticsearch index
	client := search.GetElasticClient()
	_, err = client.Delete().Index("codephile").Id(user.ID.Hex()).Do(context.TODO())
	if err != nil {
		panic(err)
	}
	// Delete user from our database
	sess := db.NewUserCollectionSession()
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
