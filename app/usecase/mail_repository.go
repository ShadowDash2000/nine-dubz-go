package usecase

type MailRepository interface {
	SendMail(from, to, subject, content string) error
}
