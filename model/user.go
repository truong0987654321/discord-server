package model

import "mime/multipart"

type User struct {
	BaseModel
	Username string `gorm:"not null" json:"username"`
	Email    string `gorm:"not null;uniqueIndex" json:"email"`
	Password string `gorm:"not null" json:"-"`
	Image    string `json:"image"`
	IsOnline bool   `gorm:"index;default:true" json:"isOnline"`
}

type UserService interface {
	Register(user *User) (*User, error)
	Login(email, password string) (*User, error)
	Get(id string) (*User, error)
	ChangeAvatar(header *multipart.FileHeader, directory string) (string, error)
	DeleteImage(key string) error
}

type UserRepository interface {
	Create(user *User) (*User, error)
	FindByEmail(email string) (*User, error)
	FindByID(id string) (*User, error)
}
