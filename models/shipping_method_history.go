package models

import (
	"time"

	"gorm.io/gorm"
)

type ShippingMethodHistory struct {
	ID               int
	ShippingMethodID int
	Note             string
	CreatedBy        int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	ShippingMethod *ShippingMethod
	CreatedByUser  *User `gorm:"foreignKey:CreatedBy"`
}

func (smh *ShippingMethodHistory) Create() error {
	return PGDB.Create(smh).Error
}
