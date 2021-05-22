package main

import (
	"context"
	"fmt"
	"github.com/globalsign/mgo/bson"
	_ "github.com/mdg-iitr/Codephile/conf"
	"github.com/mdg-iitr/Codephile/models/db"
	"github.com/mdg-iitr/Codephile/models/types"
	search "github.com/mdg-iitr/Codephile/services/elastic"
	"github.com/olivere/elastic/v7"
	"log"
)

func main() {
	sess := db.NewUserCollectionSession()
	defer sess.Close()
	coll := sess.Collection
	iter := coll.Find(nil).Select(bson.M{"_id": 1,
		"username": 1, "fullname": 1, "institute": 1, "picture": 1,
		"handle": 1}).Iter()
	var user types.User
	client := search.GetElasticClient()
	resp, err := client.DeleteIndex("codephile").Do(context.Background())
	if err != nil {
		log.Println(resp)
		return
	}
	fmt.Println(resp.Acknowledged)
	bulkRequest := client.Bulk()
	for iter.Next(&user) {
		req := elastic.NewBulkIndexRequest().Index("codephile").Doc(types.SearchDoc{
			ID:        user.ID,
			Username:  user.Username,
			FullName:  user.FullName,
			Institute: user.Institute,
			Picture:   user.Picture,
			Handle:    user.Handle,
		}).Id(user.ID.Hex())
		bulkRequest = bulkRequest.Add(req)
	}
	bulkResponse, err := bulkRequest.Do(context.Background())
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("Total documents created: %d\nTotal documents deleted: %d\nTotal documents succeded: %d\nTotal documents failed: %d\n",
		len(bulkResponse.Created()),
		len(bulkResponse.Deleted()),
		len(bulkResponse.Succeeded()),
		len(bulkResponse.Failed()))
}
