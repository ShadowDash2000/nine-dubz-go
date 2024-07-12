package controller

import (
	"fmt"
	"net/smtp"
)

type MailRepository struct {
	Host     string
	Username string
	Password string
}

func (mr *MailRepository) SendMail(from, to, subject, content string) error {
	auth := smtp.PlainAuth("", mr.Username, mr.Password, mr.Host)

	msg := fmt.Sprintf(`To: %s
	From: %s
	Subject: %s
	
	%s
	`, to, from, subject, content)

	err := smtp.SendMail(mr.Host+":587", auth, from, []string{to}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}
