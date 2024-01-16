package models

type User struct {
	UserId   int `gorm:"primarykey"`
	Login    string
	Email    string
	Password string
	IsAdmin  bool
}
