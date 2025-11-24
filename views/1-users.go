package views

type User struct {
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	// Role     string `json:"role"`
}

type SignInCodeRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type TokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	// Need2register bool   `json:"need2register"`
}
