package db

import (
	"github.com/astaxie/beego"
	"github.com/globalsign/mgo"
	"github.com/mdg-iitr/Codephile/conf"
	"log"
	"os"
	"path/filepath"
)

var maxPool int

var index = mgo.Index{
	Key:        []string{"username"},
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
	err = c.Collection.EnsureIndex(index)
	defer c.Close()
	if err != nil {
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
