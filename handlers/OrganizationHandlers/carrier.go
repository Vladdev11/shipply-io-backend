package OrganizationHandlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/tasks"
	"github.com/shipply-io/shipply-io-backend/util"

	shipengineHandlers "github.com/shipply-io/shipply-io-backend/api/shipengine/handlers"
)

func ListCarriers(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetOrganization()
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusUnauthorized)
		return
	}

	carriers, err := models.GetCarriers()
	if err != nil {
		util.ErrorResponse(w, "failed to get carriers", http.StatusInternalServerError)
		return
	}

	carriers = models.GetCarrierRequiredFields(carriers)

	carrierJSON := []models.CarrierReturnJSON{}
	for _, carrier := range carriers {
		carrierJSON = append(carrierJSON, carrier.ConvertToReturnJSON())
	}

	util.JSONResponse(w, carrierJSON, http.StatusOK)

}

func CreateCarrierConnection(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetOrganization()
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusUnauthorized)
		return
	}

	carrierID, err := util.GetIntFromPath(r, "carrier_id")
	if err != nil {
		util.ErrorResponse(w, "failed to parse carrier id", http.StatusBadRequest)
		return
	}

	carrier, err := models.GetCarrierByID(carrierID)
	if err != nil {
		util.ErrorResponse(w, "failed to get carrier", http.StatusBadRequest)
		return
	}

	var carrierConnect models.CarrierConnect
	createCarrierConnection, errs := models.CreateCarrierConnection(&carrier, r)
	if len(errs) > 0 {
		util.ErrorsResponse(w, errs, http.StatusBadRequest)
		return
	}

	carrierConnect = *createCarrierConnection

	response, err := shipengineHandlers.ConnectCarrier(carrierConnect, carrier.ShipEngineID)
	if err != nil {
		util.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	carrierNickname, err := carrierConnect.GetNickname()
	if err != nil {
		util.ErrorResponse(w, "failed to get carrier nickname", http.StatusInternalServerError)
		return
	}

	carrierConnectJSON, err := json.Marshal(carrierConnect)
	if err != nil {
		util.ErrorResponse(w, "failed to marshal carrier connect json", http.StatusInternalServerError)
	}

	carrierConnection := models.CarrierConnection{
		CarrierID:           carrier.ID,
		OwnerID:             user.Organization.ID,
		OwnerType:           1,
		ShipengineCarrierID: response.CarrierID,
		ShipengineNickname:  carrierNickname,
		Active:              true,
		Credentials:         carrierConnectJSON,
	}

	err = carrierConnection.Create()
	if err != nil {
		util.ErrorResponse(w, "failed to add carrier connection to database", http.StatusInternalServerError)
		return
	}

	//get carrier connection options
	shipengineCarrierConnectionOptions, err := shipengineHandlers.GetCarrierConnectionOptions(carrierConnection.ShipengineCarrierID)
	if err != nil {
		models.CreateSystemError(fmt.Sprintf("failed to get carrier connection options(%s): %s", carrierConnection.ShipengineCarrierID, err.Error()))
	}

	//convert shipengine carrier options to shippi carrier options
	carrierConnectionOptions := shipengineCarrierConnectionOptions.ConvertToCarrierConnectionOptions()

	//save carrier connection options
	err = carrierConnection.SaveCarrierOptions(carrierConnectionOptions)
	if err != nil {
		models.CreateSystemError(fmt.Sprintf("failed to save carrier options(%s): %s", carrierConnection.ShipengineCarrierID, err.Error()))
	}
	//sync carrier connection services
	tasks.SyncCarrierConnectionServicesByCarrierConnection(carrierConnection)

	//get updated carrier connection with the services
	carrierConnectionUpdated, err := models.GetCarrierConnectionByID(carrierConnection.ID)
	if err != nil {
		util.ErrorResponse(w, "failed to get carrier connection", http.StatusInternalServerError)
		return
	}
	carrierConnection = *carrierConnectionUpdated

	//get carrier info
	carrierConnection.GetCarrier()

	//get carrier connection settings
	carrierConnection.GetSettings()

	carrierConnectionJSON := carrierConnection.ConvertToReturnJSON()

	util.JSONResponse(w, carrierConnectionJSON, http.StatusOK)

}

