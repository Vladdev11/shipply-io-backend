package tasks

import (
	"context"
	"errors"

	"github.com/shipply-io/shipply-io-backend/models"
	"gorm.io/gorm"
)

func AllocateInventoryByProduct(ctx context.Context, productID int) error {

	//get product
	product, err := models.GetProductByID(ctx, productID)
	if err != nil {
		return errors.New("failed to find product")
	}

	//get inventory by product
	inventory, err := models.GetInventoryByProductID(ctx, product.ID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return errors.New("failed to find inventory")
	}

	//get order items by product that are not shipped
	orderItems, err := models.GetOrderItemsToShipByProductID(ctx, product.ID)
	if err != nil {
		return errors.New("failed to find order items")
	}

	//skip if no order items
	if len(orderItems) == 0 {
		//update order item inventory levels
		models.UpdateProductInventoryLevelsByProductID(ctx, product.ID)
		return nil
	}

	//order items allocation logic
	orderItems = models.SortOrderItemsForAllocation(orderItems)

	//set map of orderItemID to an inventory item
	orderItemAllocationMap := make(map[int]int)

	//go through inventory and allocate to order items
	for i := range inventory {

		//skip if order item slice is not long enough
		if i >= len(orderItems) {
			break
		}

		//increment order item allocation
		orderItemAllocationMap[int(orderItems[i].ID)]++

	}

	//go through order items and set allocation count
	for i, orderItem := range orderItems {
		//set allocation count
		orderItems[i].Allocated = orderItemAllocationMap[int(orderItem.ID)]
		orderItems[i].Backordered = orderItem.Quantity - orderItem.QuantityShipped - orderItems[i].Allocated
	}

	//batch update order items
	err = models.BatchUpdateOrderItemsAllocationCount(ctx, orderItems)
	if err != nil {
		return errors.New("failed to update order items")
	}

	//update order item inventory levels
	models.UpdateProductInventoryLevelsByProductID(ctx, product.ID)

	return nil
}
