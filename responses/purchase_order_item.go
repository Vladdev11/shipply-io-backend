package responses

import (
	"context"
	"time"

	"github.com/shipply-io/shipply-io-backend/models"
)

/* ------------------------------- GetPurchaseOrderItem ------------------------------- */

// GetPurchaseOrderItemResponse represents the response body for the GetPurchaseOrderItem endpoint
type GetPurchaseOrderItemResponse struct {
	ID                      int                                                         `json:"id"`
	Ordered                 int                                                         `json:"ordered"`
	Received                int                                                         `json:"received"`
	Damaged                 int                                                         `json:"damaged"`
	UnitPrice               float64                                                     `json:"unit_price"`
	Notes                   string                                                      `json:"notes"`
	Product                 ProductResponseForGetPurchaseOrderItem                      `json:"product"`
	InventoryLocationLevels []InventoryLocationLevelResponseForGetPurchaseOrderItem     `json:"inventory_location_levels"`
	History                 []PurchaseOrderItemHistoryResponseForGetPurchaseOrderItem   `json:"history"`
	Rejections              []PurchaseOrderItemRejectionResponseForGetPurchaseOrderItem `json:"rejections"`
	ProductLots             []ProductLotResponseForGetPurchaseOrderItem                 `json:"product_lots"`
}

// ProductResponseForGetPurchaseOrderItem represents the response body for a product in the GetPurchaseOrderItem endpoint
type ProductResponseForGetPurchaseOrderItem struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Sku        string  `json:"sku"`
	Barcode    string  `json:"barcode"`
	Weight     float64 `json:"weight"`
	WeightUnit string  `json:"weight_unit"`
	Height     float64 `json:"height"`
	Width      float64 `json:"width"`
	Length     float64 `json:"length"`
}

// InventoryLocationLevelResponseForGetPurchaseOrderItem represents the response for an inventory location level in the GetPurchaseOrderItem endpoint
type InventoryLocationLevelResponseForGetPurchaseOrderItem struct {
	LocationID   int    `json:"location_id"`
	LocationName string `json:"location_name"`
	Count        int    `json:"count"`
}

// PurchaseOrderItemHistoryResponseForGetPurchaseOrderItem represents the response for a purchase order item history in the GetPurchaseOrderItem endpoint
type PurchaseOrderItemHistoryResponseForGetPurchaseOrderItem struct {
	CreatedAt     time.Time     `json:"created_at"`
	Note          string        `json:"note"`
	ChangedByUser ChangedByUser `json:"changed_by_user"`
}

// ProductLotResponseForGetPurchaseOrderItem represents the response for a product lot in the GetPurchaseOrderItem endpoint
type ProductLotResponseForGetPurchaseOrderItem struct {
	ID         int       `json:"id"`
	LotNumber  string    `json:"lot_number"`
	ExpiryDate time.Time `json:"expiry_date"`
	Active     bool      `json:"active"`
}

// PurchaseOrderItemRejectionResponseForGetPurchaseOrderItem represents the response for a purchase order item rejection in the GetPurchaseOrderItem endpoint
type PurchaseOrderItemRejectionResponseForGetPurchaseOrderItem struct {
	ID             int           `json:"id"`
	RejectedReason string        `json:"rejected_reason"`
	Note           string        `json:"note"`
	Quantity       int           `json:"quantity"`
	LocationName   string        `json:"location_name"`
	CreatedBy      ChangedByUser `json:"created_by"`
	CreatedAt      time.Time     `json:"created_at"`
}

