package models

import (
	"context"
	"encoding/json"
	"github.com/globalsign/mgo/bson"
	. "github.com/mdg-iitr/Codephile/conf"
	. "github.com/mdg-iitr/Codephile/errors"
	"github.com/mdg-iitr/Codephile/models/db"
	"github.com/mdg-iitr/Codephile/models/types"
	search "github.com/mdg-iitr/Codephile/services/elastic"
	"github.com/olivere/elastic/v7"
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
	if err != nil {
		return "", err
	}
	u.Password = string(hash)
	err = collection.Collection.Insert(u)
	if err != nil {
		return "", UserAlreadyExistError
	}
	client := search.GetElasticClient()
	_, err = client.Index().Index("codephile").BodyJson(u).Id(u.ID.String()).Refresh("true").Do(context.Background())
	if err != nil {
		return "", err
	}

	go func() {
		for _, value := range ValidSites {
			_ = AddSubmissions(u.ID, value)
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
		return nil, err
	}
	return &user, nil
}

func GetAllUsers() ([]types.User, error) {
	var users []types.User
	collection := db.NewUserCollectionSession()
	defer collection.Close()
	err := collection.Collection.Find(nil).Select(bson.M{"_id": 1, "username": 1,
		"handle": 1, "lastfetched": 1,
		"picture": 1, "fullname": 1, "institute": 1}).All(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
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
		return nil, UserAlreadyExistError
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

func AuthenticateUser(username string, password string) (*types.User, bool) {
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

func UserExists(username string) (bool, error) {
	collection := db.NewUserCollectionSession()
	defer collection.Close()
	c, err := collection.Collection.Find(bson.M{"username": username}).Count()
	if err != nil {
		log.Println(err.Error())
		return false, err
	}
	if c > 0 {
		return true, nil
	}
	return false, nil
}

func SearchUser(query string, c int) ([]interface{}, error) {
	pq := elastic.NewQueryStringQuery("*" + query + "*").
		Field("username").Field("fullname").
		Field("handle.codechef").Field("handle.spoj").
		Field("handle.codeforces").Field("handle.hackerrank").
		Fuzziness("4")
	q := elastic.NewMultiMatchQuery(query,
		"username", "fullname",
		"handle.codechef", "handle.spoj",
		"handle.codeforces", "handle.hackerrank",
	).Fuzziness("4")
	bq := elastic.NewBoolQuery().Should(q, pq)
	client := search.GetElasticClient()
	result, err := client.Search().Index("codephile").
		Pretty(false).Query(bq).Size(c).
		Do(context.Background())
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	var results []interface{}
	for _, hit := range result.Hits.Hits {
		var result interface{}
		err := json.Unmarshal(hit.Source, &result)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}
