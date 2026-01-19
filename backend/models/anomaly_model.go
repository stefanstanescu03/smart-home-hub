package models

import "gorm.io/gorm"

type AnomalyModel struct {
	gorm.Model
	Location    string `gorm:"not null"`
	Param       string `gorm:"not null"`
	NotifyEmail bool
	UserId      uint
	User        User `gorm:"foreignKey:UserId"`
	DeviceId    uint
	Device      Device `gorm:"foreignKey:DeviceId"`
}
