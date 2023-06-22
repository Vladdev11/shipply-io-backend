package ShipengineHandlers

import (
	"fmt"

	ShipengineAPI "github.com/shipply-io/shipply-io-backend/api/shipengine/api"
	ShipengineModels "github.com/shipply-io/shipply-io-backend/api/shipengine/models"
)

func PurchaseLabelFromRate(plfrr ShipengineModels.PurchaseLabelFromRateRequest, RateID string) (*ShipengineModels.PurchaseLabelFromRateResponse, error) {

	purchaseLabelFromRateResponse, errorResponse, err := ShipengineAPI.PurchaseLabelFromRate(plfrr, RateID)
	if err != nil {
		return nil, err
	}

	if errorResponse != nil {
		return nil, fmt.Errorf("error purchasing label: %v", errorResponse.Errors[0].Message)
	}

	return purchaseLabelFromRateResponse, nil
}
