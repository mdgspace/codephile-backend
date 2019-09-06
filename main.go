package main

import (
	_ "github.com/mdg-iitr/Codephile/routers"
	"log"
	"os"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/joho/godotenv"

)

func main() {
	if err := godotenv.Load("conf/.env"); err != nil {
		log.Print("No .env file found")
	}

    port, exists := os.LookupEnv("PORT")

	if exists {
		i1, err := strconv.Atoi(port)
		if err == nil {
			beego.BConfig.Listen.HTTPPort = i1
		}
	} else {
		log.Print("No Port variable found")
	}

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
