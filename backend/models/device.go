package models

import "gorm.io/gorm"

type Device struct {
	gorm.Model
	Name         string `gorm:"not null;index:name;unique"`
	Ident        string `gorm:"not null;unique"`
	Csv_location string `gorm:"not null;unique"`
	Cloud_api    string
	Visibility   bool `gorm:"not null;default:false"`
	UserId       uint
	User         User `gorm:"foreignKey:UserId"`
}
