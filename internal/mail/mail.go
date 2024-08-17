package mail

import (
	"log"
	"os"
)

type UseCase struct {
	MailInteractor   Interactor
	DefaultEmailFrom string
}

func New() *UseCase {
	host, ok := os.LookupEnv("MAIL_HOST")
	if !ok {
		log.Println("No MAIL_HOST environment variable")
	}
	port, ok := os.LookupEnv("MAIL_PORT")
	if !ok {
		log.Println("No MAIL_PORT environment variable")
	}
	emailFrom, ok := os.LookupEnv("MAIL_EMAIL")
	if !ok {
		log.Println("No MAIL_EMAIL environment variable")
	}
	login, ok := os.LookupEnv("MAIL_LOGIN")
	if !ok {
		log.Println("No MAIL_LOGIN environment variable")
	}
	password, ok := os.LookupEnv("MAIL_PASSWORD")
	if !ok {
		log.Println("No MAIL_PASSWORD environment variable")
	}

	return &UseCase{
		MailInteractor: &Repository{
			Host:     host,
			Port:     port,
			Username: login,
			Password: password,
		},
		DefaultEmailFrom: emailFrom,
	}
}

func (uc *UseCase) SendMail(to, subject, content string) error {
	return uc.MailInteractor.SendMail("Nine Dubz", uc.DefaultEmailFrom, to, subject, content)
}
