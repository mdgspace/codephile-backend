package types

import (
	"encoding/json"
	"time"

	"github.com/globalsign/mgo/bson"
	. "github.com/mdg-iitr/Codephile/errors"
)

type User struct {
	ID                  bson.ObjectId         `bson:"_id" json:"id" schema:"-"`
	Username            string                `bson:"username" json:"username" schema:"username"`
	Email               string                `bson:"email" json:"email" schema:"email"`
	FullName            string                `bson:"fullname" json:"fullname" schema:"fullname"`
	Institute           string                `bson:"institute" json:"institute" schema:"institute"`
	Password            string                `bson:"password" json:"-" schema:"password"`
	Picture             string                `bson:"picture" json:"picture"`
	Verified            bool                  `bson:"verified" schema:"-" json:"-"`
	Handle              Handle                `bson:"handle" json:"handle" schema:"handle"`
	Submissions         []Submission          `bson:"submissions" json:"recent_submissions" schema:"-"`
	Profiles            AllProfiles           `json:"profiles" bson:"profiles" schema:"-"`
	Last                LastFetchedSubmission `bson:"lastfetched" json:"-"`
	FollowingUsers      []Following           `bson:"followingUsers" json:"-"`
	NoOfFollowing       int                   `bson:"-" json:"no_of_following"`
	SolvedProblemsCount SolvedProblemsCount   `json:"solved_problems_count"`
}
type LastFetchedSubmission struct {
	Codechef   time.Time `bson:"codechef"`
	Codeforces time.Time `bson:"codeforces"`
	Hackerrank time.Time `bson:"hackerrank"`
	Spoj       time.Time `bson:"spoj"`
	Leetcode   time.Time `bson:"leetcode"`
}
type Handle struct {
	Codeforces  string `bson:"codeforces" json:"codeforces" schema:"codeforces"`
	Codechef    string `bson:"codechef" json:"codechef" schema:"codechef"`
	Spoj        string `bson:"spoj" json:"spoj" schema:"spoj"`
	Hackerrank  string `bson:"hackerrank" json:"hackerrank" schema:"hackerrank"`
	Hackerearth string `bson:"hackerearth" json:"hackerearth" schema:"hackerearth"`
	Leetcode    string `bson:"leetcode" json:"leetcode" schema:"leetcode"`
}

func (u *User) UnmarshalJSON(b []byte) error {
	var m map[string]interface{}
	err := json.Unmarshal(b, &m)
	if val, ok := m["password"]; ok {
		u.Password = val.(string)
	} else {
		return FieldEmptyError
	}
	if val, ok := m["username"]; ok {
		u.Username = val.(string)
	} else {
		return FieldEmptyError
	}
	if val, ok := m["fullname"]; ok {
		u.FullName = val.(string)
	}
	if val, ok := m["institute"]; ok {
		u.Institute = val.(string)
	}
	if val, ok := m["email"]; ok {
		u.Email = val.(string)
	} else {
		return FieldEmptyError
	}
	if val, ok := m["handle"]; ok {
		d, _ := json.Marshal(val)
		err = json.Unmarshal(d, &u.Handle)
	}
	return err
}

type SearchDoc struct {
	ID        bson.ObjectId `json:"id"`
	Username  string        `json:"username"`
	FullName  string        `json:"fullname"`
	Institute string        `json:"institute"`
	Picture   string        `json:"picture"`
	Handle    Handle        `json:"handle"`
}

type UpdatePassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
