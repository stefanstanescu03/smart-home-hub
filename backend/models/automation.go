package models

import "gorm.io/gorm"

type Automation struct {
	gorm.Model
	Name    string `gorm:"not null;index:name;unique"`
	Content string `gorm:"not null"`
	UserId  uint
	User    User `gorm:"foreignKey:UserId"`
}
