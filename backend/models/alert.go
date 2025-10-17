package models

import "gorm.io/gorm"

type Alert struct {
	gorm.Model
	Subject   string `gorm:"not null;index:subject"`
	Content   string `gorm:"not null"`
	Condition string `gorm:"not null"`
	UserId    uint
	User      User `gorm:"foreignKey:UserId"`
}
