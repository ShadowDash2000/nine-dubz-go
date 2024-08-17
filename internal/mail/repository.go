package mail

import (
	"fmt"
	"net/smtp"
)

type Repository struct {
	Host     string
	Port     string
	Username string
	Password string
}

func (mr *Repository) SendMail(from, to, subject, content string) error {
	auth := smtp.PlainAuth("", mr.Username, mr.Password, mr.Host)

	msg := fmt.Sprintf(
		"To: %s\r\n"+
			"From: %s\r\n"+
			"Subject: %s\r\n"+
			"\r\n"+
			"%s\r\n",
		to, from, subject, content)

	err := smtp.SendMail(mr.Host+":"+mr.Port, auth, from, []string{to}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}
