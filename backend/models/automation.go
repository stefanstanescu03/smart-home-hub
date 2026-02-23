package models

import "gorm.io/gorm"

type Automation struct {
	gorm.Model
	Name      string `gorm:"not null;index:name;unique"`
	Device1Id uint
	Device1   Device `gorm:"foreignKey:Device1Id"`
	Device2Id uint
	Device2   Device `gorm:"foreignKey:Device2Id"`
	Payload   string `gorm:"not null"`
	Condition string `gorm:"not null"`
	UserId    uint
	User      User `gorm:"foreignKey:UserId"`
}
