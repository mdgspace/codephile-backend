package db

import (
	"log"
	"os"
	"path/filepath"

	"github.com/astaxie/beego"
	"github.com/getsentry/sentry-go"
	"github.com/mdg-iitr/Codephile/conf"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var maxPool int

// var usernameIndex = mgo.Index{
// 	Key:        []string{"username"},
// 	Unique:     true,
// 	Background: true,
// }

var emailIndex = mongo.IndexModel{
	Keys:        []string{"email"},
	Options:     options.Index().SetUnique(true).SetBackground(true),
}

var usernameIndex = mongo.IndexModel{
	Keys:        []string{"username"},
	Options:     options.Index().SetUnique(true).SetBackground(true),
}

func init() {
	var err error
	maxPool, err = beego.AppConfig.Int("DBMaxPool")
	if err != nil {
		if err := beego.LoadAppConfig("ini", filepath.Join(conf.AppRootDir, "conf/app.conf")); err != nil {
			panic(err)
		}
		maxPool = beego.AppConfig.DefaultInt("DBMaxPool", 30)
	}
	// init method to start db
	checkAndInitServiceConnection()
	c := NewUserCollectionSession()
	//TODO: ensure that indexes are present
	defer c.Close()
	if err != nil {
		sentry.CurrentHub().CaptureException(err)
		log.Println(err.Error())
	}

}

func checkAndInitServiceConnection() {
	if service.baseClient == nil {
		service.URL = os.Getenv("DBPath")
		err := service.New()
		if err != nil {
			panic(err)
		}
	}
}
