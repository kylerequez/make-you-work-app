package utils

import (
	"bytes"
	"strconv"
	"text/template"

	"github.com/go-gomail/gomail"
)

const EMAIL_LOCATION string = "src/api/privates/emails/"

type EmailDetails struct {
	From       string
	To         []string
	Subject    string
	Body       string
	Attachment *string
}

type SMTPCredentials struct {
	host     string
	port     int
	username string
	password string
}

func SendEmail(subject string, from string, to []string, attachment *string, filename string, data any) error {
	m := gomail.NewMessage()

	m.SetHeader("Subject", subject)
	m.SetHeader("To", to...)
	m.SetHeader("From", from)

	body, err := ParseHTML(filename, data)
	if err != nil {
		return err
	}

	m.SetBody("text/html", body)
	if attachment != nil {
		m.Attach(*attachment)
	}

	smtp, err := GetSMTPCredentials()
	if err != nil {
		return err
	}

	d := gomail.NewDialer(smtp.host, smtp.port, smtp.username, smtp.password)
	if err = d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func GetSMTPCredentials() (*SMTPCredentials, error) {
	host, err := GetEnv("SMTP_HOST")
	if err != nil {
		return nil, err
	}
	strPort, err := GetEnv("SMTP_PORT")
	if err != nil {
		return nil, err
	}
	username, err := GetEnv("SMTP_USER")
	if err != nil {
		return nil, err
	}
	password, err := GetEnv("SMTP_PASS")
	if err != nil {
		return nil, err
	}

	intPort, err := strconv.Atoi(*strPort)
	if err != nil {
		return nil, err
	}

	return &SMTPCredentials{
		host:     *host,
		port:     intPort,
		username: *username,
		password: *password,
	}, nil
}

func ParseHTML(filename string, data any) (string, error) {
	tmpl, err := template.ParseFiles(EMAIL_LOCATION + filename)
	if err != nil {
		return "", err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return "", err
	}

	return body.String(), nil
}

