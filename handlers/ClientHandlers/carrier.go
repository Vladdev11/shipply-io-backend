package ClientHandlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/responses"
	"github.com/shipply-io/shipply-io-backend/tasks"
	"github.com/shipply-io/shipply-io-backend/util"

	shipengineHandlers "github.com/shipply-io/shipply-io-backend/api/shipengine/handlers"
)

func ListCarriers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetClient(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusUnauthorized)
		return
	}

	carriers, err := models.GetCarriers(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusInternalServerError)
		return
	}

	carriers = models.GetCarrierRequiredFields(carriers)

	response := responses.GenerateListCarriersResponse(carriers)
	util.JSONResponse(w, response, http.StatusOK)

}

func CreateCarrierConnection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetClient(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusUnauthorized)
		return
	}

	carrierID, err := util.GetIntFromPath(r, "carrier_id")
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	carrier, err := models.GetCarrierByID(ctx, carrierID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	var carrierConnect models.CarrierConnect
	createCarrierConnection, errs := models.CreateCarrierConnection(&carrier, r)
	if errs != nil {
		util.ErrResponse(w, errs, http.StatusBadRequest)
		return
	}
	carrierConnect = *createCarrierConnection

	shipEngineResponse, err := shipengineHandlers.ConnectCarrier(ctx, carrierConnect, carrier.ShipEngineID)
	if err != nil {
		util.ErrResponse(w, ErrShipEngineConnectCarrier, http.StatusBadRequest)
		return
	}

	carrierNickname, err := carrierConnect.GetNickname()
	if err != nil {
		util.ErrResponse(w, err, http.StatusInternalServerError)
		return
	}

	carrierConnectJSON, err := json.Marshal(carrierConnect)
	if err != nil {
		util.ErrResponse(w, ErrMarshalJSON, http.StatusInternalServerError)
	}

	carrierConnection := models.CarrierConnection{
		CarrierID:           carrier.ID,
		OwnerID:             user.Client.ID,
		OwnerType:           2,
		ShipengineCarrierID: shipEngineResponse.CarrierID,
		ShipengineNickname:  carrierNickname,
		Active:              true,
		Credentials:         carrierConnectJSON,
	}

	err = carrierConnection.Create(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusInternalServerError)
		return
	}

	//get carrier connection options
	shipengineCarrierConnectionOptions, err := shipengineHandlers.GetCarrierConnectionOptions(ctx, carrierConnection.ShipengineCarrierID)
	if err != nil {
		models.CreateSystemError(ctx, fmt.Sprintf("%e", err))
	}

	//convert shipengine carrier options to shippi carrier options
	carrierConnectionOptions := shipengineCarrierConnectionOptions.ConvertToCarrierConnectionOptions()

	//save carrier connection options
	err = carrierConnection.SaveCarrierOptions(ctx, carrierConnectionOptions)
	if err != nil {
		models.CreateSystemError(ctx, fmt.Sprintf("%e", err))
	}

	//sync carrier connection services
	tasks.SyncCarrierConnectionServicesByCarrierConnection(ctx, carrierConnection)

	response := responses.GenerateCreateCarrierConnectionResponse(carrierConnection)
	util.JSONResponse(w, response, http.StatusOK)

}

func DisconnectCarrierConnection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetClient(ctx)
	if err != nil {
		util.ErrResponse(w, ErrGetClient, http.StatusUnauthorized)
		return
	}

	carrierID, err := util.GetIntFromPath(r, "carrier_connection_id")
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	carrierConnection, err := models.GetCarrierConnectionByID(ctx, carrierID)
	if err != nil {
		util.ErrResponse(w, ErrGetCarrier, http.StatusBadRequest)
		return
	}

	//if carrierconnection is owned by the client, make sure it belongs to the organization
	if carrierConnection.OwnerType == 1 {
		util.ErrResponse(w, ErrCarrierConnectionBelongsToOrganization, http.StatusUnauthorized)
		return
	}

	if carrierConnection.OwnerID != user.OwnerID {
		util.ErrResponse(w, ErrCarrierConnectionDoesNotBelongToClient, http.StatusUnauthorized)
		return
	}

	carrier, err := models.GetCarrierByID(ctx, carrierConnection.CarrierID)
	if err != nil {
		util.ErrResponse(w, ErrGetCarrier, http.StatusBadRequest)
		return
	}

	//TODO make sure there are no shipments waiting to use this carrier for manifests

	err = shipengineHandlers.DeleteCarrier(ctx, carrier.ShipEngineID, carrierConnection.ShipengineCarrierID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	err = carrierConnection.Delete(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusInternalServerError)
		return
	}

	util.SuccessResponse(w, http.StatusOK)

}

func ListCarrierConnections(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetClient(ctx)
	if err != nil {
		util.ErrResponse(w, ErrGetClient, http.StatusUnauthorized)
		return
	}

	carrierConnections, err := user.Client.GetCarrierConnections(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusInternalServerError)
		return
	}

	response := responses.GenerateListCarrierConnectionsResponse(carrierConnections)
	util.JSONResponse(w, response, http.StatusOK)

}

func GetCarrierConnection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetClient(ctx)
	if err != nil {
		util.ErrResponse(w, ErrGetClient, http.StatusUnauthorized)
		return
	}

	carrierConnectionID, err := util.GetIntFromPath(r, "carrier_connection_id")
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	carrierConnection, err := models.GetCarrierConnectionByID(ctx, carrierConnectionID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if carrierConnection.OwnerID != user.OwnerID {
		util.ErrResponse(w, ErrCarrierConnectionDoesNotBelongToClient, http.StatusUnauthorized)
		return
	}

	response := responses.GenerateGetCarrierConnectionResponse(*carrierConnection)
	util.JSONResponse(w, response, http.StatusOK)

}
