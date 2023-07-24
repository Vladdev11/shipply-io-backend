package ClientHandlers

import (
	"errors"
	"net/http"

	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/responses"
	"github.com/shipply-io/shipply-io-backend/util"
)

var (
	//ErrPurchaseOrderStatusInUse is returned when a purchase order status is in use and cannot be deleted
	ErrPurchaseOrderStatusInUse = errors.New("purchase order status is in use and cannot be deleted")
)

func ListPurchaseOrderStatuses(w http.ResponseWriter, r *http.Request) {
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

	request := models.PurchaseOrderStatusListRequest{}
	request.ClientID = user.Client.ID
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	statuses, total, count, err := user.Client.GetPurchaseOrderStatuses(ctx, request)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	response := responses.GenerateListPurchaseOrderStatusesResponse(statuses, total, count)
	util.JSONResponse(w, response, http.StatusOK)
}

func PurchaseOrderStatusCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetOrganization(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusUnauthorized)
		return
	}

	request := models.PurchaseOrderStatusCreateRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	if user.OwnerID != request.ClientID {
		util.ErrResponse(w, models.ErrUserDoesNotBelongToClient, http.StatusForbidden)
		return
	}

	purchaseOrderStatus := &models.PurchaseOrderStatus{
		ClientID:    request.ClientID,
		Name:        request.Name,
		StatusColor: request.StatusColor,
		TextColor:   request.TextColor,
	}

	err = purchaseOrderStatus.Create(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusInternalServerError)
		return
	}

	response := responses.GenerateCreatePurchaseOrderStatusResponse(*purchaseOrderStatus)
	util.JSONResponse(w, response, http.StatusOK)
}

func PurchaseOrderStatusDelete(w http.ResponseWriter, r *http.Request) {
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

	posID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	purchaseOrderStatus, err := models.GetPurchaseOrderStatusByID(ctx, posID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if user.OwnerID != purchaseOrderStatus.ClientID {
		util.ErrResponse(w, models.ErrUserDoesNotBelongToClient, http.StatusForbidden)
		return
	}

	if purchaseOrderStatus.IsInUse(ctx) {
		util.ErrResponse(w, ErrPurchaseOrderStatusInUse, http.StatusBadRequest)
		return
	}

	err = purchaseOrderStatus.Delete(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusInternalServerError)
		return
	}

	util.SuccessResponse(w, http.StatusOK)

}

func PurchaseOrderStatusUpdate(w http.ResponseWriter, r *http.Request) {
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

	posID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	purchaseOrderStatus, err := models.GetPurchaseOrderStatusByID(ctx, posID)
	if err != nil {
		util.ErrorResponse(w, "failed to find purchase order status", http.StatusBadRequest)
		return
	}

	if user.OwnerID != purchaseOrderStatus.ClientID {
		util.ErrResponse(w, err, http.StatusForbidden)
		return
	}

	request := models.PurchaseOrderStatusUpdateRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	err = purchaseOrderStatus.UpdateWithRequest(ctx, &request)
	if err != nil {
		util.ErrorResponse(w, "failed to update purchase order status", http.StatusInternalServerError)
		return
	}

	response := responses.GenerateUpdatePurchaseOrderStatusResponse(*purchaseOrderStatus)
	util.JSONResponse(w, response, http.StatusOK)
}
