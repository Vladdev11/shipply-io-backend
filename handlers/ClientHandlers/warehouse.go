package ClientHandlers

import (
	"net/http"

	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/responses"
	"github.com/shipply-io/shipply-io-backend/util"
)

func ListWarehouses(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetClient(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get client", http.StatusUnauthorized)
		return
	}

	err = user.Client.GetOrganization(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusUnauthorized)
		return
	}

	request := models.WarehouseListRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	request.OrganizationID = user.Client.Organization.ID
	warehouses, total, count, err := user.Client.Organization.GetWarehouses(ctx, request)

	response := responses.GenerateListWarehousesResponse(warehouses, total, count)
	util.JSONResponse(w, response, http.StatusOK)

}

func GetWarehouse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetClient(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get client", http.StatusUnauthorized)
		return
	}

	warehouseID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "failed to get warehouse id", http.StatusBadRequest)
		return
	}

	warehouse, err := models.GetWarehouseByID(ctx, warehouseID)
	if err != nil {
		util.ErrorResponse(w, "failed to get warehouse", http.StatusBadRequest)
		return
	}

	err = warehouse.GetShipFromAddress(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get ship from address", http.StatusBadRequest)
		return
	}

	err = warehouse.GetReturnAddress(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get return address", http.StatusBadRequest)
		return
	}

	if user.Client.OrganizationID != warehouse.OrganizationID {
		util.ErrorResponse(w, "user does not have access to warehouse", http.StatusForbidden)
		return
	}

	response := responses.GenerateGetWarehouseResponse(*warehouse)
	util.JSONResponse(w, response, http.StatusOK)
}
