package mail

type Interactor interface {
	SendMail(from, to, subject, content string) error
}
