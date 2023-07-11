package OrganizationHandlers

import (
	"net/http"

	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/util"
)

func VendorList(w http.ResponseWriter, r *http.Request) {
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

	request := models.VendorListRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	request.OrganizationID = user.Organization.ID
	vendors, total, count, err := user.Organization.GetVendors(ctx, request)

	searchResults, err := models.ConvertVendorsToSearchResults(vendors, total, count)
	if err != nil {
		util.ErrorResponse(w, "failed to convert vendors to search results", http.StatusBadRequest)
		return
	}

	util.JSONResponse(w, searchResults, http.StatusOK)

}

func VendorCreate(w http.ResponseWriter, r *http.Request) {
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

	request := models.VendorCreateRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	if !user.Organization.IsClientOwner(ctx, request.ClientID) {
		util.ErrorResponse(w, "user does not have access to client", http.StatusForbidden)
		return
	}

	vendor := &models.Vendor{
		ClientID:        request.ClientID,
		Name:            request.Name,
		VendorAccountID: request.VendorAccountID,
	}

	err = vendor.Create(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to create vendor", http.StatusBadRequest)
		return
	}

	util.JSONResponse(w, vendor.ConvertToReturnJSON(), http.StatusOK)

}

func VendorGet(w http.ResponseWriter, r *http.Request) {
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

	if !user.Organization.IsClientOwner(ctx, vendor.ClientID) {
		util.ErrorResponse(w, "user does not have access to client", http.StatusForbidden)
		return
	}

	util.JSONResponse(w, vendor, http.StatusOK)

}

func VendorUpdate(w http.ResponseWriter, r *http.Request) {
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

	vendorID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "invalid vendor id", http.StatusBadRequest)
		return
	}

	vendor, err := models.GetVendorByID(ctx, vendorID)
	if err != nil {
		util.ErrorResponse(w, "failed to find vendor", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsClientOwner(ctx, vendor.ClientID) {
		util.ErrorResponse(w, "user does not have access to this client", http.StatusForbidden)
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

	vendorJSON := vendor.ConvertToReturnJSON()
	util.JSONResponse(w, vendorJSON, http.StatusOK)

}

func VendorDelete(w http.ResponseWriter, r *http.Request) {
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

	vendorID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "invalid vendor id", http.StatusBadRequest)
		return
	}

	vendor, err := models.GetVendorByID(ctx, vendorID)
	if err != nil {
		util.ErrorResponse(w, "failed to find vendor", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsClientOwner(ctx, vendor.ClientID) {
		util.ErrorResponse(w, "user does not have access to this client", http.StatusForbidden)
		return
	}

	purchaseOrders, err := vendor.GetPurchaseOrders(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get purchase orders", http.StatusInternalServerError)
		return
	}

	if len(purchaseOrders) > 0 {
		util.ErrorResponse(w, "cannot delete vendor because there is one or more purchase orders associated with it", http.StatusBadRequest)
		return
	}

	err = vendor.Delete(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to delete vendor", http.StatusInternalServerError)
		return
	}

	util.SuccessResponse(w, http.StatusOK)

}
