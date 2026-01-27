package event

import "time"

const ForgotPasswordSubject = "auth.forgot_password"

type ForgotPasswordEvent struct {
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
}
