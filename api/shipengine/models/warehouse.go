package ShipengineModels

import (
	"time"

	"github.com/shipply-io/shipply-io-backend/models"
)

type Warehouse struct {
	Name          string  `json:"name"`
	OriginAddress Address `json:"origin_address"`
	ReturnAddress Address `json:"return_address"`
}

type WarehouseResponse struct {
	WarehouseID   string    `json:"warehouse_id"`
	IsDefault     bool      `json:"is_default"`
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"created_at"`
	OriginAddress Address   `json:"origin_address"`
	ReturnAddress Address   `json:"return_address"`
}

func ConvertWarehouseToShipengineWarehouse(warehouse *models.Warehouse) *Warehouse {
	return &Warehouse{
		Name:          warehouse.Name,
		OriginAddress: ConvertAddressToShipengineAddress(&warehouse.ShipFromAddress),
		ReturnAddress: ConvertAddressToShipengineAddress(&warehouse.ShipFromAddress),
	}
}
