package utils

type AuthInput struct {
	UserName string `binding:"required"`
	Password string `binding:"required"`
	Email    string `binding:"required"`
	Bio      string `json:"bio"`
}

type LoginInput struct {
	UserName string
	Email    string
	Password string `binding:"required"`
}

type UserResponse struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
}
