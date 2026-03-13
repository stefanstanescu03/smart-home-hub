package models

import "gorm.io/gorm"

type Widget struct {
	gorm.Model
	Widgettype  string `gorm:"check:widget_type_check,widgettype IN ('ta', 'ca', 'btn', 'sw')"`
	Payload     string
	Payload2    string
	Label       string
	Channel     string
	DeviceId    uint
	Device      Device `gorm:"foreignKey:DeviceId"`
	DashboardId uint
	Dashboard   Dashboard `gorm:"foreignKey:DashboardId"`
}
