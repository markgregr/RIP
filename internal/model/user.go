package model

type User struct {
	UserID   uint   `gorm:"autoIncrement;primarykey" json:"user_id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}