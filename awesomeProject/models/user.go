package models

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username" binding:"required" example:"john_doe"`
	Password string `json:"password" binding:"required" example:"password123"`
	Email    string `json:"email"`
}
