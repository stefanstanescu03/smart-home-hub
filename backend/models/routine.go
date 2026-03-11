package models

import "gorm.io/gorm"

type Routine struct {
	gorm.Model
	Name      string `gorm:"not null;index:name;unique"`
	DeviceId  uint
	Device    Device `gorm:"foreignKey:DeviceId"`
	Payload   string `gorm:"not null"`
	StartTime string `gorm:"not null"`
	UserId    uint
	User      User `gorm:"foreignKey:UserId"`
}
