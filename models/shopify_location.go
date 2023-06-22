package models

import (
	"errors"
	"time"

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

func (sl *ShopifyLocation) Create() error {

	// Check if the location already exists
	var existingLocation ShopifyLocation
	if err := PGDB.Where("warehouse_id = ? AND store_id = ?", sl.WarehouseID, sl.StoreID).First(&existingLocation).Error; err == nil {
		return errors.New("Location already exists")
	}

	return PGDB.Create(sl).Error
}

func GetShopifyLocationsByStoreID(storeID int) ([]ShopifyLocation, error) {
	var locations []ShopifyLocation
	err := PGDB.Where("store_id = ?", storeID).Find(&locations).Error
	return locations, err
}
