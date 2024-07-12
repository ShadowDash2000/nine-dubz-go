package usecase

type MailInteractor struct {
	MailRepository MailRepository
}

func (mi *MailInteractor) SendMail(from, to, subject, content string) error {
	return mi.MailRepository.SendMail(from, to, subject, content)
}
