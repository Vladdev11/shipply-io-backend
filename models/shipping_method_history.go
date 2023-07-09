package models

import (
	"context"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
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

func (smh *ShippingMethodHistory) Create(ctx context.Context) error {
	return util.DBFromContext(ctx).Create(smh).Error
}
