package search

import (
	"github.com/getsentry/sentry-go"
	"github.com/olivere/elastic/v7"
	"github.com/olivere/elastic/v7/config"
	"log"
	"os"
)

var client *elastic.Client

func GetElasticClient() *elastic.Client {
	if client == nil {
		var err error
		conf, err := config.Parse(os.Getenv("ELASTICURL"))
		if err != nil {
			panic("invalid elastic search connection string")
		}
		client, err = elastic.NewClientFromConfig(conf)
		if err != nil {
			sentry.CaptureException(err)
			log.Println(err.Error())
		}
	}
	return client
}
