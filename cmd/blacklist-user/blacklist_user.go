package blacklist_user

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	_ "github.com/mdg-iitr/Codephile/conf"
	"github.com/mdg-iitr/Codephile/services/auth"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run ./blacklist_user <uid>")
		os.Exit(1)
	}
	err := auth.BlacklistUser(bson.ObjectId(os.Args[1]))
	if err != nil {
		fmt.Println(err.Error())
	}

}
