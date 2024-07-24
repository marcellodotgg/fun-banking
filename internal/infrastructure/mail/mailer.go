package mail

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"text/template"

	"gopkg.in/mail.v2"
)

type Mailer struct{}

func (m Mailer) Send(to, subject, templateName string, data interface{}) error {
	username := os.Getenv("EMAIL_USERNAME")
	password := os.Getenv("EMAIL_PASSWORD")
	host := os.Getenv("EMAIL_HOST")
	port, _ := strconv.Atoi(os.Getenv("EMAIL_PORT"))

	email := mail.NewMessage()
	email.SetHeader("From", username)
	email.SetHeader("To", to)
	email.SetHeader("Subject", subject)

	if err := m.textBody(email, templateName, data); err != nil {
		return err
	}

	if err := m.htmlBody(email, templateName, data); err != nil {
		return err
	}

	dailer := mail.NewDialer(host, port, username, password)
	if err := dailer.DialAndSend(email); err != nil {
		fmt.Println("[ERROR]: Unable to send email to:", to)
		return err
	}

	fmt.Println("[SUCCESS] Sent an email to:", to)
	return nil
}

func (m Mailer) textBody(email *mail.Message, templateName string, data interface{}) error {
	return m.setEmailBody(email, data, templateName+".txt", "text/plain")
}

func (m Mailer) htmlBody(email *mail.Message, templateName string, data interface{}) error {
	return m.setEmailBody(email, data, templateName+".html", "text/html")
}

func (m Mailer) setEmailBody(email *mail.Message, data interface{}, templateName, contentType string) error {
	templ, err := template.ParseFiles(fmt.Sprintf("templates/emails/%s", templateName))

	if err != nil {
		return err
	}

	buffer := &bytes.Buffer{}
	if err := templ.Execute(buffer, &data); err != nil {
		return err
	}

	email.SetBody(contentType, buffer.String())

	return nil
}
