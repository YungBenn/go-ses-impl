package main

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"gopkg.in/gomail.v2"
)

type Recipient struct {
	toEmails  []string
	ccEmails  []string
	bccEmails []string
}

func main() {
	senderName := "Admin"
	subject := "Sample Email using AWS SES and Go"
	fromEmail := "test@gmail.com"
	messageTitle := "Hello, World!"
	messageBody := "This is a sample email sent using AWS SES and Go."

	message := fmt.Sprintf(`
	<html>
		<body style="background-color: red;">
			<h1 style="color: #333;">%s</h1>
			<p style="color: #666;">%s</p>
		</body>
	</html>
	`, messageTitle, messageBody)

	recipient := Recipient{
		toEmails:  []string{"adisuryo22@gmail.com"},
		ccEmails:  []string{"adisuryo22@gmail.com"},
		bccEmails: []string{"adisuryo22@gmail.com"},
	}

	attachments := []string{}

	SendEmailRawSES(senderName, message, subject, fromEmail, recipient, attachments)
}

func AwsConfig() (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedCredentialsFiles(
			[]string{"aws/credentials"},
		),
		config.WithSharedConfigFiles(
			[]string{"aws/config"},
		),
	)
	if err != nil {
		return aws.Config{}, err
	}

	return cfg, nil
}

// SendEmailSES sends email to specified email IDs
func SendEmailRawSES(senderName string, messageBody string, subject string, fromEmail string, recipient Recipient, attachments []string) {
	cfg, err := AwsConfig()
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// create raw message
	msg := gomail.NewMessage()

	// set to section
	var recipients []string
	for _, r := range recipient.toEmails {
		recipient := r
		recipients = append(recipients, recipient)
	}

	// Set to emails
	msg.SetHeader("To", recipient.toEmails...)

	// cc mails mentioned
	if len(recipient.ccEmails) != 0 {
		// Need to add cc mail IDs also in recipient list
		for _, r := range recipient.ccEmails {
			recipient := r
			recipients = append(recipients, recipient)
		}
		msg.SetHeader("cc", recipient.ccEmails...)
	}

	// bcc mails mentioned
	if len(recipient.bccEmails) != 0 {
		// Need to add bcc mail IDs also in recipient list
		for _, r := range recipient.bccEmails {
			recipient := r
			recipients = append(recipients, recipient)
		}
		msg.SetHeader("bcc", recipient.bccEmails...)
	}

	// create an SES session.
	svc := ses.NewFromConfig(cfg)

	msg.SetAddressHeader("From", fromEmail, senderName)
	msg.SetHeader("To", recipient.toEmails...)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", messageBody)

	// If attachments exists
	if len(attachments) != 0 {
		for _, f := range attachments {
			msg.Attach(f)
		}
	}

	// create a new buffer to add raw data
	var emailRaw bytes.Buffer
	msg.WriteTo(&emailRaw)

	// create new raw message
	message := types.RawMessage{
		Data: emailRaw.Bytes(),
	}

	input := &ses.SendRawEmailInput{
		Source:       &fromEmail,
		Destinations: recipients,
		RawMessage:   &message,
	}

	// send raw email
	_, err = svc.SendRawEmail(context.TODO(), input)
	if err != nil {
		log.Println("Error sending mail - ", err)
		return
	}

	log.Println("Email sent successfully to: ", recipient.toEmails)
}
