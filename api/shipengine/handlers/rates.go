package ShipengineHandlers

import (
	"fmt"

	ShipengineAPI "github.com/shipply-io/shipply-io-backend/api/shipengine/api"
	ShipengineModels "github.com/shipply-io/shipply-io-backend/api/shipengine/models"
)

func ShopRates(rsr ShipengineModels.RateShopRequest) (*ShipengineModels.RateShopResponse, error) {

	ratesShopResponse, errorResponse, err := ShipengineAPI.ShopRates(rsr)
	if err != nil {
		return nil, err
	}

	if errorResponse != nil {
		return nil, fmt.Errorf("error rate shopping: %v", errorResponse.Errors[0].Message)
	}

	return ratesShopResponse, nil
}
