package models

import (
	"context"
	"errors"
	"github.com/globalsign/mgo/bson"
	. "github.com/mdg-iitr/Codephile/conf"
	"github.com/mdg-iitr/Codephile/models/db"
	"github.com/mdg-iitr/Codephile/models/types"
	search "github.com/mdg-iitr/Codephile/services/elastic"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func AddUser(u types.User) (string, error) {
	u.ID = bson.NewObjectId()
	collection := db.NewUserCollectionSession()
	defer collection.Close()
	//hashing the password
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	//data type of hash is []byte
	u.Password = string(hash)
	if err != nil {
		log.Println(err)
	}
	err = collection.Collection.Insert(u)
	if err != nil {
		log.Println(err)
		return "", errors.New("Could not create user: Username already exists")
	}
	client := search.GetElasticClient()
	_, err = client.Index().Index("codephile").BodyJson(u).Id(u.ID.String()).Refresh("true").Do(context.Background())
	if err != nil {
		log.Println(err.Error())
	}

	var valid_sites = []string{HACKERRANK, CODECHEF, CODEFORCES, SPOJ}

	go func() {
		for _, value := range valid_sites {
			_ = AddSubmissions(&u, value)
		}
	}()

	return u.ID.Hex(), nil
}

func GetUser(uid bson.ObjectId) (*types.User, error) {
	var user types.User
	collection := db.NewUserCollectionSession()
	defer collection.Close()
	err := collection.Collection.FindId(uid).Select(bson.M{"_id": 1, "username": 1,
		"handle": 1, "lastfetched": 1,
		"picture": 1, "fullname": 1, "institute": 1}).One(&user)
	//fmt.Println(err.Error())
	if err != nil {
		return nil, errors.New("user not exists")
	}
	return &user, nil
}

func GetAllUsers() []types.User {
	var users []types.User
	collection := db.NewUserCollectionSession()
	defer collection.Close()
	err := collection.Collection.Find(nil).Select(bson.M{"_id": 1, "username": 1,
		"handle": 1, "lastfetched": 1,
		"picture": 1, "fullname": 1, "institute": 1}).All(&users)
	if err != nil {
		panic(err)
	}
	return users
}

func UpdateUser(uid bson.ObjectId, uu *types.User) (a *types.User, err error) {
	var updateDoc = bson.M{}
	var elasticDoc = map[string]interface{}{}
	var newHandle types.Handle
	if uu.Username != "" {
		updateDoc["username"] = uu.Username
		elasticDoc["username"] = uu.Username
	}
	if uu.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(uu.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		uu.Password = string(hash)
		updateDoc["password"] = uu.Password
		elasticDoc["password"] = uu.Password
	}
	if uu.Institute != "" {
		updateDoc["institute"] = uu.Institute
		elasticDoc["institute"] = uu.Institute
	}
	if uu.FullName != "" {
		updateDoc["fullname"] = uu.FullName
		elasticDoc["fullname"] = uu.FullName
	}
	if uu.Handle.Codechef != "" {
		updateDoc["handle.codechef"] = uu.Handle.Codechef
		newHandle.Codechef = uu.Handle.Codechef
	}
	if uu.Handle.Codeforces != "" {
		updateDoc["handle.codeforces"] = uu.Handle.Codeforces
		newHandle.Codeforces = uu.Handle.Codeforces
	}
	if uu.Handle.Hackerearth != "" {
		updateDoc["handle.hackerearth"] = uu.Handle.Hackerearth
		newHandle.Hackerearth = uu.Handle.Hackerearth
	}
	if uu.Handle.Hackerrank != "" {
		updateDoc["handle.hackerrank"] = uu.Handle.Hackerrank
		newHandle.Hackerrank = uu.Handle.Hackerrank
	}
	if uu.Handle.Spoj != "" {
		updateDoc["handle.spoj"] = uu.Handle.Spoj
		newHandle.Spoj = uu.Handle.Spoj
	}
	elasticDoc["handle"] = newHandle

	collection := db.NewUserCollectionSession()
	defer collection.Close()
	err = collection.Collection.UpdateId(uid, bson.M{"$set": updateDoc})
	if err != nil {
		log.Println(err.Error())
		err = errors.New("username already exists")
		return nil, err
	}
	client := search.GetElasticClient()
	_, err = client.Update().Index("codephile").Id(uid.String()).Doc(elasticDoc).Do(context.Background())
	if err != nil {
		log.Println(err.Error())
	}
	u, err := GetUser(uid)
	if err != nil {
		return nil, err
	}
	return u, err
}

func AutheticateUser(username string, password string) (*types.User, bool) {
	var user types.User
	collection := db.NewUserCollectionSession()
	defer collection.Close()
	err := collection.Collection.Find(bson.M{"username": username}).One(&user)
	//fmt.Println(err.Error())
	if err != nil {
		log.Println(err)
		return nil, false
	}

	err2 := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err2 != nil {
		log.Println(err2)
		return nil, false
	} else {
		return &user, true
	}

}

func UpdatePicture(uid bson.ObjectId, url string) error {
	client := search.GetElasticClient()
	_, err := client.Update().Index("codephile").Id(uid.String()).Doc(map[string]interface{}{"picture": url}).Do(context.Background())
	if err != nil {
		log.Println(err.Error())
	}
	coll := db.NewUserCollectionSession()
	defer coll.Close()
	_, err = coll.Collection.UpsertId(uid, bson.M{"$set": bson.M{"picture": url}})
	if err != nil {
		return err
	}
	return nil
}

func GetPicture(uid bson.ObjectId) string {
	var user types.User
	coll := db.NewUserCollectionSession()
	defer coll.Close()
	err := coll.Collection.FindId(uid).Select(bson.M{"picture": 1}).One(&user)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	return user.Picture
}