// GenerateGetPurchaseOrderItemResponse generates the response body for the GetPurchaseOrderItem endpoint
func GenerateGetPurchaseOrderItemResponse(ctx context.Context, PurchaseOrderItem models.PurchaseOrderItem) *GetPurchaseOrderItemResponse {

	//set product
	product := ProductResponseForGetPurchaseOrderItem{
		ID:         PurchaseOrderItem.Product.ID,
		Name:       PurchaseOrderItem.Product.Name,
		Sku:        PurchaseOrderItem.Product.Sku,
		Barcode:    PurchaseOrderItem.Product.Barcode,
		Weight:     PurchaseOrderItem.Product.Weight,
		WeightUnit: PurchaseOrderItem.Product.WeightUnit,
		Height:     PurchaseOrderItem.Product.Height,
		Width:      PurchaseOrderItem.Product.Width,
		Length:     PurchaseOrderItem.Product.Length,
	}

	//set inventory location levels
	inventoryLocationsLevels := make([]InventoryLocationLevelResponseForGetPurchaseOrderItem, len(PurchaseOrderItem.InventoryLocationLevels))
	for i, inventoryLocationLevel := range PurchaseOrderItem.InventoryLocationLevels {
		inventoryLocationsLevels[i] = InventoryLocationLevelResponseForGetPurchaseOrderItem{
			LocationID:   inventoryLocationLevel.LocationID,
			LocationName: inventoryLocationLevel.LocationName,
			Count:        inventoryLocationLevel.Count,
		}
	}

	//set product history
	history := make([]PurchaseOrderItemHistoryResponseForGetPurchaseOrderItem, len(PurchaseOrderItem.History))
	for i, purchaseOrderItemHistory := range PurchaseOrderItem.History {
		history[i] = PurchaseOrderItemHistoryResponseForGetPurchaseOrderItem{
			CreatedAt: purchaseOrderItemHistory.CreatedAt,
			Note:      purchaseOrderItemHistory.Note,
			ChangedByUser: ChangedByUser{
				ID:        purchaseOrderItemHistory.CreatedByUser.ID,
				FirstName: purchaseOrderItemHistory.CreatedByUser.FirstName,
				LastName:  purchaseOrderItemHistory.CreatedByUser.LastName,
				ImageURL:  purchaseOrderItemHistory.CreatedByUser.GetAvatarFileURL(ctx),
			},
		}
	}

	//set product lots
	productLots := make([]ProductLotResponseForGetPurchaseOrderItem, len(PurchaseOrderItem.ProductLots))
	for i, productLot := range PurchaseOrderItem.ProductLots {
		productLots[i] = ProductLotResponseForGetPurchaseOrderItem{
			ID:         productLot.ID,
			LotNumber:  productLot.LotNumber,
			ExpiryDate: productLot.ExpiryDate,
			Active:     productLot.Active,
		}
	}

	//set rejections
	rejections := make([]PurchaseOrderItemRejectionResponseForGetPurchaseOrderItem, len(PurchaseOrderItem.Rejections))
	for i, rejection := range PurchaseOrderItem.Rejections {
		rejections[i] = PurchaseOrderItemRejectionResponseForGetPurchaseOrderItem{
			ID:             rejection.ID,
			RejectedReason: rejection.RejectedReaseon,
			Note:           rejection.Note,
			Quantity:       rejection.Quantity,
			LocationName:   rejection.Location.Name,
			CreatedBy: ChangedByUser{
				ID:        rejection.CreatedByUser.ID,
				FirstName: rejection.CreatedByUser.FirstName,
				LastName:  rejection.CreatedByUser.LastName,
				ImageURL:  rejection.CreatedByUser.GetAvatarFileURL(ctx),
			},
			CreatedAt: rejection.CreatedAt,
		}
	}

	return &GetPurchaseOrderItemResponse{
		ID:                      PurchaseOrderItem.ID,
		Ordered:                 PurchaseOrderItem.Ordered,
		Received:                PurchaseOrderItem.Received,
		Damaged:                 PurchaseOrderItem.Damaged,
		UnitPrice:               PurchaseOrderItem.UnitPrice,
		Notes:                   PurchaseOrderItem.Notes,
		Product:                 product,
		InventoryLocationLevels: inventoryLocationsLevels,
		History:                 history,
		Rejections:              rejections,
		ProductLots:             productLots,
	}

}

