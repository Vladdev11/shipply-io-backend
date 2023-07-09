package OrganizationHandlers

import (
	"net/http"

	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/util"
)

func LocationList(w http.ResponseWriter, r *http.Request) {
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

	request := models.LocationListRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	locations, total, count, err := user.Organization.GetLocations(ctx, request)
	if err != nil {
		util.ErrorResponse(w, "failed to get locations", http.StatusBadRequest)
		return
	}

	searchResults, err := models.ConvertLocationsToSearchResults(ctx, locations, total, count)
	if err != nil {
		util.ErrorResponse(w, "failed to convert locations to search results", http.StatusBadRequest)
		return
	}

	util.JSONResponse(w, searchResults, http.StatusOK)

}

func LocationCreate(w http.ResponseWriter, r *http.Request) {
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

	request := models.LocationCreateRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	location := models.Location{
		Name:           request.Name,
		WarehouseID:    request.WarehouseID,
		LocationTypeID: request.LocationTypeID,
		Pickable:       request.Pickable,
		Sellable:       request.Sellable,
		IsTote:         request.IsTote,
	}

	err = location.Create(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to create location", http.StatusBadRequest)
		return
	}

	err = location.GetLocationType(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get location type", http.StatusBadRequest)
		return
	}

	util.JSONResponse(w, location.ConvertToReturnJSON(ctx), http.StatusOK)

}

func LocationGet(w http.ResponseWriter, r *http.Request) {
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

	locationID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "Invalid location id", http.StatusBadRequest)
		return
	}

	location, err := models.GetLocationByID(ctx, locationID)
	if err != nil {
		util.ErrorResponse(w, "failed to get location", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsWarehouseOwner(ctx, location.WarehouseID) {
		util.ErrorResponse(w, "user does not have access to location", http.StatusForbidden)
		return
	}

	err = location.GetLocationType(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get location type", http.StatusBadRequest)
		return
	}

	util.JSONResponse(w, location.ConvertToReturnJSON(ctx), http.StatusOK)

}

func LocationUpdate(w http.ResponseWriter, r *http.Request) {
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

	locationID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "Invalid location id", http.StatusBadRequest)
		return
	}

	location, err := models.GetLocationByID(ctx, locationID)
	if err != nil {
		util.ErrorResponse(w, "failed to get location", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsWarehouseOwner(ctx, location.WarehouseID) {
		util.ErrorResponse(w, "user does not have access to this location", http.StatusForbidden)
		return
	}

	request := &models.LocationUpdateRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	err = location.UpdateWithRequest(ctx, request)
	if err != nil {
		util.ErrorResponse(w, "failed to update location", http.StatusBadRequest)
		return
	}

	err = location.GetLocationType(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get location type", http.StatusBadRequest)
		return
	}

	util.JSONResponse(w, location.ConvertToReturnJSON(ctx), http.StatusOK)

}

func LocationDelete(w http.ResponseWriter, r *http.Request) {
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

	locationID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "Invalid location id", http.StatusBadRequest)
		return
	}

	location, err := models.GetLocationByID(ctx, locationID)
	if err != nil {
		util.ErrorResponse(w, "failed to get location", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsWarehouseOwner(ctx, location.WarehouseID) {
		util.ErrorResponse(w, "user does not have access to this location", http.StatusForbidden)
		return
	}

	if location.HasInventory(ctx) {
		util.ErrorResponse(w, "location has inventory", http.StatusBadRequest)
		return
	}

	if location.HasActivePickSessionOrder(ctx) {
		util.ErrorResponse(w, "location has active pick session order", http.StatusBadRequest)
		return
	}

	err = location.Delete(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to delete location", http.StatusBadRequest)
		return
	}

	util.SuccessResponse(w, http.StatusOK)
}
