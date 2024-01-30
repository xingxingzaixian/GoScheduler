package models

import "gorm.io/gorm"

const PasswordSaltLength = 6

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"type:varchar(32);not null;unique"`
	Password string `json:"-" gorm:"type:char(32);not null"`
	Salt     string `json:"-" gorm:"type:char(6);not null"`
	Email    string `json:"email" gorm:"type:varchar(50);not null;unique;default ''"`
	IsAdmin  string `json:"is_admin" gorm:"type:tinyint;not null;default 0"`
	Status   string `json:"status" gorm:"type:tinyint;not null;default 1"`
}
