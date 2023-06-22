package tasks

import (
	"github.com/shipply-io/shipply-io-backend/models"

	shipengineHandlers "github.com/shipply-io/shipply-io-backend/api/shipengine/handlers"
)

func SyncCarrierConnectionOptions() {

	//get all carrier connections
	carrierConnections, err := models.GetAllCarrierConnections()
	if err != nil {
		return
	}

	//for each carrier connection
	for _, carrierConnection := range carrierConnections {

		//get carrier connection options
		shipengineCarrierConnectionOptions, err := shipengineHandlers.GetCarrierConnectionOptions(carrierConnection.ShipengineCarrierID)
		if err != nil {
			continue
		}

		//convert shipengine carrier options to shippi carrier options
		carrierConnectionOptions := shipengineCarrierConnectionOptions.ConvertToCarrierConnectionOptions()

		//save carrier connection options
		err = carrierConnection.SaveCarrierOptions(carrierConnectionOptions)
		if err != nil {
			continue
		}

	}
}

func SyncCarrierConnectionPackageTypes() {

	//get all carrier connections
	carrierConnections, err := models.GetAllCarrierConnections()
	if err != nil {
		return
	}

	//for each carrier connection
	for _, carrierConnection := range carrierConnections {

		//get carrier package types
		shipengineCarrierPackageTypes, err := shipengineHandlers.GetCarrierConnectionPackageTypes(carrierConnection.ShipengineCarrierID)
		if err != nil {
			continue
		}

		//convert shipengine carrier package types to shippi carrier package types
		carrierConnectionPackageTypes := shipengineCarrierPackageTypes.ConvertToCarrierConnectionPackageTypes()

		//save carrier connection package types
		err = carrierConnection.SaveCarrierPackageTypes(carrierConnectionPackageTypes)
		if err != nil {
			continue
		}

	}

}

func SyncCarrierConnectionServices() {

	//get all carrier connections
	carrierConnections, err := models.GetAllCarrierConnections()
	if err != nil {
		return
	}

	//for each carrier connection
	for _, carrierConnection := range carrierConnections {
		err = SyncCarrierConnectionServicesByCarrierConnection(carrierConnection)
		if err != nil {
			continue
		}
	}

}

func SyncCarrierConnectionServicesByCarrierConnection(cc models.CarrierConnection) error {
	//get carrier services
	shipengineCarrierServices, err := shipengineHandlers.GetCarrierConnectionServices(cc.ShipengineCarrierID)
	if err != nil {
		return err
	}

	//convert shipengine carrier services to shippi carrier services
	carrierConnectionServices := shipengineCarrierServices.ConvertToCarrierConnectionServices()

	//save carrier connection services
	err = cc.SaveCarrierServices(carrierConnectionServices)
	if err != nil {
		return err
	}

	return nil
}
