package firebase

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"errors"
	firebase "firebase.google.com/go"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

var app *firebase.App
var conf map[string]string

func init() {
	err := json.Unmarshal([]byte(os.Getenv("FIREBASE_CONFIG")), &conf)
	if err != nil {
		log.Println("bad firebase configuration")
	}
	opt := option.WithCredentialsJSON([]byte(os.Getenv("FIREBASE_CREDENTIALS")))
	app, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Println(err.Error())
	}
}

func AddFile(f multipart.File, fh *multipart.FileHeader, oldPic string) (string, error) {
	client, err := app.Storage(context.Background())
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	bucket, _ := client.DefaultBucket()
	publicURL := "https://storage.googleapis.com/" + conf["storageBucket"] + "/"
	if bucket == nil {
		err := errors.New("nil bucket")
		log.Println(err)
		return "", err
	}
	if oldPic != "" {
		err := bucket.Object(strings.Split(oldPic, publicURL)[1]).Delete(context.Background())
		if err != nil {
			log.Println(err)
		}
	}
	// random filename, retaining existing extension.
	name := "profile/" + uuid.New().String() + path.Ext(fh.Filename)
	w := bucket.Object(name).NewWriter(context.Background())
	w.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	w.ContentType = fh.Header.Get("Content-Type")
	// Entries are immutable, be aggressive about caching (1 day).
	w.CacheControl = "public, max-age=86400"
	if _, err := io.Copy(w, f); err != nil {
		log.Println(err)
		return "", err
	}
	if err := w.Close(); err != nil {
		log.Println(err)
		return "", err
	}
	picUrl := publicURL + name
	return picUrl, nil
}
