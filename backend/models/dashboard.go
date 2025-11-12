package models

import "gorm.io/gorm"

type Dashboard struct {
	gorm.Model
	Name       string
	Visibility bool `gorm:"not null;default:false"`
	UserId     uint
	User       User `gorm:"foreignKey:UserId"`
}
