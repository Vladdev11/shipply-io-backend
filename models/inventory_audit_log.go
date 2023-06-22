package models

import (
	"time"
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

func (ial *InventoryAuditLog) Create() error {
	return PGDB.Create(ial).Error
}
