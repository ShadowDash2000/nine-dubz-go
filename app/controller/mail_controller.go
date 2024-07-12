package controller

import (
	"fmt"
	"nine-dubz/app/usecase"
	"os"
)

type MailController struct {
	MailInteractor   usecase.MailInteractor
	DefaultEmailFrom string
}

func NewMailController() *MailController {
	host, ok := os.LookupEnv("MAIL_HOST")
	if !ok {
		fmt.Println("No MAIL_HOST environment variable")
	}
	emailFrom, ok := os.LookupEnv("MAIL_EMAIL")
	if !ok {
		fmt.Println("No MAIL_EMAIL environment variable")
	}
	login, ok := os.LookupEnv("MAIL_LOGIN")
	if !ok {
		fmt.Println("No MAIL_LOGIN environment variable")
	}
	password, ok := os.LookupEnv("MAIL_PASSWORD")
	if !ok {
		fmt.Println("No MAIL_PASSWORD environment variable")
	}

	return &MailController{
		MailInteractor: usecase.MailInteractor{
			MailRepository: &MailRepository{
				Host:     host,
				Username: login,
				Password: password,
			},
		},
		DefaultEmailFrom: emailFrom,
	}
}

func (mc *MailController) SendMail(from, to, subject, content string) error {
	return mc.MailInteractor.SendMail(from, to, subject, content)
}