/* ------------------------------- UpdatePurchaseOrderItem ------------------------------- */

// UpdatePurchaseOrderItemResponse represents the response body for the UpdatePurchaseOrderItem endpoint
type UpdatePurchaseOrderItemResponse struct {
	ID                      int                                                            `json:"id"`
	Ordered                 int                                                            `json:"ordered"`
	Received                int                                                            `json:"received"`
	Damaged                 int                                                            `json:"damaged"`
	UnitPrice               float64                                                        `json:"unit_price"`
	Notes                   string                                                         `json:"notes"`
	Product                 ProductResponseForUpdatePurchaseOrderItem                      `json:"product"`
	InventoryLocationLevels []InventoryLocationLevelResponseForUpdatePurchaseOrderItem     `json:"inventory_location_levels"`
	History                 []PurchaseOrderItemHistoryResponseForUpdatePurchaseOrderItem   `json:"history"`
	Rejections              []PurchaseOrderItemRejectionResponseForUpdatePurchaseOrderItem `json:"rejections"`
	ProductLots             []ProductLotResponseForUpdatePurchaseOrderItem                 `json:"product_lots"`
}

// ProductResponseForUpdatePurchaseOrderItem represents the response for a product in the UpdatePurchaseOrderItem endpoint
type ProductResponseForUpdatePurchaseOrderItem struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Sku        string  `json:"sku"`
	Barcode    string  `json:"barcode"`
	Weight     float64 `json:"weight"`
	WeightUnit string  `json:"weight_unit"`
	Height     float64 `json:"height"`
	Width      float64 `json:"width"`
	Length     float64 `json:"length"`
}

// InventoryLocationLevelResponseForUpdatePurchaseOrderItem represents the response for an inventory location level in the UpdatePurchaseOrderItem endpoint
type InventoryLocationLevelResponseForUpdatePurchaseOrderItem struct {
	LocationID   int    `json:"location_id"`
	LocationName string `json:"location_name"`
	Count        int    `json:"count"`
}

// PurchaseOrderItemHistoryResponseForUpdatePurchaseOrderItem represents the response for a purchase order item history in the UpdatePurchaseOrderItem endpoint
type PurchaseOrderItemHistoryResponseForUpdatePurchaseOrderItem struct {
	CreatedAt     time.Time     `json:"created_at"`
	Note          string        `json:"note"`
	ChangedByUser ChangedByUser `json:"changed_by_user"`
}

// ProductLotResponseForUpdatePurchaseOrderItem represents the response for a product lot in the UpdatePurchaseOrderItem endpoint
type ProductLotResponseForUpdatePurchaseOrderItem struct {
	ID         int       `json:"id"`
	LotNumber  string    `json:"lot_number"`
	ExpiryDate time.Time `json:"expiry_date"`
	Active     bool      `json:"active"`
}

// PurchaseOrderItemRejectionResponseForUpdatePurchaseOrderItem represents the response for a purchase order item rejection in the UpdatePurchaseOrderItem endpoint
type PurchaseOrderItemRejectionResponseForUpdatePurchaseOrderItem struct {
	ID             int           `json:"id"`
	RejectedReason string        `json:"rejected_reason"`
	Note           string        `json:"note"`
	Quantity       int           `json:"quantity"`
	LocationName   string        `json:"location_name"`
	CreatedBy      ChangedByUser `json:"created_by"`
	CreatedAt      time.Time     `json:"created_at"`
}

