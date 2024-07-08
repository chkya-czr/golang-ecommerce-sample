package model

import "gorm.io/gorm"

type ParentProduct struct {
	Id          int
	Name        string
	Description string
}

type ParentProductModel struct {
	gorm.Model
	ParentProduct ParentProduct `gorm:"embedded"`
}
