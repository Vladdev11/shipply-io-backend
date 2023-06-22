package models

import (
	"sort"
)

func SortOrderItemsForAllocation(orderItems []OrderItem) []OrderItem {

	//create new slice
	singleOrderItems := []OrderItem{}

	//add to new slice for every instance of order item quantity
	for _, orderItem := range orderItems {
		//quantity is amount that we actually need inventory for
		quantity := orderItem.Quantity - orderItem.QuantityShipped
		//for each unit of an inventory item we add a new struct to the slices
		for j := 0; j < quantity; j++ {
			newOrderItem := orderItem
			//set quantity to 1 since it is a single occurence
			newOrderItem.Quantity = 1
			singleOrderItems = append(singleOrderItems, newOrderItem)
		}
	}

	//sort order items by allocation by created at date
	sort.Slice(singleOrderItems, func(i, j int) bool {
		return singleOrderItems[i].CreatedAt.Before(singleOrderItems[j].CreatedAt)
	})

	return singleOrderItems
}
