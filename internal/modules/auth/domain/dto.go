package domain

type RegisterInput struct {
	Name                 string
	Email                string
	Password             string
	PasswordConfirmation string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthOutput struct {
	AccessToken  string
	RefreshToken string
}

type ResetPasswordInput struct {
	Token       string
	NewPassword string
}
