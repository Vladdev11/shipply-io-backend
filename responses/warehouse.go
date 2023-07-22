package responses

import (
	"time"

	"github.com/shipply-io/shipply-io-backend/models"
)

/* ----------------------------- ListWarehouses ----------------------------- */

// ListWarehousesResponse represents the expected response body for the ListWarehouses endpoint
type ListWarehousesResponse struct {
	TotalCount    int                                  `json:"total_count"`
	FilteredCount int                                  `json:"filtered_count"`
	Data          []WarehouseResponseForListWarehouses `json:"data"`
}

// WarehouseResponseForListWarehouses represents the expected response body for the ListWarehouses endpoint
type WarehouseResponseForListWarehouses struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GenerateListWarehousesResponse generates the response body for the ListWarehouses endpoint
func GenerateListWarehousesResponse(warehouses []models.Warehouse, totalCount int, filteredCount int) *ListWarehousesResponse {

	var warehouseResponses []WarehouseResponseForListWarehouses
	for _, warehouse := range warehouses {
		warehouseResponses = append(warehouseResponses, WarehouseResponseForListWarehouses{
			ID:        warehouse.ID,
			Name:      warehouse.Name,
			UpdatedAt: warehouse.UpdatedAt,
		})
	}

	return &ListWarehousesResponse{
		TotalCount:    totalCount,
		FilteredCount: filteredCount,
		Data:          warehouseResponses,
	}

}

/* ------------------------------ GetWarehouse ------------------------------ */

// GetWarehouseResponse represents the expected response body for the GetWarehouse endpoint
type GetWarehouseResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`

	ShipFromAddress *Address `json:"ship_from_address"`
	ReturnAddress   *Address `json:"return_address"`
}

// GenerateGetWarehouseResponse generates the response body for the GetWarehouse endpoint
func GenerateGetWarehouseResponse(warehouse models.Warehouse) *GetWarehouseResponse {

	response := &GetWarehouseResponse{
		ID:        warehouse.ID,
		Name:      warehouse.Name,
		UpdatedAt: warehouse.UpdatedAt,
	}

	if &warehouse.ShipFromAddressID != nil {
		response.ShipFromAddress = &Address{
			ID:         warehouse.ShipFromAddress.ID,
			FirstName:  warehouse.ShipFromAddress.FirstName,
			LastName:   warehouse.ShipFromAddress.LastName,
			Company:    warehouse.ShipFromAddress.Company,
			Street1:    warehouse.ShipFromAddress.Street1,
			Street2:    warehouse.ShipFromAddress.Street2,
			Street3:    warehouse.ShipFromAddress.Street3,
			City:       warehouse.ShipFromAddress.City,
			State:      warehouse.ShipFromAddress.State,
			PostalCode: warehouse.ShipFromAddress.PostalCode,
			Country:    warehouse.ShipFromAddress.Country,
			Phone:      warehouse.ShipFromAddress.Phone,
		}
	}

	if &warehouse.ReturnAddress != nil {
		response.ReturnAddress = &Address{
			ID:         warehouse.ReturnAddress.ID,
			FirstName:  warehouse.ReturnAddress.FirstName,
			LastName:   warehouse.ReturnAddress.LastName,
			Company:    warehouse.ReturnAddress.Company,
			Street1:    warehouse.ReturnAddress.Street1,
			Street2:    warehouse.ReturnAddress.Street2,
			Street3:    warehouse.ReturnAddress.Street3,
			City:       warehouse.ReturnAddress.City,
			State:      warehouse.ReturnAddress.State,
			PostalCode: warehouse.ReturnAddress.PostalCode,
			Country:    warehouse.ReturnAddress.Country,
			Phone:      warehouse.ReturnAddress.Phone,
		}
	}

	return response

}

/* ------------------------------ CreateWarehouse ------------------------------ */

// CreateWarehouseResponse represents the expected response body for the CreateWarehouse endpoint
type CreateWarehouseResponse struct {
	ID int `json:"id"`
}

// GenerateCreateWarehouseResponse generates the response body for the CreateWarehouse endpoint
func GenerateCreateWarehouseResponse(warehouse models.Warehouse) *CreateWarehouseResponse {

	return &CreateWarehouseResponse{
		ID: warehouse.ID,
	}

}

/* ------------------------------ UpdateWarehouse ------------------------------ */
type UpdateWarehouseResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`

	ShipFromAddress *Address `json:"ship_from_address"`
	ReturnAddress   *Address `json:"return_address"`
}

func GenerateUpdateWarehouseResponse(warehouse models.Warehouse) *UpdateWarehouseResponse {

	response := &UpdateWarehouseResponse{
		ID:        warehouse.ID,
		Name:      warehouse.Name,
		UpdatedAt: warehouse.UpdatedAt,
	}

	if &warehouse.ShipFromAddressID != nil {
		response.ShipFromAddress = &Address{
			ID:         warehouse.ShipFromAddress.ID,
			FirstName:  warehouse.ShipFromAddress.FirstName,
			LastName:   warehouse.ShipFromAddress.LastName,
			Company:    warehouse.ShipFromAddress.Company,
			Street1:    warehouse.ShipFromAddress.Street1,
			Street2:    warehouse.ShipFromAddress.Street2,
			Street3:    warehouse.ShipFromAddress.Street3,
			City:       warehouse.ShipFromAddress.City,
			State:      warehouse.ShipFromAddress.State,
			PostalCode: warehouse.ShipFromAddress.PostalCode,
			Country:    warehouse.ShipFromAddress.Country,
			Phone:      warehouse.ShipFromAddress.Phone,
		}
	}

	if &warehouse.ReturnAddress != nil {
		response.ReturnAddress = &Address{
			ID:         warehouse.ReturnAddress.ID,
			FirstName:  warehouse.ReturnAddress.FirstName,
			LastName:   warehouse.ReturnAddress.LastName,
			Company:    warehouse.ReturnAddress.Company,
			Street1:    warehouse.ReturnAddress.Street1,
			Street2:    warehouse.ReturnAddress.Street2,
			Street3:    warehouse.ReturnAddress.Street3,
			City:       warehouse.ReturnAddress.City,
			State:      warehouse.ReturnAddress.State,
			PostalCode: warehouse.ReturnAddress.PostalCode,
			Country:    warehouse.ReturnAddress.Country,
			Phone:      warehouse.ReturnAddress.Phone,
		}
	}

	return response

}
