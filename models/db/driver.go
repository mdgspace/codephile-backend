package db

import (
	"github.com/astaxie/beego"
	"github.com/globalsign/mgo"
	"log"
	"os"
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
		panic(err)
	}
	// init method to start db
	checkAndInitServiceConnection()
	err = NewCollectionSession("coduser").Session.EnsureIndex(index)
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
