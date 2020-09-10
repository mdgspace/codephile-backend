package mail

import (
	"fmt"
	"net/smtp"
	"os"
)

//send a text mail to the receiver's email id
func SendMail(to string, subject string, body string) {
	// Sender data.
	from := os.Getenv("EMAIL_SERVER_USER")
	password := os.Getenv("EMAIL_SERVER_PASS")
	smtpHost := os.Getenv("EMAIL_SMTP_HOST")
	smtpPort := os.Getenv("EMAIL_SMTP_PORT")

	// Message.
	msg := "From: Codephile" + "\n" + "To: " + to+ "\n" + "Subject: "+ subject + "\n\n" + body
	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost + ":" + smtpPort, auth, from, []string{to}, []byte(msg))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")
}