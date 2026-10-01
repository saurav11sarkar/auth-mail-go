package auth

// Service inputs are independent of HTTP JSON and validation tags.
type RegisterInput struct {
	Name     string
	Email    string
	Password string
}
type LoginInput struct {
	Email    string
	Password string
}
type RefreshInput struct{ RefreshToken string }
type ForgetPasswordInput struct{ Email string }
type ResetPasswordInput struct {
	Email    string
	OTP      string
	Password string
}

// LoginResult is mapped to a public response by the handler.
type LoginResult struct {
	Account      Account
	AccessToken  string
	RefreshToken string
}
