package firebase

import (
	gcpStorage "cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"errors"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/storage"
	"github.com/astaxie/beego"
	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

var client *storage.Client
var conf map[string]string

func init() {
	err := json.Unmarshal([]byte(os.Getenv("FIREBASE_CONFIG")), &conf)
	if err != nil {
		log.Println("bad firebase configuration")
	}
	opt := option.WithCredentialsJSON([]byte(os.Getenv("FIREBASE_CREDENTIALS")))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		sentry.CaptureException(err)
		log.Println(err.Error())
		return
	}
	client, err = app.Storage(context.Background())
	if err != nil {
		sentry.CaptureException(err)
		log.Println(err.Error())
	}
}

func URLFromName(name string) string {
	return "https://storage.googleapis.com/" + conf["storageBucket"] + "/profile/" + name
}

func AddFile(f multipart.File, fh *multipart.FileHeader, oldPic string) (string, error) {
	if client == nil {
		return "", errors.New("firebase conf not available")
	}
	bucket, _ := client.DefaultBucket()
	publicURL := "https://storage.googleapis.com/" + conf["storageBucket"] + "/"
	if bucket == nil {
		err := errors.New("nil bucket")
		log.Println(err)
		return "", err
	}
	oldPicName := strings.Split(oldPic, publicURL)[1]
	isSpecialPic := false
	specialPics := beego.AppConfig.DefaultStrings("DEFAULT_PICS", []string{})
	for _, pic := range specialPics {
		if oldPicName == "profile/"+pic {
			isSpecialPic = true
			break
		}
	}
	if oldPic != "" && !isSpecialPic {
		err := bucket.Object(oldPicName).Delete(context.Background())
		if err != nil {
			log.Println(err)
		}
	}
	// random filename, retaining existing extension.
	name := "profile/" + uuid.New().String() + path.Ext(fh.Filename)
	w := bucket.Object(name).NewWriter(context.Background())
	w.ACL = []gcpStorage.ACLRule{{Entity: gcpStorage.AllUsers, Role: gcpStorage.RoleReader}}
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

func DeleteObject(key string) error {
	bucket, _ := client.DefaultBucket()
	return bucket.Object(key).Delete(context.Background())
}
