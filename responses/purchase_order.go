package responses

import (
	"context"
	"time"

	"github.com/shipply-io/shipply-io-backend/models"
)

/* ----------------------------- GetPurchaseOrder ---------------------------- */

// GetPurchaseOrderResponse represents the expected response body for the GetPurchaseOrder endpoint
type GetPurchaseOrderResponse struct {
	ID             int    `json:"id"`
	PONumber       string `json:"po_number"`
	VendorID       int    `json:"vendor_id"`
	WarehouseID    int    `json:"warehouse_id"`
	TrackingNumber string `json:"tracking_number"`
	TrackingURL    string `json:"tracking_url"`
	WarehouseNotes string `json:"warehouse_notes"`
	UniqueItems    int    `json:"unique_items"`
	TotalItems     int    `json:"total_items"`

	ExpectedDate time.Time `json:"expected_date"`
	ShipDate     time.Time `json:"ship_date"`
	ClosedDate   time.Time `json:"closed_date"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Client  ClientResponseForGetPurchaseOrder    `json:"client"`
	Tags    []string                             `json:"tags"`
	Status  StatusResponseForGetPurchaseOrder    `json:"status"`
	Items   []ItemResponseForGetPurchaseOrder    `json:"items"`
	History []HistoryResponseForGetPurchaseOrder `json:"history"`
}

// StatusResponseForGetPurchaseOrder represents the expected response body for an individual status in the GetPurchaseOrder endpoint
type StatusResponseForGetPurchaseOrder struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	StatusColor string `json:"status_color"`
	TextColor   string `json:"text_color"`
}

// ItemResponseForGetPurchaseOrder represents the expected response body for an individual item in the GetPurchaseOrder endpoint
type ItemResponseForGetPurchaseOrder struct {
	ID        int                                               `json:"id"`
	Ordered   int                                               `json:"ordered"`
	Received  int                                               `json:"received"`
	Damaged   int                                               `json:"damaged"`
	UnitPrice float64                                           `json:"unit_price"`
	Notes     string                                            `json:"notes"`
	Product   ProductResponseForItemResponseForGetPurchaseOrder `json:"product"`
}

// ProductResponseForGetPurchaseOrder represents the expected response body for an individual product in the GetPurchaseOrder endpoint
type ProductResponseForItemResponseForGetPurchaseOrder struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Sku         string     `json:"sku"`
	Barcode     string     `json:"barcode"`
	Active      bool       `json:"active"`
	Weight      Weight     `json:"weight"`
	Dimensions  Dimensions `json:"dimensions"`
	LastCounted time.Time  `json:"last_counted"`
}

// HistoryResponseForGetPurchaseOrder represents the expected response body for an individual history in the GetPurchaseOrder endpoint
type HistoryResponseForGetPurchaseOrder struct {
	Note      string        `json:"note"`
	CreatedAt time.Time     `json:"created_at"`
	CreatedBy ChangedByUser `json:"created_by"`
}

type ClientResponseForGetPurchaseOrder struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GenerateGetPurchaseOrderResponse generates the response body for the GetPurchaseOrder endpoint
func GenerateGetPurchaseOrderResponse(ctx context.Context, purchaseOrder models.PurchaseOrder) *GetPurchaseOrderResponse {

	statusResponse := StatusResponseForGetPurchaseOrder{
		ID:          purchaseOrder.Status.ID,
		Name:        purchaseOrder.Status.Name,
		StatusColor: purchaseOrder.Status.StatusColor,
		TextColor:   purchaseOrder.Status.TextColor,
	}

	client := ClientResponseForGetPurchaseOrder{
		ID:   purchaseOrder.Client.ID,
		Name: purchaseOrder.Client.Name,
	}

	itemResponses := []ItemResponseForGetPurchaseOrder{}
	for _, item := range purchaseOrder.Items {

		weight := Weight{
			Value: item.Product.Weight,
			Unit:  item.Product.WeightUnit,
		}

		dimensions := Dimensions{
			Length: item.Product.Length,
			Width:  item.Product.Width,
			Height: item.Product.Height,
		}

		productResponse := ProductResponseForItemResponseForGetPurchaseOrder{
			ID:          item.Product.ID,
			Name:        item.Product.Name,
			Sku:         item.Product.Sku,
			Barcode:     item.Product.Barcode,
			Active:      item.Product.Active,
			Weight:      weight,
			Dimensions:  dimensions,
			LastCounted: item.Product.LastCounted,
		}

		itemResponses = append(itemResponses, ItemResponseForGetPurchaseOrder{
			ID:        item.ID,
			Ordered:   item.Ordered,
			Received:  item.Received,
			Damaged:   item.Damaged,
			UnitPrice: item.UnitPrice,
			Notes:     item.Notes,
			Product:   productResponse,
		})
	}

	historyResponses := []HistoryResponseForGetPurchaseOrder{}
	for _, history := range purchaseOrder.History {

		changedByUser := ChangedByUser{
			ID:        history.CreatedByUser.ID,
			FirstName: history.CreatedByUser.FirstName,
			LastName:  history.CreatedByUser.LastName,
			ImageURL:  history.CreatedByUser.GetAvatarFileURL(ctx),
		}

		historyResponses = append(historyResponses, HistoryResponseForGetPurchaseOrder{
			Note:      history.Note,
			CreatedAt: history.CreatedAt,
			CreatedBy: changedByUser,
		})
	}

	tagResponses := []string{}
	for _, tag := range purchaseOrder.Tags {
		tagResponses = append(tagResponses, tag.Tag)
	}

	return &GetPurchaseOrderResponse{
		ID:             purchaseOrder.ID,
		PONumber:       purchaseOrder.PONumber,
		VendorID:       purchaseOrder.VendorID,
		WarehouseID:    purchaseOrder.WarehouseID,
		TrackingNumber: purchaseOrder.TrackingNumber,
		TrackingURL:    purchaseOrder.TrackingURL,
		WarehouseNotes: purchaseOrder.WarehouseNotes,
		UniqueItems:    purchaseOrder.UniqueItems,
		TotalItems:     purchaseOrder.TotalItems,

		ExpectedDate: purchaseOrder.ExpectedDate,
		ShipDate:     purchaseOrder.ShipDate,
		ClosedDate:   purchaseOrder.ClosedDate,
		CreatedAt:    purchaseOrder.CreatedAt,
		UpdatedAt:    purchaseOrder.UpdatedAt,

		Tags:    tagResponses,
		Status:  statusResponse,
		Items:   itemResponses,
		History: historyResponses,
		Client:  client,
	}

}

/* --------------------------- ListPurchaseOrders --------------------------- */

// ListPurchaseOrdersResponse represents the expected response body for the ListPurchaseOrders endpoint
type ListPurchaseOrdersResponse struct {
	TotalCount    int                                          `json:"total_count"`
	FilteredCount int                                          `json:"filtered_count"`
	Data          []PurchaseOrderResponseForListPurchaseOrders `json:"data"`
}

// PurchaseOrderResponseForListPurchaseOrders represents the expected response body for an individual purchase order in the ListPurchaseOrders endpoint
type PurchaseOrderResponseForListPurchaseOrders struct {
	ID        int                                 `json:"id"`
	PONumber  string                              `json:"po_number"`
	CreatedAt time.Time                           `json:"created_at"`
	Status    StatusResponseForListPurchaseOrders `json:"status"`
}

// StatusResponseForListPurchaseOrders represents the expected response body for an individual status in the ListPurchaseOrders endpoint
type StatusResponseForListPurchaseOrders struct {
	Name        string `json:"name"`
	StatusColor string `json:"status_color"`
	TextColor   string `json:"text_color"`
}

// GenerateListPurchaseOrdersResponse generates the response body for the ListPurchaseOrders endpoint
func GenerateListPurchaseOrdersResponse(purchaseOrders []models.PurchaseOrder, totalCount int, filteredCount int) *ListPurchaseOrdersResponse {

	purchaseOrderResponses := []PurchaseOrderResponseForListPurchaseOrders{}
	for _, purchaseOrder := range purchaseOrders {

		statusResponse := StatusResponseForListPurchaseOrders{
			Name:        purchaseOrder.Status.Name,
			StatusColor: purchaseOrder.Status.StatusColor,
			TextColor:   purchaseOrder.Status.TextColor,
		}

		purchaseOrderResponses = append(purchaseOrderResponses, PurchaseOrderResponseForListPurchaseOrders{
			ID:        purchaseOrder.ID,
			PONumber:  purchaseOrder.PONumber,
			CreatedAt: purchaseOrder.CreatedAt,
			Status:    statusResponse,
		})
	}

	return &ListPurchaseOrdersResponse{
		TotalCount:    totalCount,
		FilteredCount: filteredCount,
		Data:          purchaseOrderResponses,
	}

}

/* --------------------------- CreatePurchaseOrder -------------------------- */

// CreatePurchaseOrderResponse represents the expected response body for the CreatePurchaseOrder endpoint
type CreatePurchaseOrderResponse struct {
	ID             int       `json:"id"`
	PONumber       string    `json:"po_number"`
	ExpectedDate   time.Time `json:"expected_date"`
	ShipDate       time.Time `json:"ship_date"`
	ClosedDate     time.Time `json:"closed_date"`
	VendorID       int       `json:"vendor_id"`
	WarehouseID    int       `json:"warehouse_id"`
	TrackingNumber string    `json:"tracking_number"`
	TrackingURL    string    `json:"tracking_url"`

	Status StatusResponseForCreatePurchaseOrder `json:"status"`
}

// StatusResponseForCreatePurchaseOrder represents the expected response body for an individual status in the CreatePurchaseOrder endpoint
type StatusResponseForCreatePurchaseOrder struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	StatusColor string `json:"status_color"`
	TextColor   string `json:"text_color"`
}

// GenerateCreatePurchaseOrderResponse generates the response body for the CreatePurchaseOrder endpoint
func GenerateCreatePurchaseOrderResponse(purchaseOrder models.PurchaseOrder) *CreatePurchaseOrderResponse {

	statusResponse := StatusResponseForCreatePurchaseOrder{
		ID:          purchaseOrder.Status.ID,
		Name:        purchaseOrder.Status.Name,
		StatusColor: purchaseOrder.Status.StatusColor,
		TextColor:   purchaseOrder.Status.TextColor,
	}

	return &CreatePurchaseOrderResponse{
		ID:             purchaseOrder.ID,
		PONumber:       purchaseOrder.PONumber,
		ExpectedDate:   purchaseOrder.ExpectedDate,
		ShipDate:       purchaseOrder.ShipDate,
		ClosedDate:     purchaseOrder.ClosedDate,
		VendorID:       purchaseOrder.VendorID,
		WarehouseID:    purchaseOrder.WarehouseID,
		TrackingNumber: purchaseOrder.TrackingNumber,
		TrackingURL:    purchaseOrder.TrackingURL,

		Status: statusResponse,
	}

}

/* --------------------------- UpdatePurchaseOrder -------------------------- */

// UpdatePurchaseOrderResponse represents the expected response body for the UpdatePurchaseOrder endpoint
type UpdatePurchaseOrderResponse struct {
	ID             int    `json:"id"`
	PONumber       string `json:"po_number"`
	VendorID       int    `json:"vendor_id"`
	WarehouseID    int    `json:"warehouse_id"`
	TrackingNumber string `json:"tracking_number"`
	TrackingURL    string `json:"tracking_url"`
	WarehouseNotes string `json:"warehouse_notes"`
	UniqueItems    int    `json:"unique_items"`
	TotalItems     int    `json:"total_items"`

	ExpectedDate time.Time `json:"expected_date"`
	ShipDate     time.Time `json:"ship_date"`
	ClosedDate   time.Time `json:"closed_date"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Tags    []string                                `json:"tags"`
	Status  StatusResponseForUpdatePurchaseOrder    `json:"status"`
	Items   []ItemResponseForUpdatePurchaseOrder    `json:"items"`
	History []HistoryResponseForUpdatePurchaseOrder `json:"history"`
}

