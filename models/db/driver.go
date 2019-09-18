package db

import (
	"github.com/astaxie/beego"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var maxPool int

func init() {
	var err error
	maxPool, err = beego.AppConfig.Int("DBMaxPool")
	if err != nil {
		panic(err)
	}
	// init method to start db
	err = godotenv.Load("conf/.env")
	if err != nil {
		log.Println("No .env file found")
	}
	checkAndInitServiceConnection()
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
