package models

import (
	"context"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
)

// for purpose of this inventory, we only care about on hand
type InventoryAuditLog struct {
	ID         int `gorm:"primarykey"`
	ProductID  int
	LocationID int
	Delta      int
	Note       string
	ChangedBy  int

	Product  Product
	Location Location
	User     User `gorm:"foreignKey:ChangedBy"`

	CreatedAt time.Time
}

func (ial *InventoryAuditLog) Create(ctx context.Context) error {
	return util.DBFromContext(ctx).Create(ial).Error
}

type InventoryAuditLogReturnJSON struct {
	ID         int    `json:"id"`
	ProductID  int    `json:"product_id"`
	LocationID int    `json:"location_id"`
	Delta      int    `json:"delta"`
	Note       string `json:"note"`
	ChangedBy  int    `json:"changed_by"`
}

func (ial *InventoryAuditLog) GetChangedByUser(ctx context.Context) error {
	return util.DBFromContext(ctx).Model(ial).Association("User").Find(&ial.User)
}