// StatusResponseForUpdatePurchaseOrder represents the expected response body for an individual status in the UpdatePurchaseOrder endpoint
type StatusResponseForUpdatePurchaseOrder struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	StatusColor string `json:"status_color"`
	TextColor   string `json:"text_color"`
}

// ItemResponseForUpdatePurchaseOrder represents the expected response body for an individual item in the UpdatePurchaseOrder endpoint
type ItemResponseForUpdatePurchaseOrder struct {
	ID        int                                                  `json:"id"`
	Ordered   int                                                  `json:"ordered"`
	Received  int                                                  `json:"received"`
	Damaged   int                                                  `json:"damaged"`
	UnitPrice float64                                              `json:"unit_price"`
	Notes     string                                               `json:"notes"`
	Product   ProductResponseForItemResponseForUpdatePurchaseOrder `json:"product"`
}

// ProductResponseForItemResponseForUpdatePurchaseOrder represents the expected response body for an individual product in the UpdatePurchaseOrder endpoint
type ProductResponseForItemResponseForUpdatePurchaseOrder struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Sku         string     `json:"sku"`
	Barcode     string     `json:"barcode"`
	Active      bool       `json:"active"`
	Weight      Weight     `json:"weight"`
	Dimensions  Dimensions `json:"dimensions"`
	LastCounted time.Time  `json:"last_counted"`
}

