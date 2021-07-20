package infra

type MailDriver interface {
	SendMail(e Email) error
}

type Email struct {
	From    string
	To      []string
	Message []byte
}
