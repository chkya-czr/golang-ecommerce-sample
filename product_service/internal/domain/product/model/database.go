package model

import "time"

type ParentProductResource struct {
	ParentId    int `gorm:"primaryKey"`
	Name        string
	Description string
	CategoryId  string `gorm:"index"`
}

type ParentProductModel struct {
	ParentProduct ParentProductResource `gorm:"embedded"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type ChildProductResource struct {
	SkuId         int `gorm:"primaryKey"`
	ParentProduct ParentProductResource
	Size          string
	Color         string
	// Price             string
	// AvailableQuantity int
}

type ChildProductModel struct {
	ChildProduct  ChildProductResource
	ParentProduct ParentProductResource `gorm:"embedded"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
