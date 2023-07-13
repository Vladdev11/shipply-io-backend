package models

import (
	"time"

	"gorm.io/gorm"
)

type OrderTag struct {
	ID      int
	OrderID int
	Tag     string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Order *Order
}
