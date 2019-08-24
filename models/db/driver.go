package db

import ("github.com/astaxie/beego")

var maxPool int

func init() {
	var err error
	maxPool, err = beego.AppConfig.Int("DBMaxPool")
	if err != nil {
		panic(err)
	}
	// init method to start db
	checkAndInitServiceConnection()
}

func checkAndInitServiceConnection() {
	if service.baseSession == nil {
		service.URL = beego.AppConfig.String("DBPath")
		err := service.New()
		if err != nil {
			panic(err)
		}
	}
}