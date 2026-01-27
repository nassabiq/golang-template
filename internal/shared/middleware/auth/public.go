package auth

func isPublicMethod(method string) bool {
	switch method {
	case
		"/auth.v1.AuthService/Login",
		"/auth.v1.AuthService/Register",
		"/auth.v1.AuthService/Refresh",
		"/auth.v1.AuthService/ForgotPassword",
		"/auth.v1.AuthService/ResetPassword":
		return true
	default:
		return false
	}
}
