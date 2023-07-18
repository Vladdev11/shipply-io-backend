package responses

import "github.com/shipply-io/shipply-io-backend/models"

/* ------------------------ ListPurchaseOrderStatuses ----------------------- */

// ListPurchaseOrderStatusesResponse represents the response body of the ListPurchaseOrderStatuses endpoint.
type ListPurchaseOrderStatusesResponse struct {
	TotalCount    int                                                       `json:"total_count"`
	FilteredCount int                                                       `json:"filtered_count"`
	Data          []PurchaseOrderStatusResponseForListPurchaseOrderStatuses `json:"data"`
}

// PurchaseOrderStatusResponseForListPurchaseOrderStatuses represents an individual PurchaseOrderStatus returned from the ListPurchaseOrderStatuses endpoint.
type PurchaseOrderStatusResponseForListPurchaseOrderStatuses struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	StatusColor string `json:"status_color"`
	TextColor   string `json:"text_color"`
}

// GenerateListPurchaseOrderStatusesResponse  generates the response body for the ListPurchaseOrderStatuses endpoint.
func GenerateListPurchaseOrderStatusesResponse(purchaseOrderStatuses []models.PurchaseOrderStatus, totalCount int, filteredCount int) *ListPurchaseOrderStatusesResponse {

	var purchaseOrderStatusResponses []PurchaseOrderStatusResponseForListPurchaseOrderStatuses
	for _, purchaseOrderStatus := range purchaseOrderStatuses {
		purchaseOrderStatusResponses = append(purchaseOrderStatusResponses, PurchaseOrderStatusResponseForListPurchaseOrderStatuses{
			ID:          purchaseOrderStatus.ID,
			Name:        purchaseOrderStatus.Name,
			StatusColor: purchaseOrderStatus.StatusColor,
			TextColor:   purchaseOrderStatus.TextColor,
		})
	}

	return &ListPurchaseOrderStatusesResponse{
		TotalCount:    totalCount,
		FilteredCount: filteredCount,
		Data:          purchaseOrderStatusResponses,
	}

}

/* -------------------- CreatePurchaseOrderStatus ------------------- */

// CreatePurchaseOrderStatusResponse represents the response body of the CreatePurchaseOrderStatus endpoint.
type CreatePurchaseOrderStatusResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	StatusColor string `json:"status_color"`
	TextColor   string `json:"text_color"`
}

// GenerateCreatePurchaseOrderStatusResponse generates the response body for the CreatePurchaseOrderStatus endpoint.
func GenerateCreatePurchaseOrderStatusResponse(purchaseOrderStatus models.PurchaseOrderStatus) *CreatePurchaseOrderStatusResponse {

	return &CreatePurchaseOrderStatusResponse{
		ID:          purchaseOrderStatus.ID,
		Name:        purchaseOrderStatus.Name,
		StatusColor: purchaseOrderStatus.StatusColor,
		TextColor:   purchaseOrderStatus.TextColor,
	}

}

/* -------------------- UpdatePurchaseOrderStatus ------------------- */

// UpdatePurchaseOrderStatusResponse represents the response body of the UpdatePurchaseOrderStatus endpoint.
type UpdatePurchaseOrderStatusResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	StatusColor string `json:"status_color"`
	TextColor   string `json:"text_color"`
}

// GenerateUpdatePurchaseOrderStatusResponse generates the response body for the UpdatePurchaseOrderStatus endpoint.
func GenerateUpdatePurchaseOrderStatusResponse(purchaseOrderStatus models.PurchaseOrderStatus) *UpdatePurchaseOrderStatusResponse {

	return &UpdatePurchaseOrderStatusResponse{
		ID:          purchaseOrderStatus.ID,
		Name:        purchaseOrderStatus.Name,
		StatusColor: purchaseOrderStatus.StatusColor,
		TextColor:   purchaseOrderStatus.TextColor,
	}

}
