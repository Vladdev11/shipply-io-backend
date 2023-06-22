package ShipengineHandlers

import (
	"fmt"

	ShipengineAPI "github.com/shipply-io/shipply-io-backend/api/shipengine/api"
	ShipengineModels "github.com/shipply-io/shipply-io-backend/api/shipengine/models"
)

func CreateShipment(shipment ShipengineModels.CreateShipmentRequest) (*ShipengineModels.CreateShipmentResponse, error) {

	shipmentResponse, errorResponse, err := ShipengineAPI.CreateShipment(shipment)
	if err != nil {
		return nil, err
	}

	if errorResponse != nil {
		return nil, fmt.Errorf("error creating shipment: %v", errorResponse.Errors[0].Message)
	}

	return shipmentResponse, nil
}
