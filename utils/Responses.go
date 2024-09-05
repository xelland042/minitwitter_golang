package utils

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserInfoResponse struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
	Picture  string `json:"picture"`
}
