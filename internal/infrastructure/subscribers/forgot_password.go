package subscribers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nassabiq/golang-template/internal/infrastructure/mail"
	"github.com/nats-io/nats.go"
)

type ForgotPasswordSubscriber struct {
	mailer mail.Mailer
}

func NewForgotPasswordSubscriber(mailer mail.Mailer) *ForgotPasswordSubscriber {
	return &ForgotPasswordSubscriber{mailer: mailer}
}

func (subscriber *ForgotPasswordSubscriber) Subject() string {
	return "auth.forgot_password"
}

func (subscriber *ForgotPasswordSubscriber) Durable() string {
	return "email-forgot-password"
}

func (sub *ForgotPasswordSubscriber) Subscribe(js nats.JetStreamContext) error {
	_, err := js.Subscribe(sub.Subject(),
		func(msg *nats.Msg) {
			var event struct {
				Email string `json:"email"`
				Token string `json:"token"`
			}

			if err := json.Unmarshal(msg.Data, &event); err != nil {
				log.Println(err)
				return
			}

			link := "http://localhost:3000/reset-password?token=" + event.Token

			_ = sub.mailer.Send(event.Email, "Reset Password", fmt.Sprintf("Klik link berikut:\n\n%s", link))
			msg.Ack()
		},
		nats.Durable(sub.Durable()),
		nats.ManualAck(),
	)
	return err
}
