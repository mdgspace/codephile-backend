package controllers

import (
	// "encoding/json"
	"github.com/astaxie/beego"
	//"github.com/globalsign/mgo/bson"
	"github.com/mdg-iitr/Codephile/models"
	"net/http"
	"io/ioutil"
	"log"
)

//Controller to display contests
type ContestController struct {
	beego.Controller
}

// @Title GetContests
// @Description displays all contests
// @Success 200 {object} models.S
// @Failure 403 error
// @router / [get]
func (u *ContestController) GetContests() {
	fetchDataAndParse();
	contests := models.ReturnContests()
	u.Data["json"] = contests
	u.ServeJSON()
}

//function to fetch data from URL and parse
func fetchDataAndParse()  {
	resp, err := http.Get("https://contesttrackerapi.herokuapp.com/")

	if err != nil {
		log.Println("Error")
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return
	}	
	models.ParseContests(body);
}



