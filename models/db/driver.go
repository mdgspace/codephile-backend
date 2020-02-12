package db

import (
	"github.com/astaxie/beego"
	"github.com/globalsign/mgo"
	"log"
	"os"
	"path/filepath"
	"runtime"
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
		_, file, _, _ := runtime.Caller(0)
		if err := beego.LoadAppConfig("ini", filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(file))), "conf/app.conf")); err != nil {
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
