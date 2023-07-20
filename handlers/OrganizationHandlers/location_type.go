package OrganizationHandlers

import (
	"net/http"

	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/responses"
	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"
)

func ListLocationTypes(w http.ResponseWriter, r *http.Request) {
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

	// TODO change to typical list response with total count, filtered count, and data
	locationTypes, err := user.Organization.GetLocationTypes(ctx)
	if err != nil && err != gorm.ErrRecordNotFound {
		util.ErrorResponse(w, "failed to get location types", http.StatusBadRequest)
		return
	}

	response := responses.GenerateListLocationTypesResponse(locationTypes)
	util.JSONResponse(w, response, http.StatusOK)
}

func GetLocationType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	locationTypeID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "failed to get location type id", http.StatusBadRequest)
		return
	}

	locationType, err := models.GetLocationTypeByID(ctx, locationTypeID)
	if err != nil {
		util.ErrorResponse(w, "failed to get location type", http.StatusBadRequest)
		return
	}

	if locationType.OrganizationID != user.OwnerID {
		util.ErrorResponse(w, "user does not have access to this view location type", http.StatusForbidden)
		return
	}

	response := responses.GenerateGetLocationTypeResponse(*locationType)
	util.JSONResponse(w, response, http.StatusOK)
}

func CreateLocationType(w http.ResponseWriter, r *http.Request) {
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

	request := models.LocationTypeCreateRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	locationType := models.LocationType{
		Name:           request.Name,
		OrganizationID: user.Organization.ID,
	}

	err = locationType.Create(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to create location type", http.StatusBadRequest)
		return
	}

	response := responses.GenerateCreateLocationTypeResponse(locationType)
	util.JSONResponse(w, response, http.StatusOK)
}

func UpdateLocationType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	locationTypeID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "failed to get location type id", http.StatusBadRequest)
		return
	}

	locationType, err := models.GetLocationTypeByID(ctx, locationTypeID)
	if err != nil {
		util.ErrorResponse(w, "failed to get location type", http.StatusBadRequest)
		return
	}

	if locationType.OrganizationID != user.OwnerID {
		util.ErrorResponse(w, "user does not have access to this view location type", http.StatusForbidden)
		return
	}

	request := models.LocationTypeUpdateRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	err = locationType.UpdateWithRequest(request)
	if err != nil {
		util.ErrorResponse(w, "failed to update location type", http.StatusBadRequest)
		return
	}

	response := responses.GenerateUpdateLocationTypeResponse(*locationType)
	util.JSONResponse(w, response, http.StatusOK)
}

func DeleteLocationType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	locationTypeID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "failed to get location type id", http.StatusBadRequest)
		return
	}

	locationType, err := models.GetLocationTypeByID(ctx, locationTypeID)
	if err != nil {
		util.ErrorResponse(w, "failed to get location type", http.StatusBadRequest)
		return
	}

	if locationType.OrganizationID != user.OwnerID {
		util.ErrorResponse(w, "user does not have access to this view location type", http.StatusForbidden)
		return
	}

	err = locationType.GetLocations(ctx)
	if err != nil && err != gorm.ErrRecordNotFound {
		util.ErrorResponse(w, "failed to get locations", http.StatusBadRequest)
		return
	}

	if len(locationType.Locations) > 0 {
		util.ErrorResponse(w, "can't delete location type because it is associated with one or more locations", http.StatusBadRequest)
		return
	}

	err = locationType.Delete(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to delete location type", http.StatusBadRequest)
		return
	}

	util.SuccessResponse(w, http.StatusOK)
}
