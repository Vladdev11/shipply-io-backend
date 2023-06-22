package models

import (
	"errors"
)

type PickSessionOrder struct {
	ID                      int  `json:"id" gorm:"primary_key"`
	OrderID                 int  `json:"order_id"`
	PickSessionID           int  `json:"pick_session_id"`
	LocationID              *int `json:"location_id"`
	Picked                  bool `json:"picked"`
	Shipped                 bool `json:"shipped"`
	HasError                bool `json:"has_error"`
	PickSessionOrderErrorID *int `json:"pick_session_order_error_id"`
	ShippingRateID          int  `json:"shipping_rate_id"`

	Order                 Order                  `json:"order"`
	PickSession           PickSession            `json:"pick_session"`
	Location              *Location              `json:"location"`
	PickSessionOrderError *PickSessionOrderError `json:"pick_session_order_error"`
	PickSessionOrderItems []PickSessionOrderItem `json:"pick_session_order_items"`
	ShippingRate          ShippingRate           `json:"shipping_rate"`
}

func (pso *PickSessionOrder) Create() error {
	return PGDB.Create(pso).Error
}

func (pso *PickSessionOrder) Update() error {
	return PGDB.Save(pso).Error
}

func (pso *PickSessionOrder) GetShippingRate() error {

	var shippingRate ShippingRate
	err := PGDB.Where("id = ?", pso.ShippingRateID).First(&shippingRate).Error
	if err != nil {
		return err
	}

	pso.ShippingRate = shippingRate

	return nil

}

func (pso *PickSessionOrder) GetPickSession() error {
	return PGDB.Where("id = ?", pso.PickSessionID).First(&pso.PickSession).Error
}

func (pso *PickSessionOrder) GetOrder() (Order, error) {
	var order Order
	err := PGDB.Where("id = ?", pso.OrderID).First(&order).Error
	return order, err
}

func (pso *PickSessionOrder) GetPickSessionOrderItems() error {
	return PGDB.Where("pick_session_order_id = ?", pso.ID).Find(&pso.PickSessionOrderItems).Error
}

// TODO clean up this function
func (pso *PickSessionOrder) CreatePickSessionOrderItems(warehouseID int) error {
	var result []struct {
		ID         int
		Allocated  int
		ProductID  int
		LocationID int
	}
	err := PGDB.Raw(`SELECT order_items.id, order_items.allocated, order_items.product_id, locations.id as location_id
	FROM order_items
	JOIN locations ON locations.warehouse_id = ? AND locations.pickable = true AND locations.is_tote IS NOT TRUE
	JOIN inventory ON inventory.product_id = order_items.product_id
	JOIN (
		SELECT inventory.product_id, inventory.location_id, COUNT(*) AS available_count
		FROM inventory
		WHERE inventory.order_item_id IS NULL
		GROUP BY inventory.product_id, inventory.location_id
	) AS inventory_counts ON inventory_counts.product_id = order_items.product_id AND inventory_counts.location_id = locations.id AND inventory_counts.available_count >= order_items.allocated
	WHERE order_items.quantity = order_items.allocated
	AND order_items.order_id = ?
	GROUP BY order_items.id, locations.id
	`, warehouseID, pso.OrderID).Scan(&result).Error
	if err != nil {
		return errors.New("failed to get order items for pick session order")
	}

	for _, orderItem := range result {
		pickSessionOrderItem := PickSessionOrderItem{
			PickSessionID:      pso.PickSessionID,
			PickSessionOrderID: pso.ID,
			OrderItemID:        orderItem.ID,
			QuantityToPick:     orderItem.Allocated,
			QuantityPicked:     0,
			ProductID:          orderItem.ProductID,
			LocationID:         orderItem.LocationID,
		}
		err := pickSessionOrderItem.Create()
		if err != nil {
			return err
		}

		for i := 0; i < orderItem.Allocated; i++ {

			var inventoryItem Inventory
			err := PGDB.Where("product_id = ? AND location_id = ? AND order_item_id IS NULL", orderItem.ProductID, orderItem.LocationID).First(&inventoryItem).Error
			if err != nil {
				return errors.New("failed to get inventory item for pick session order item")
			}

			err = PGDB.Model(&inventoryItem).Update("order_item_id", orderItem.ID).Error
			if err != nil {
				return errors.New("failed to update inventory item for pick session order item")
			}
		}

	}

	return nil
}

func GetPickSessionOrderByID(id int) (PickSessionOrder, error) {
	var pickSessionOrder PickSessionOrder
	err := PGDB.Where("id = ?", id).First(&pickSessionOrder).Error
	return pickSessionOrder, err
}
