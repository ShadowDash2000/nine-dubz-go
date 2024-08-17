package mail

type Interactor interface {
	SendMail(name, from, to, subject, content string) error
}
