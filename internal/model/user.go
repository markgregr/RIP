package model

type User struct {
	UserID   uint   `gorm:"autoIncrement;primarykey" json:"user_id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password`
	Role     Role   `json:"role"`
}
type UserRegisterRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password`
}
type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UserLoginResponse struct {
	AccessToken string `json:"access_token"`
	FullName    string `json:"full_name"`
	Role        Role   `json:"role"`
}