package OrganizationHandlers

import (
	"fmt"
	"net/http"

	ShipengineHandlers "github.com/shipply-io/shipply-io-backend/api/shipengine/handlers"
	ShipengineModels "github.com/shipply-io/shipply-io-backend/api/shipengine/models"
	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/util"
)

func WarehouseList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetOrganization(ctx)
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

	request.OrganizationID = user.Organization.ID
	warehouses, total, count, err := user.Organization.GetWarehouses(ctx, request)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	searchResults, err := models.ConvertWarehousesToSearchResults(warehouses, total, count)
	if err != nil {
		util.ErrorResponse(w, "failed to convert warehouses to search results", http.StatusBadRequest)
		return
	}

	util.JSONResponse(w, searchResults, http.StatusOK)

}

func WarehouseCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetOrganization(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusUnauthorized)
		return
	}

	request := models.WarehouseCreateRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	err = request.ShipFromAddress.Create(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to create ship from address", http.StatusBadRequest)
		return
	}

	warehouse := &models.Warehouse{
		OrganizationID:    user.Organization.ID,
		Name:              request.Name,
		ShipFromAddressID: request.ShipFromAddress.ID,
		ShipFromAddress:   *request.ShipFromAddress,
		ReturnAddressID:   request.ReturnAddress.ID,
		ReturnAddress:     *request.ReturnAddress,
	}

	shipengineWarehouseCreateRequest := ShipengineModels.ConvertWarehouseToShipengineWarehouse(warehouse)
	shipengineWarehouse, err := ShipengineHandlers.CreateWarehouse(ctx, shipengineWarehouseCreateRequest)
	if err != nil {
		fmt.Println(err)
		util.ErrorResponse(w, "error registering warehouse", http.StatusBadRequest)
		return
	}

	warehouse.ShipengineWarehouseID = shipengineWarehouse.WarehouseID

	err = warehouse.Create(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to create warehouse", http.StatusBadRequest)
		return
	}

	warehouseJson := warehouse.ConvertToReturnJSON()

	util.JSONResponse(w, warehouseJson, http.StatusOK)
}

func WarehouseGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
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

	if user.OwnerID != warehouse.OrganizationID {
		util.ErrorResponse(w, "user does not have access to warehouse", http.StatusForbidden)
		return
	}

	warehouseJson := warehouse.ConvertToReturnJSON()

	util.JSONResponse(w, warehouseJson, http.StatusOK)
}

func WarehouseUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	warehouseID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "invalid warehouse id", http.StatusBadRequest)
		return
	}

	warehouse, err := models.GetWarehouseByID(ctx, warehouseID)
	if err != nil {
		util.ErrorResponse(w, "failed to find warehouse", http.StatusBadRequest)
		return
	}

	if user.OwnerID != warehouse.OrganizationID {
		util.ErrorResponse(w, "user does not have access to this warehouse", http.StatusForbidden)
		return
	}

	request := &models.WarehouseUpdateRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	err = warehouse.UpdateWithRequest(ctx, request)
	if err != nil {
		util.ErrorResponse(w, "failed to update warehouse", http.StatusInternalServerError)
		return
	}

	err = warehouse.GetShipFromAddress(ctx)
	if err != nil {
		fmt.Print(err)
		util.ErrorResponse(w, "failed to get ship from address", http.StatusInternalServerError)
		return
	}

	err = warehouse.GetReturnAddress(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get return address", http.StatusBadRequest)
		return
	}

	warehouseJSON := warehouse.ConvertToReturnJSON()
	util.JSONResponse(w, warehouseJSON, http.StatusOK)
}

func WarehouseDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	warehouseID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "invalid warehouse id", http.StatusBadRequest)
		return
	}

	warehouse, err := models.GetWarehouseByID(ctx, warehouseID)
	if err != nil {
		util.ErrorResponse(w, "failed to find warehouse", http.StatusBadRequest)
		return
	}

	if user.OwnerID != warehouse.OrganizationID {
		util.ErrorResponse(w, "user does not have access to this warehouse", http.StatusForbidden)
		return
	}

	locations, err := models.GetLocationsByWarehouseID(ctx, warehouse.ID)
	if err != nil {
		util.ErrorResponse(w, "failed to get locations", http.StatusInternalServerError)
		return
	}

	if len(locations) > 0 {
		util.ErrorResponse(w, "cannot delete location because it is associated with one or more locations", http.StatusBadRequest)
		return
	}

	err = warehouse.Delete(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to delete warehouse", http.StatusInternalServerError)
		return
	}

	util.SuccessResponse(w, http.StatusOK)
}