// HistoryResponseForUpdatePurchaseOrder represents the expected response body for an individual history in the UpdatePurchaseOrder endpoint
type HistoryResponseForUpdatePurchaseOrder struct {
	Note      string        `json:"note"`
	CreatedAt time.Time     `json:"created_at"`
	CreatedBy ChangedByUser `json:"created_by"`
}

// GenerateUpdatePurchaseOrderResponse generates the response body for the UpdatePurchaseOrder endpoint
func GenerateUpdatePurchaseOrderResponse(ctx context.Context, purchaseOrder models.PurchaseOrder) *UpdatePurchaseOrderResponse {

	statusResponse := StatusResponseForUpdatePurchaseOrder{
		ID:          purchaseOrder.Status.ID,
		Name:        purchaseOrder.Status.Name,
		StatusColor: purchaseOrder.Status.StatusColor,
		TextColor:   purchaseOrder.Status.TextColor,
	}

	itemResponses := []ItemResponseForUpdatePurchaseOrder{}
	for _, item := range purchaseOrder.Items {

		weight := Weight{
			Value: item.Product.Weight,
			Unit:  item.Product.WeightUnit,
		}

		dimensions := Dimensions{
			Length: item.Product.Length,
			Width:  item.Product.Width,
			Height: item.Product.Height,
		}

		productResponse := ProductResponseForItemResponseForUpdatePurchaseOrder{
			ID:          item.Product.ID,
			Name:        item.Product.Name,
			Sku:         item.Product.Sku,
			Barcode:     item.Product.Barcode,
			Active:      item.Product.Active,
			Weight:      weight,
			Dimensions:  dimensions,
			LastCounted: item.Product.LastCounted,
		}

		itemResponses = append(itemResponses, ItemResponseForUpdatePurchaseOrder{
			ID:        item.ID,
			Ordered:   item.Ordered,
			Received:  item.Received,
			Damaged:   item.Damaged,
			UnitPrice: item.UnitPrice,
			Notes:     item.Notes,
			Product:   productResponse,
		})
	}

	historyResponses := []HistoryResponseForUpdatePurchaseOrder{}
	for _, history := range purchaseOrder.History {

		changedByUser := ChangedByUser{
			ID:        history.CreatedByUser.ID,
			FirstName: history.CreatedByUser.FirstName,
			LastName:  history.CreatedByUser.LastName,
			ImageURL:  history.CreatedByUser.GetAvatarFileURL(ctx),
		}

		historyResponses = append(historyResponses, HistoryResponseForUpdatePurchaseOrder{
			Note:      history.Note,
			CreatedAt: history.CreatedAt,
			CreatedBy: changedByUser,
		})
	}

	tagResponses := []string{}
	for _, tag := range purchaseOrder.Tags {
		tagResponses = append(tagResponses, tag.Tag)
	}

	return &UpdatePurchaseOrderResponse{
		ID:             purchaseOrder.ID,
		PONumber:       purchaseOrder.PONumber,
		VendorID:       purchaseOrder.VendorID,
		WarehouseID:    purchaseOrder.WarehouseID,
		TrackingNumber: purchaseOrder.TrackingNumber,
		TrackingURL:    purchaseOrder.TrackingURL,
		WarehouseNotes: purchaseOrder.WarehouseNotes,
		UniqueItems:    purchaseOrder.UniqueItems,
		TotalItems:     purchaseOrder.TotalItems,

		ExpectedDate: purchaseOrder.ExpectedDate,
		ShipDate:     purchaseOrder.ShipDate,
		ClosedDate:   purchaseOrder.ClosedDate,
		CreatedAt:    purchaseOrder.CreatedAt,
		UpdatedAt:    purchaseOrder.UpdatedAt,

		Tags:    tagResponses,
		Status:  statusResponse,
		Items:   itemResponses,
		History: historyResponses,
	}

}
