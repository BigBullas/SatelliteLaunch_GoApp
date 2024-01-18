package models

type User struct {
	UserId   int `gorm:"primarykey"`
	Login    string  `json:"login" binding:"required,max=64"`
	Email    string  `json:"name,omitempty"`
	Password string  `json:"password,omitempty" binding:"required,min=8,max=64"`
	IsAdmin  bool
}
