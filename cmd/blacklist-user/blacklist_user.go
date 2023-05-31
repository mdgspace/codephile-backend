package main

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	_ "github.com/mdg-iitr/Codephile/conf"
	"github.com/mdg-iitr/Codephile/services/auth"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run ./blacklist_user <uid>")
		os.Exit(1)
	}
	// TODO: upgrade driver
	id, _ := primitive.ObjectIDFromHex(os.Args[1])
	err := auth.BlacklistUser(id)
	if err != nil {
		fmt.Println(err.Error())
	}

}
