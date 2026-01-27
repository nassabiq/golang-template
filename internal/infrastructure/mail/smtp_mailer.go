package mail

import (
	"fmt"
	"net/smtp"
)

type SMTPMailer struct {
	host string
	port int
	from string
}

func NewSMTPMailer(host string, port int, from string) *SMTPMailer {
	return &SMTPMailer{
		host: host,
		port: port,
		from: from,
	}
}

func (mailer *SMTPMailer) Send(to, subject, body string) error {
	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		mailer.from,
		to,
		subject,
		body,
	)

	addr := fmt.Sprintf("%s:%d", mailer.host, mailer.port)

	return smtp.SendMail(addr, nil, mailer.from, []string{to}, []byte(msg))
}
