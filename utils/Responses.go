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

type ChangePasswordInput struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

type TweetResponse struct {
	ID        uint   `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	File      string `json:"file"`
}
