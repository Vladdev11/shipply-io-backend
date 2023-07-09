package tasks

import (
	"context"

	"github.com/shipply-io/shipply-io-backend/models"

	shipengineHandlers "github.com/shipply-io/shipply-io-backend/api/shipengine/handlers"
)

func SyncCarrierConnectionOptions(ctx context.Context) func() {
	return func() {
		//get all carrier connections
		carrierConnections, err := models.GetAllCarrierConnections(ctx)
		if err != nil {
			return
		}

		//for each carrier connection
		for _, carrierConnection := range carrierConnections {

			//get carrier connection options
			shipengineCarrierConnectionOptions, err := shipengineHandlers.GetCarrierConnectionOptions(ctx, carrierConnection.ShipengineCarrierID)
			if err != nil {
				continue
			}

			//convert shipengine carrier options to shippi carrier options
			carrierConnectionOptions := shipengineCarrierConnectionOptions.ConvertToCarrierConnectionOptions()

			//save carrier connection options
			err = carrierConnection.SaveCarrierOptions(ctx, carrierConnectionOptions)
			if err != nil {
				continue
			}

		}
	}
}

func SyncCarrierConnectionPackageTypes(ctx context.Context) func() {
	return func() {
		//get all carrier connections
		carrierConnections, err := models.GetAllCarrierConnections(ctx)
		if err != nil {
			return
		}

		//for each carrier connection
		for _, carrierConnection := range carrierConnections {

			//get carrier package types
			shipengineCarrierPackageTypes, err := shipengineHandlers.GetCarrierConnectionPackageTypes(ctx, carrierConnection.ShipengineCarrierID)
			if err != nil {
				continue
			}

			//convert shipengine carrier package types to shippi carrier package types
			carrierConnectionPackageTypes := shipengineCarrierPackageTypes.ConvertToCarrierConnectionPackageTypes()

			//save carrier connection package types
			err = carrierConnection.SaveCarrierPackageTypes(ctx, carrierConnectionPackageTypes)
			if err != nil {
				continue
			}

		}
	}

}

func SyncCarrierConnectionServices(ctx context.Context) func() {
	return func() {
		//get all carrier connections
		carrierConnections, err := models.GetAllCarrierConnections(ctx)
		if err != nil {
			return
		}

		//for each carrier connection
		for _, carrierConnection := range carrierConnections {
			err = SyncCarrierConnectionServicesByCarrierConnection(ctx, carrierConnection)
			if err != nil {
				continue
			}
		}
	}
}

func SyncCarrierConnectionServicesByCarrierConnection(ctx context.Context, cc models.CarrierConnection) error {
	//get carrier services
	shipengineCarrierServices, err := shipengineHandlers.GetCarrierConnectionServices(ctx, cc.ShipengineCarrierID)
	if err != nil {
		return err
	}

	//convert shipengine carrier services to shippi carrier services
	carrierConnectionServices := shipengineCarrierServices.ConvertToCarrierConnectionServices()

	//save carrier connection services
	err = cc.SaveCarrierServices(ctx, carrierConnectionServices)
	if err != nil {
		return err
	}

	return nil
}
