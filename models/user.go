package models

import (
	"errors"
	"github.com/globalsign/mgo/bson"
	"github.com/mdg-iitr/Codephile/models/db"
)

type User struct {
	ID       bson.ObjectId `bson:"_id" json:"id" schema:"-"`
	Username string        `bson:"username" json:"username" schema:"username"`
	Password string        `bson:"password" json:"-" schema:"password"`
	Handle   Handle        `bson:"handle" json:"handle" schema:"handle"`
}
type Handle struct {
	Codeforces  string `bson:"codeforces" json:"codeforces" schema:"codeforces"`
	Codechef    string `bson:"codechef" json:"codechef" schema:"codechef"`
	Spoj        string `bson:"spoj" json:"spoj" schema:"spoj"`
	Hackerrank  string `bson:"hackerrank" json:"hackerrank" schema:"hackerrank"`
	Hackerearth string `bson:"hackerearth" json:"hackerearth" schema:"hackerearth"`
}

func AddUser(u User) (string, error) {
	u.ID = bson.NewObjectId()
	collection := db.NewCollectionSession("coduser")
	defer collection.Close()
	err := collection.Session.Insert(u)
	if err != nil {
		panic(err)
	}
	return u.ID.Hex(), nil
}

func GetUser(uid bson.ObjectId) (*User, error) {
	var user User
	collection := db.NewCollectionSession("coduser")
	defer collection.Close()
	err := collection.Session.FindId(uid).One(&user)
	//fmt.Println(err.Error())
	if err != nil {
		return nil, errors.New("user not exists")
	}
	return &user, nil
}

func GetAllUsers() []User {
	var users []User
	collection := db.NewCollectionSession("coduser")
	defer collection.Close()
	err := collection.Session.Find(nil).All(&users)
	if err != nil {
		panic(err)
	}
	return users
}

func UpdateUser(uid bson.ObjectId, uu *User) (a *User, err error) {
	if u, err := GetUser(uid); err == nil {
		if uu.Username != "" {
			u.Username = uu.Username
		}
		if uu.Handle.Codechef != "" {
			u.Handle.Codechef = uu.Handle.Codechef
		}
		if uu.Handle.Codeforces != "" {
			u.Handle.Codeforces = uu.Handle.Codeforces
		}
		if uu.Handle.Hackerearth != "" {
			u.Handle.Hackerearth = uu.Handle.Hackerearth
		}
		if uu.Handle.Hackerrank != "" {
			u.Handle.Hackerrank = uu.Handle.Hackerrank
		}
		if uu.Handle.Hackerearth != "" {
			u.Handle.Hackerearth = uu.Handle.Hackerearth
		}
		collection := db.NewCollectionSession("coduser")
		_, err := collection.Session.UpsertId(uid, &u)
		return u, err
	}
	return nil, errors.New("User Not Exist")
}

// func Login(username, password string) bool {
// 	for _, u := range UserList {
// 		if u.Username == username && u.Password == password {
// 			return true
// 		}
// 	}
// 	return false
// }

// func DeleteUser(uid string) {
// 	delete(UserList, uid)
// }
