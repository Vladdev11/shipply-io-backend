package models

import (
	"context"
	"errors"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"
)

type ShopifyLocation struct {
	ID                int    `json:"id" gorm:"primaryKey"`
	WarehouseID       int    `json:"warehouse_id"`
	StoreID           int    `json:"storeID"`
	ShopifyLocationID string `json:"shopify_location_id"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	Warehouse Warehouse `json:"warehouse" gorm:"foreignKey:WarehouseID"`
	Store     Store     `json:"store" gorm:"foreignKey:StoreID"`
}

func (sl *ShopifyLocation) Create(ctx context.Context) error {

	// Check if the location already exists
	var existingLocation ShopifyLocation
	if err := util.DBFromContext(ctx).Where("warehouse_id = ? AND store_id = ?", sl.WarehouseID, sl.StoreID).First(&existingLocation).Error; err == nil {
		return errors.New("Location already exists")
	}

	return util.DBFromContext(ctx).Create(sl).Error
}

func GetShopifyLocationsByStoreID(ctx context.Context, storeID int) ([]ShopifyLocation, error) {
	var locations []ShopifyLocation
	err := util.DBFromContext(ctx).Where("store_id = ?", storeID).Find(&locations).Error
	return locations, err
}
