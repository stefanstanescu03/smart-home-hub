package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"not null;index:username;unique"`
	Password string `gorm:"not null"`
	Email    string `gorm:"size:250;index;not null;index:email;unique"`
	Role     string `gorm:"type:varchar(50);not null;default:'user'"`
}
