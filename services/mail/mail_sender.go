package mail

import (
	"context"
	"encoding/base64"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

//EMAIL_SECRET=C65v5DEIHlDdEU1shHYrlnns
//EMAIL_CLIENT=918242223374-k6ihlaib00kvt65j82gnqnrukms6ocrg.apps.googleusercontent.com

//send a text mail to the receiver's email id
func SendMail(to string, subject string, body string) {
	var GmailService *gmail.Service

	config := oauth2.Config{
		ClientID:     os.Getenv("EMAIL_CLIENT"),
		ClientSecret: os.Getenv("EMAIL_SECRET"),
		Endpoint:     google.Endpoint,
		RedirectURL:  os.Getenv("REDIRECT_URL"),
	}

	token := oauth2.Token{
		AccessToken:  os.Getenv("EMAIL_ACCESS_TOKEN"),
		RefreshToken: os.Getenv("EMAIL_REFRESH_TOKEN"),
		TokenType:    "Bearer",
		Expiry:       time.Now(),
	}

	//refreshing token
	var tokenSource = config.TokenSource(context.TODO(), &token)

	updatedToken, Error := tokenSource.Token()
	if Error != nil {
		log.Fatalf(Error.Error())
	} else if (*updatedToken).AccessToken != token.AccessToken {
		os.Setenv("EMAIL_ACCESS_TOKEN", (*updatedToken).AccessToken)
		token = *updatedToken
	}

	//sending mail
	tokenSource = config.TokenSource(context.Background(), &token)

	srv, err := gmail.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		log.Printf("Unable to retrieve Gmail client: %v", err)
	}

	GmailService = srv
	if GmailService != nil {
		log.Println("Email service is initialized")
	}

	var message gmail.Message

	emailTo := "To: " + to + "\r\n"
	sub := "Subject: " + subject + "\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	msg := []byte(emailTo + sub + mime + "\n" + body)

	message.Raw = base64.URLEncoding.EncodeToString(msg)

	if GmailService != nil {
		_, Err := GmailService.Users.Messages.Send("me", &message).Do()
		if Err != nil {
			log.Fatalf(err.Error())
		}
	}

}
