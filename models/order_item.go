package models

import (
	"time"

	"gorm.io/gorm"
)

type OrderItem struct {
	ID                           int
	APIID                        string
	OrderID                      int
	Name                         string
	Status                       string
	Backordered                  int
	ProductID                    *int
	MarketplaceFulfillmentStatus string
	Quantity                     int
	ItemPrice                    float64
	QuantityShipped              int `gorm:"default:0"`
	Allocated                    int `gorm:"default:0"`
	Priority                     int
	//TODO remove shipped / convert to fulfilled
	Shipped   bool
	ShippedAt time.Time
	Fulfilled bool
	Sku       string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	Order   *Order
	Product *Product
}

type OrderItemReturnJSON struct {
	ID                           int       `json:"id"`
	OrderID                      int       `json:"order_id,omitempty"`
	ProductID                    *int      `json:"product_id,omitempty"`
	Name                         string    `json:"name"`
	Status                       string    `json:"status"`
	Backordered                  int       `json:"backordered"`
	MarketplaceFulfillmentStatus string    `json:"marketplace_fulfillment_status"`
	Quantity                     int       `json:"quantity"`
	ItemPrice                    float64   `json:"item_price"`
	QuantityShipped              int       `json:"quantity_shipped"`
	Allocated                    int       `json:"allocated"`
	Priority                     int       `json:"priority"`
	Shipped                      bool      `json:"shipped"`
	ShippedAt                    time.Time `json:"shipped_at"`
	Fulfilled                    bool      `json:"fulfilled"`
	Sku                          string    `json:"sku"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Order   *OrderReturnJSON   `json:"order,omitempty"`
	Product *ProductReturnJSON `json:"product,omitempty"`
}

func (oi *OrderItem) ConvertToReturnJSON() OrderItemReturnJSON {
	returnJSON := OrderItemReturnJSON{
		ID:                           oi.ID,
		OrderID:                      oi.OrderID,
		ProductID:                    oi.ProductID,
		Name:                         oi.Name,
		Status:                       oi.Status,
		Backordered:                  oi.Backordered,
		MarketplaceFulfillmentStatus: oi.MarketplaceFulfillmentStatus,
		Quantity:                     oi.Quantity,
		ItemPrice:                    oi.ItemPrice,
		QuantityShipped:              oi.QuantityShipped,
		Allocated:                    oi.Allocated,
		Priority:                     oi.Priority,
		Shipped:                      oi.Shipped,
		ShippedAt:                    oi.ShippedAt,
		Fulfilled:                    oi.Fulfilled,
		Sku:                          oi.Sku,
		CreatedAt:                    oi.CreatedAt,
		UpdatedAt:                    oi.UpdatedAt,
	}

	if oi.Order != nil {
		returnJSON.Order = oi.Order.ConvertToReturnJSON()
	}

	if oi.Product != nil {
		returnJSON.Product = oi.Product.ConvertToReturnJSON()
	}

	return returnJSON

}

func GetOrderItemsByOrderID(orderID int) ([]OrderItem, error) {
	var orderItems []OrderItem
	err := PGDB.Where("order_id = ?", orderID).Find(&orderItems).Error
	if err != nil {
		return nil, err
	}
	return orderItems, nil
}

func GetOrderItemsToShipByProductID(productID int) ([]OrderItem, error) {
	var orderItems []OrderItem
	err := PGDB.Where("product_id = ?", productID).Where("quantity > quantity_shipped").Find(&orderItems).Error
	if err != nil {
		return nil, err
	}
	return orderItems, nil
}

func BatchUpdateOrderItemsAllocationCount(orderItems []OrderItem) error {
	//start transaction
	tx := PGDB.Begin()
	for _, orderItem := range orderItems {

		//update the columns
		err := tx.Model(&orderItem).Update("allocated", orderItem.Allocated).Error
		if err != nil {
			//rollback transaction if error
			tx.Rollback()
			return err
		}
	}
	//commit transaction
	return tx.Commit().Error
}

func GetOrderItemsByProductID(productID int) ([]OrderItem, error) {
	var orderItems []OrderItem
	err := PGDB.Where("product_id = ?", productID).Find(&orderItems).Error
	if err != nil {
		return nil, err
	}
	return orderItems, nil
}

func (oi *OrderItem) Create() error {
	err := PGDB.Create(oi).Error
	if err != nil {
		return err
	}
	return nil
}

func GetOrderItemByOrderIDAndAPIID(orderID int, apiID string) (*OrderItem, error) {
	orderItem := OrderItem{}
	err := PGDB.Where("order_id = ? AND api_id = ?", orderID, apiID).First(&orderItem).Error
	if err != nil {
		return nil, err
	}
	return &orderItem, nil
}
func (oi *OrderItem) Update() error {
	return PGDB.Save(oi).Error
}
