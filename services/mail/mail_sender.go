package mail

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"log"
	"os"
)

//send a text mail to the receiver's email id
func SendMail(to string, subject string, body string) {
	b := []byte(os.Getenv("EMAIL_CONFIG"))
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Println(err.Error())
		return
	}
	client := config.Client(context.TODO(), &oauth2.Token{
		AccessToken:  os.Getenv("EMAIL_ACCESS_TOKEN"),
		RefreshToken: os.Getenv("EMAIL_REFRESH_TOKEN"),
	})
	srv, _ := gmail.New(client)
	r, err := srv.Users.Labels.List("me").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}
	if len(r.Labels) == 0 {
		fmt.Println("No labels found.")
		return
	}
	fmt.Println("Labels:")
	for _, l := range r.Labels {
		fmt.Printf("- %s\n", l.Name)
	}
}

