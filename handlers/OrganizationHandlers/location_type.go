package OrganizationHandlers

import (
	"net/http"

	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"
)

func LocationTypeList(w http.ResponseWriter, r *http.Request) {
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

	locationTypes, err := user.Organization.GetLocationTypes()
	if err != nil && err != gorm.ErrRecordNotFound {
		util.ErrorResponse(w, "failed to get location types", http.StatusBadRequest)
		return
	}

	var locationTypesReturnJSON []*models.LocationTypeReturnJSON
	for _, locationType := range locationTypes {
		locationTypesReturnJSON = append(locationTypesReturnJSON, locationType.ConvertToReturnJSON())
	}

	if len(locationTypesReturnJSON) == 0 {
		locationTypesReturnJSON = make([]*models.LocationTypeReturnJSON, 0)
	}

	util.JSONResponse(w, locationTypesReturnJSON, http.StatusOK)
}

func LocationTypeCreate(w http.ResponseWriter, r *http.Request) {
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

	err = locationType.Create()
	if err != nil {
		util.ErrorResponse(w, "failed to create location type", http.StatusBadRequest)
		return
	}

	util.JSONResponse(w, locationType.ConvertToReturnJSON(), http.StatusOK)
}

func LocationTypeGet(w http.ResponseWriter, r *http.Request) {
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

	locationType, err := models.GetLocationTypeByID(locationTypeID)
	if err != nil {
		util.ErrorResponse(w, "failed to get location type", http.StatusBadRequest)
		return
	}

	if locationType.OrganizationID != user.OwnerID {
		util.ErrorResponse(w, "user does not have access to this view location type", http.StatusForbidden)
		return
	}

	util.JSONResponse(w, locationType.ConvertToReturnJSON(), http.StatusOK)
}

func LocationTypeUpdate(w http.ResponseWriter, r *http.Request) {
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

	locationType, err := models.GetLocationTypeByID(locationTypeID)
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

	util.JSONResponse(w, locationType.ConvertToReturnJSON(), http.StatusOK)
}

func LocationTypeDelete(w http.ResponseWriter, r *http.Request) {
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

	locationType, err := models.GetLocationTypeByID(locationTypeID)
	if err != nil {
		util.ErrorResponse(w, "failed to get location type", http.StatusBadRequest)
		return
	}

	if locationType.OrganizationID != user.OwnerID {
		util.ErrorResponse(w, "user does not have access to this view location type", http.StatusForbidden)
		return
	}

	err = locationType.GetLocations()
	if err != nil && err != gorm.ErrRecordNotFound {
		util.ErrorResponse(w, "failed to get locations", http.StatusBadRequest)
		return
	}

	if len(locationType.Locations) > 0 {
		util.ErrorResponse(w, "can't delete location type because it is associated with one or more locations", http.StatusBadRequest)
		return
	}

	err = locationType.Delete()
	if err != nil {
		util.ErrorResponse(w, "failed to delete location type", http.StatusBadRequest)
		return
	}

	util.SuccessResponse(w, http.StatusOK)
}
