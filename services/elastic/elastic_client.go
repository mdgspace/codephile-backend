package search

import (
	"github.com/olivere/elastic/v7"
	"log"
	"os"
)

var client *elastic.Client

func GetElasticClient() *elastic.Client {
	if client == nil {
		var err error
		client, err = elastic.NewClient(
			elastic.SetSniff(false),
			elastic.SetBasicAuth(os.Getenv("ELASTICUSERNAME"), os.Getenv("ELASTICPASSWORD")),
			elastic.SetURL(os.Getenv("ELASTICCLOUDURL")))
		if err != nil {
			log.Println(err.Error())
		}
	}
	return client
}
