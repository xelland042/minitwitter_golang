package utils

type AuthInput struct {
	UserName string `binding:"required"`
	Password string `binding:"required"`
	Email    string `binding:"required"`
	Bio      string `json:"bio"`
}

type LoginInput struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `binding:"required"`
}

type UserResponse struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
}

type TweetCreate struct {
	Title string `json:"title" binding:"required"`
	Body  string `json:"body" binding:"required"`
}
