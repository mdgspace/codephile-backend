package test

import (
	"fmt"

	"github.com/globalsign/mgo/bson"
	"github.com/mdg-iitr/Codephile/conf"
	"github.com/mdg-iitr/Codephile/models/db"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/astaxie/beego"
	_ "github.com/mdg-iitr/Codephile/conf"
	"github.com/mdg-iitr/Codephile/models"
	_ "github.com/mdg-iitr/Codephile/models"
	"github.com/mdg-iitr/Codephile/models/types"
	_ "github.com/mdg-iitr/Codephile/routers"
	"github.com/mdg-iitr/Codephile/services/auth"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	beego.TestBeegoInit(conf.AppRootDir)
	db.NewUserCollectionSession().DropDatabase()
}

// TestGet is a sample to run an endpoint test
func TestGetAllUsers(t *testing.T) {
	uid, _ := models.AddUser(types.User{
		ID:        bson.NewObjectId(),
		Email:     "test@abc.com",
		Username:  "test",
		FullName:  "Test User",
		Institute: "IIT Roorkee",
		Password:  "password",
	})
	token := auth.GenerateToken(uid)
	r, _ := http.NewRequest("GET", "/v1/user/all", nil)
	r.Header.Set("Authorization", token)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	beego.Trace("testing", "TestGet", "Code[%d]\n%s", w.Code, w.Body.String())
	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
	})
}

// A sample query to check if database is connected
func TestFind(t *testing.T) {
	var user types.User
	collection := db.NewUserCollectionSession()
	defer collection.Close()
	err := collection.Collection.Find(bson.M{"username": "nano_nish"}).Select(bson.M{"_id": 1, "username": 1, "email": 1,
		"handle": 1, "lastfetched": 1, "profiles": 1,
		"picture": 1, "fullname": 1, "institute": 1, "submissions": bson.M{"$slice": 5}}).One(&user)
	// fmt.Println(err.Error())
	fmt.Print(user)
	if err != nil {
		// return nil, err
		t.Error(err)
	}
}
