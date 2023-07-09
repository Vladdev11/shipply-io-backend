package ShipengineHandlers

import (
	"context"
	"fmt"

	ShipengineAPI "github.com/shipply-io/shipply-io-backend/api/shipengine/api"
	ShipengineModels "github.com/shipply-io/shipply-io-backend/api/shipengine/models"
)

func PurchaseLabelFromRate(ctx context.Context, plfrr ShipengineModels.PurchaseLabelFromRateRequest, RateID string) (*ShipengineModels.PurchaseLabelFromRateResponse, error) {

	purchaseLabelFromRateResponse, errorResponse, err := ShipengineAPI.PurchaseLabelFromRate(ctx, plfrr, RateID)
	if err != nil {
		return nil, err
	}

	if errorResponse != nil {
		return nil, fmt.Errorf("error purchasing label: %v", errorResponse.Errors[0].Message)
	}

	return purchaseLabelFromRateResponse, nil
}
