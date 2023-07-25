package ShipengineHandlers

import (
	"context"
	"fmt"

	ShipengineAPI "github.com/shipply-io/shipply-io-backend/api/shipengine/api"
	ShipengineModels "github.com/shipply-io/shipply-io-backend/api/shipengine/models"
)

func ConnectCarrier(ctx context.Context, carrier interface{}, carrierName string) (*ShipengineModels.CarrierConnectResponse, error) {

	if carrier == nil {
		return nil, ErrEmptyCarrierInformation
	}

	carrierConnect := ShipengineModels.CarrierConnect{
		Carrier: carrier,
	}

	connectCarrierResponse, errorResponse, err := ShipengineAPI.ConnectCarrier(ctx, carrierConnect, carrierName)

	if err != nil {
		return nil, err
	}

	if errorResponse != nil {
		if len(errorResponse.Errors) == 0 {
			return nil, fmt.Errorf("error connecting carrier: %v", errorResponse)
		} else {
			return nil, fmt.Errorf("error connecting carrier: %v", errorResponse.Errors[0].Message)
		}
	}

	return connectCarrierResponse, nil

}

func DeleteCarrier(ctx context.Context, carrierName string, carrierId string) error {

	_, errorResponse, err := ShipengineAPI.DeleteCarrier(ctx, carrierName, carrierId)

	if err != nil {
		return err
	}

	if errorResponse != nil {
		return fmt.Errorf("%e:%v", ErrDeleteCarrierFailed, errorResponse.Errors[0].Message)
	}

	return nil

}

func GetCarrierConnectionOptions(ctx context.Context, carrier_id string) (*ShipengineModels.CarrierOptionsResponse, error) {

	carrierOptionsResponse, errorResponse, err := ShipengineAPI.GetCarrierConnectionOptions(ctx, carrier_id)
	if err != nil {
		return nil, err
	}

	if errorResponse != nil {
		return nil, fmt.Errorf("error getting carrier options: %v", errorResponse.Errors[0].Message)
	}

	return carrierOptionsResponse, nil

}

func GetCarrierConnectionPackageTypes(ctx context.Context, carrier_id string) (*ShipengineModels.CarrierPackageTypesResponse, error) {

	carrierPackageTypesResponse, errorResponse, err := ShipengineAPI.GetCarrierConnectionPackageTypes(ctx, carrier_id)
	if err != nil {
		return nil, err
	}

	if errorResponse != nil {
		return nil, fmt.Errorf("error getting carrier package types: %v", errorResponse.Errors[0].Message)
	}

	return carrierPackageTypesResponse, nil

}

func GetCarrierConnectionServices(ctx context.Context, carrier_id string) (*ShipengineModels.CarrierServicesResponse, error) {

	carrierServicesResponse, errorResponse, err := ShipengineAPI.GetCarrierConnectionServices(ctx, carrier_id)
	if err != nil {
		return nil, err
	}

	if errorResponse != nil {
		return nil, fmt.Errorf("error getting carrier services: %v", errorResponse.Errors[0].Message)
	}

	return carrierServicesResponse, nil

}
