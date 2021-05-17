package db

import (
	"github.com/astaxie/beego"
	"github.com/getsentry/sentry-go"
	"github.com/globalsign/mgo"
	"github.com/mdg-iitr/Codephile/conf"
	"log"
	"os"
	"path/filepath"
)

var maxPool int

var usernameIndex = mgo.Index{
	Key:        []string{"username"},
	Unique:     true,
	Background: true,
}

var emailIndex = mgo.Index{
	Key:        []string{"email"},
	Unique:     true,
	Background: true,
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
	err = c.Collection.EnsureIndex(usernameIndex)
	if err != nil {
		log.Println(err.Error())
		sentry.CurrentHub().CaptureException(err)
	}
	err = c.Collection.EnsureIndex(emailIndex)
	if err != nil {
		log.Println(err.Error())
		sentry.CurrentHub().CaptureException(err)
	}
	defer c.Close()
	if err != nil {
		sentry.CurrentHub().CaptureException(err)
		log.Println(err.Error())
	}

}

func checkAndInitServiceConnection() {
	if service.baseSession == nil {
		service.URL = os.Getenv("DBPath")
		err := service.New()
		if err != nil {
			panic(err)
		}
	}
}