// GenerateUpdatePurchaseOrderItemResponse generates the response body for the UpdatePurchaseOrderItem endpoint
func GenerateUpdatePurchaseOrderItemResponse(ctx context.Context, PurchaseOrderItem models.PurchaseOrderItem) *UpdatePurchaseOrderItemResponse {

	//set product
	product := ProductResponseForUpdatePurchaseOrderItem{
		ID:         PurchaseOrderItem.Product.ID,
		Name:       PurchaseOrderItem.Product.Name,
		Sku:        PurchaseOrderItem.Product.Sku,
		Barcode:    PurchaseOrderItem.Product.Barcode,
		Weight:     PurchaseOrderItem.Product.Weight,
		WeightUnit: PurchaseOrderItem.Product.WeightUnit,
		Height:     PurchaseOrderItem.Product.Height,
		Width:      PurchaseOrderItem.Product.Width,
		Length:     PurchaseOrderItem.Product.Length,
	}

	//set inventory location levels
	inventoryLocationsLevels := make([]InventoryLocationLevelResponseForUpdatePurchaseOrderItem, len(PurchaseOrderItem.InventoryLocationLevels))
	for i, inventoryLocationLevel := range PurchaseOrderItem.InventoryLocationLevels {
		inventoryLocationsLevels[i] = InventoryLocationLevelResponseForUpdatePurchaseOrderItem{
			LocationID:   inventoryLocationLevel.LocationID,
			LocationName: inventoryLocationLevel.LocationName,
			Count:        inventoryLocationLevel.Count,
		}
	}

	//set product history
	history := make([]PurchaseOrderItemHistoryResponseForUpdatePurchaseOrderItem, len(PurchaseOrderItem.History))
	for i, purchaseOrderItemHistory := range PurchaseOrderItem.History {
		history[i] = PurchaseOrderItemHistoryResponseForUpdatePurchaseOrderItem{
			CreatedAt: purchaseOrderItemHistory.CreatedAt,
			Note:      purchaseOrderItemHistory.Note,
			ChangedByUser: ChangedByUser{
				ID:        purchaseOrderItemHistory.CreatedByUser.ID,
				FirstName: purchaseOrderItemHistory.CreatedByUser.FirstName,
				LastName:  purchaseOrderItemHistory.CreatedByUser.LastName,
				ImageURL:  purchaseOrderItemHistory.CreatedByUser.GetAvatarFileURL(ctx),
			},
		}
	}

	//set product lots
	productLots := make([]ProductLotResponseForUpdatePurchaseOrderItem, len(PurchaseOrderItem.ProductLots))
	for i, productLot := range PurchaseOrderItem.ProductLots {
		productLots[i] = ProductLotResponseForUpdatePurchaseOrderItem{
			ID:         productLot.ID,
			LotNumber:  productLot.LotNumber,
			ExpiryDate: productLot.ExpiryDate,
			Active:     productLot.Active,
		}
	}

	//set rejections
	rejections := make([]PurchaseOrderItemRejectionResponseForUpdatePurchaseOrderItem, len(PurchaseOrderItem.Rejections))
	for i, rejection := range PurchaseOrderItem.Rejections {
		rejections[i] = PurchaseOrderItemRejectionResponseForUpdatePurchaseOrderItem{
			ID:             rejection.ID,
			RejectedReason: rejection.RejectedReaseon,
			Note:           rejection.Note,
			Quantity:       rejection.Quantity,
			LocationName:   rejection.Location.Name,
			CreatedBy: ChangedByUser{
				ID:        rejection.CreatedByUser.ID,
				FirstName: rejection.CreatedByUser.FirstName,
				LastName:  rejection.CreatedByUser.LastName,
				ImageURL:  rejection.CreatedByUser.GetAvatarFileURL(ctx),
			},
			CreatedAt: rejection.CreatedAt,
		}
	}

	return &UpdatePurchaseOrderItemResponse{
		ID:                      PurchaseOrderItem.ID,
		Ordered:                 PurchaseOrderItem.Ordered,
		Received:                PurchaseOrderItem.Received,
		Damaged:                 PurchaseOrderItem.Damaged,
		UnitPrice:               PurchaseOrderItem.UnitPrice,
		Notes:                   PurchaseOrderItem.Notes,
		Product:                 product,
		InventoryLocationLevels: inventoryLocationsLevels,
		History:                 history,
		Rejections:              rejections,
		ProductLots:             productLots,
	}

}
