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

	//start transaction
	tx := util.DBFromContext(ctx).Begin()
	for _, inventoryUnit := range inventory {
		//update inventory unit
		err := tx.Model(&inventoryUnit).Update("order_item_id", inventoryUnit.OrderItemID).Error
		if err != nil {
			//rollback transaction if error
			tx.Rollback()
		}
	}
	//commit transaction
	err := tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func CreateInventory(ctx context.Context, productID int, quantity int, locationID int, damaged bool, rejectionID int) error {

	//start transaction
	tx := util.DBFromContext(ctx).Begin()

	//flag to indicate if an error has occurred
	errorOccurred := false

	//define err
	var err error

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
		err := tx.Create(&inventory).Error
		if err != nil {
			errorOccurred = true
			break
		}
	}

	if errorOccurred {
		//rollback transaction if error
		tx.Rollback()
		return err
	}

	//commit transaction
	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func RemoveInventory(ctx context.Context, productID int, quantity int, locationID int, damaged bool) error {

	//start transaction
	tx := util.DBFromContext(ctx).Begin()

	//get inventory ids to delete
	var deletedIDs []uint
	util.DBFromContext(ctx).Model(&Inventory{}).Where("product_id = ? AND location_id = ?", 1, 1).Where("damaged = ?", damaged).Order("created_at DESC").Limit(quantity).Pluck("id", &deletedIDs)

	//delete inventory
	err := tx.Where("id IN (?)", deletedIDs).Delete(&Inventory{}).Error
	if err != nil {
		//rollback transaction if error
		tx.Rollback()
		return err
	}

	//commit transaction
	return tx.Commit().Error

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
