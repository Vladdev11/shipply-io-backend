package ShipengineHandlers

import (
	"fmt"

	ShipengineAPI "github.com/shipply-io/shipply-io-backend/api/shipengine/api"
	ShipengineModels "github.com/shipply-io/shipply-io-backend/api/shipengine/models"
)

func ValidateAddress(inputAddress ShipengineModels.Address) (*ShipengineModels.AddressResponse, error) {

	addressRequest := []ShipengineModels.Address{
		inputAddress,
	}

	addressResponse, errorResponse, err := ShipengineAPI.ValidateAddress(addressRequest)
	if err != nil {
		return nil, err
	}

	if errorResponse != nil {
		return nil, fmt.Errorf("error validating address: %v", errorResponse.Errors[0].Message)
	}

	address := addressResponse[0]

	return address, nil

}

func ParseAddress(addressText string) (ShipengineModels.Address, error) {

	addressParseRequest := ShipengineModels.AddressParseRequest{
		Text: addressText,
	}

	addressParseResponse, errorResponse, err := ShipengineAPI.ParseAddress(addressParseRequest)
	if err != nil {
		return ShipengineModels.Address{}, err
	}

	if errorResponse != nil {
		return ShipengineModels.Address{}, fmt.Errorf("error parsing address: %v", errorResponse.Errors[0].Message)
	}

	//TODO logic around when we use / don't use the parsed address based on score

	return addressParseResponse.Address, nil

}
