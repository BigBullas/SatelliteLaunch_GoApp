package models

type User struct {
	UserId   int    `gorm:"primarykey"`
	Login    string `json:"login" binding:"required,max=64"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty" binding:"required,min=6,max=64"`
	IsAdmin  bool	`json:"is_admin"`
}
