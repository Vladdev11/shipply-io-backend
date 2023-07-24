package ClientHandlers

import (
	"errors"
	"net/http"

	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/responses"
	"github.com/shipply-io/shipply-io-backend/util"
)

var (
	//ErrClientDoesNotHaveAccessToVendor is returned when a client does not have access to a vendor
	ErrClientDoesNotHaveAccessToVendor = errors.New("client does not have access to vendor")
)

func ListVendors(w http.ResponseWriter, r *http.Request) {
  
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

	request := models.VendorListRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	request.ClientID = user.Client.ID
	vendors, total, count, err := user.Client.GetVendors(ctx, request)

	response := responses.GenerateListVendorsResponse(vendors, total, count)
	util.JSONResponse(w, response, http.StatusOK)
}

func GetVendor(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	vendorID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "failed to get vendor id", http.StatusBadRequest)
		return
	}

	vendor, err := models.GetVendorByID(ctx, vendorID)
	if err != nil {
		util.ErrorResponse(w, "failed to get vendor", http.StatusBadRequest)
		return
	}

	if vendor.ClientID != user.OwnerID {
		util.ErrorResponse(w, "user does not have access to client", http.StatusForbidden)
		return
	}

	response := responses.GenerateGetVendorResponse(*vendor)
	util.JSONResponse(w, response, http.StatusOK)
}

func CreateVendor(w http.ResponseWriter, r *http.Request) {
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

	// Adds ClientID to request (always the same as user.Client.ID)
	request := models.VendorCreateRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	vendor := &models.Vendor{
		ClientID:        request.ClientID,
		Name:            request.Name,
		VendorAccountID: request.VendorAccountID,
	}

	err = vendor.Create(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	response := responses.GenerateCreateVendorResponse(*vendor)
	util.JSONResponse(w, response, http.StatusOK)
}

func UpdateVendor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	vendorID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	vendor, err := models.GetVendorByID(ctx, vendorID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if user.OwnerID != vendor.ClientID {
		util.ErrResponse(w, ErrClientDoesNotHaveAccessToVendor, http.StatusForbidden)
		return
	}

	request := &models.VendorUpdateRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	err = vendor.UpdateWithRequest(ctx, request)
	if err != nil {
		util.ErrorResponse(w, "failed to update vendor", http.StatusInternalServerError)
		return
	}

	response := responses.GenerateUpdateVendorResponse(*vendor)
	util.JSONResponse(w, response, http.StatusOK)
}

func DeleteVendor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	vendorID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	vendor, err := models.GetVendorByID(ctx, vendorID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if user.OwnerID != vendor.ClientID {
		util.ErrResponse(w, ErrClientDoesNotHaveAccessToVendor, http.StatusForbidden)
		return
	}

	err = vendor.Delete(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusInternalServerError)
		return
	}

	util.SuccessResponse(w, http.StatusOK)

}
