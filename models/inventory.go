package models

import (
	"context"
	"fmt"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"
)

// override table name for inventory
func (Inventory) TableName() string {
	return "inventory"
}

// inventory struct
type Inventory struct {
	ID          int
	ProductID   int
	OrderItemID int
	LocationID  int
	Damaged     bool
	RejectionID int
	ReceivedAt  time.Time `gorm:"default:null"`
	ShippedAt   time.Time `gorm:"default:null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	Product  Product
	Location Location
}

type InventoryProductStats struct {
	ProductID   int
	OnHand      int
	Available   int
	Allocated   int
	Reserved    int
	Backordered int
	OnOrder     int
	SellAhead   int
}

func GetInventoryByProductID(ctx context.Context, productID int) ([]Inventory, error) {
	var inventory []Inventory
	err := util.DBFromContext(ctx).Where("product_id = ?", productID).Where("shipped_at IS NULL").Find(&inventory).Error
	if err != nil {
		return nil, err
	}
	return inventory, nil
}

// function to batch update inventory gien a slice of inventory
func BatchUpdateInventoryOrderItemIDs(ctx context.Context, inventory []Inventory) error {
	db := util.DBFromContext(ctx)

	for i, inventoryUnit := range inventory {
		//update inventory unit
		err := db.Model(&inventoryUnit).Update("order_item_id", inventoryUnit.OrderItemID).Error
		if err != nil {
			return ErrUpdateFailed{Err: err, Object: fmt.Sprintf("inventory unit #%d", i)}
		}
	}

	return nil
}

func CreateInventory(ctx context.Context, productID int, quantity int, locationID int, damaged bool, rejectionID int) error {
	db := util.DBFromContext(ctx)
	//loop through quantity
	for i := 0; i < quantity; i++ {

		//create inventory
		inventory := Inventory{
			ProductID:   productID,
			LocationID:  locationID,
			Damaged:     damaged,
			ReceivedAt:  time.Now(),
			RejectionID: rejectionID,
		}

		//create inventory unit
		if err := db.Create(&inventory).Error; err != nil {
			return ErrCreateFailed{Err: err, Object: fmt.Sprintf("inventory #%d for product %d", i, productID)}
		}
	}

	return nil
}

func RemoveInventory(ctx context.Context, productID int, quantity int, locationID int, damaged bool) error {
	db := util.DBFromContext(ctx)

	//get inventory ids to delete
	var deletedIDs []uint
	if err := db.Model(&Inventory{}).
		Where("product_id = ? AND location_id = ?", 1, 1).
		Where("damaged = ?", damaged).
		Order("created_at DESC").
		Limit(quantity).
		Pluck("id", &deletedIDs).
		Error; err != nil {
		return ErrQueryFailed{Err: err, Object: "inventory ids to delete"}
	}

	//delete inventory
	err := db.Where("id IN (?)", deletedIDs).Delete(&Inventory{}).Error
	if err != nil {
		return ErrDeleteFailed{Err: err, Object: "inventory"}
	}

	return nil
}

func GetProductLocationsAndLevelsByProductID(ctx context.Context, productID int) ([]InventoryLocationLevel, error) {

	var result []InventoryLocationLevel

	err := util.DBFromContext(ctx).Raw(fmt.Sprintf(`
	SELECT
		locations.ID AS location_id,
		locations.NAME AS location_name,
		COUNT ( inventory.ID ) AS COUNT 
	FROM
		inventory
		INNER JOIN locations ON locations.ID = inventory.location_id 
	WHERE
		product_id = %d 
		AND shipped_at IS NULL 
	GROUP BY
		locations.ID,
		locations.NAME 
	ORDER BY
	COUNT DESC`, productID)).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil

}
