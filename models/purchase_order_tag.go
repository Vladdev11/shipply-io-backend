package models

import (
	"time"

	"gorm.io/gorm"
)

type PurchaseOrderTag struct {
	ID              int
	PurchaseOrderID int
	Tag             string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	PurchaseOrder *PurchaseOrder
}
