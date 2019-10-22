package firebase

import (
	"cloud.google.com/go/storage"
	"context"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"log"
	"os"
)

var bucket *storage.BucketHandle

func GetStorageBucket() *storage.BucketHandle {
	if bucket == nil {

		opt := option.WithCredentialsJSON([]byte(os.Getenv("FIREBASE_CREDENTIALS")))
		app, err := firebase.NewApp(context.Background(), nil, opt)

		if err != nil {
			log.Println(err.Error())
		}
		client, err := app.Storage(context.Background())
		if err != nil {
			log.Println(err.Error())
		}
		bucket, err = client.DefaultBucket()
		if err != nil {
			log.Println(err.Error())
		}
	}
	return bucket
}
