package ShipengineHandlers

import (
	"errors"

	ShipengineAPI "github.com/shipply-io/shipply-io-backend/api/shipengine/api"
	shipengineModels "github.com/shipply-io/shipply-io-backend/api/shipengine/models"
)

func CreateWarehouse(warehouseInput *shipengineModels.Warehouse) (*shipengineModels.WarehouseResponse, error) {

	if warehouseInput == nil {
		return nil, errors.New("warehouse information is empty")
	}

	warehouseCreateResponse, errorResponse, err := ShipengineAPI.CreateWarehouse(warehouseInput)
	if err != nil {
		return nil, err
	}

	if errorResponse != nil {
		return nil, errors.New(errorResponse.Errors[0].Message)
	}

	return warehouseCreateResponse, nil
}
