package ShipengineHandlers

import (
	"context"
	"errors"

	ShipengineAPI "github.com/shipply-io/shipply-io-backend/api/shipengine/api"
	shipengineModels "github.com/shipply-io/shipply-io-backend/api/shipengine/models"
)

func CreateWarehouse(ctx context.Context, warehouseInput *shipengineModels.Warehouse) (*shipengineModels.WarehouseResponse, error) {

	if warehouseInput == nil {
		return nil, errors.New("warehouse information is empty")
	}

	warehouseCreateResponse, errorResponse, err := ShipengineAPI.CreateWarehouse(ctx, warehouseInput)
	if err != nil {
		return nil, err
	}

	if errorResponse != nil {
		return nil, errors.New(errorResponse.Errors[0].Message)
	}

	return warehouseCreateResponse, nil
}