func GetCarrierConnection(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetOrganization()
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusUnauthorized)
		return
	}

	carrierConnectionId, err := util.GetIntFromPath(r, "carrier_connection_id")
	if err != nil {
		util.ErrorResponse(w, "failed to parse carrier connection id", http.StatusBadRequest)
		return
	}

	carrierConnection, err := models.GetCarrierConnectionByID(carrierConnectionId)
	if err != nil {
		util.ErrorResponse(w, "failed to get carrier connection", http.StatusBadRequest)
		return
	}

	//if carrierconnection is owned by the client, make sure it belongs to the organization
	if carrierConnection.OwnerType == 2 {
		err = user.Organization.GetClients()
		if err != nil {
			util.ErrorResponse(w, "failed to get clients", http.StatusInternalServerError)
			return
		}

		for _, client := range user.Organization.Clients {
			if client.ID == carrierConnection.OwnerID {
				break
			}
		}

		util.ErrorResponse(w, "carrier connection does not belong to organization", http.StatusUnauthorized)
		return
	} else {
		if carrierConnection.OwnerID != user.Organization.ID {
			util.ErrorResponse(w, "carrier connection does not belong to organization", http.StatusUnauthorized)
			return
		}
	}

	err = carrierConnection.GetSettings()

	carrierConnectionJSON := carrierConnection.ConvertToReturnJSON()

	util.JSONResponse(w, carrierConnectionJSON, http.StatusOK)

}

func DisconnectCarrierConnection(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetOrganization()
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusUnauthorized)
		return
	}

	carrierID, err := util.GetIntFromPath(r, "carrier_connection_id")
	if err != nil {
		util.ErrorResponse(w, "failed to parse carrier id", http.StatusBadRequest)
		return
	}

	carrierConnection, err := models.GetCarrierConnectionByID(carrierID)
	if err != nil {
		util.ErrorResponse(w, "failed to get carrier", http.StatusBadRequest)
		return
	}

	//if carrierconnection is owned by the client, make sure it belongs to the organization
	if carrierConnection.OwnerType == 2 {
		err = user.Organization.GetClients()
		if err != nil {
			util.ErrorResponse(w, "failed to get clients", http.StatusInternalServerError)
			return
		}

		for _, client := range user.Organization.Clients {
			if client.ID == carrierConnection.OwnerID {
				break
			}
		}

		util.ErrorResponse(w, "carrier connection does not belong to organization", http.StatusUnauthorized)
		return
	} else {
		if carrierConnection.OwnerID != user.Organization.ID {
			util.ErrorResponse(w, "carrier connection does not belong to organization", http.StatusUnauthorized)
			return
		}
	}

	carrier, err := models.GetCarrierByID(carrierConnection.CarrierID)
	if err != nil {
		util.ErrorResponse(w, "failed to get carrier", http.StatusBadRequest)
		return
	}

	//TODO make sure there are no shipments waiting to use this carrier for manifests

	err = shipengineHandlers.DeleteCarrier(carrier.ShipEngineID, carrierConnection.ShipengineCarrierID)
	if err != nil {
		util.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = carrierConnection.Delete()
	if err != nil {
		util.ErrorResponse(w, "failed to delete carrier connection from database", http.StatusInternalServerError)
		return
	}

	util.SuccessResponse(w, http.StatusOK)

}

func ListCarrierConnections(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetOrganization()
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusUnauthorized)
		return
	}

	carrierConnections, err := user.Organization.GetCarrierConnections()
	if err != nil {
		util.ErrorResponse(w, "failed to get carrier connections", http.StatusInternalServerError)
		return
	}

	carrierConnectionsJSON := []models.CarrierConnectionReturnJSON{}
	for _, carrierConnection := range carrierConnections {
		carrierConnection.GetCarrier()
		carrierConnectionsJSON = append(carrierConnectionsJSON, carrierConnection.ConvertToReturnJSON())
	}

	util.JSONResponse(w, carrierConnectionsJSON, http.StatusOK)

}
