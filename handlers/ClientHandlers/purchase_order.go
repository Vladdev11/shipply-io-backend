package ClientHandlers

import (
	"errors"
	"net/http"

	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/responses"
	"github.com/shipply-io/shipply-io-backend/util"
)

var (
	//ErrPurchaseOrderAttachmentDoesNotBelongToPurchaseOrder is returned when a purchase order attachment does not belong to a purchase order
	ErrPurchaseOrderAttachmentDoesNotBelongToPurchaseOrder = errors.New("purchase order attachment does not belong to purchase order")
)

func GetPurchaseOrder(w http.ResponseWriter, r *http.Request) {

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

	purchaseOrderID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	purchaseOrder, err := models.GetPurchaseOrderByID(ctx, purchaseOrderID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if user.OwnerID != purchaseOrder.ClientID {
		util.ErrResponse(w, models.ErrUserDoesNotBelongToClient, http.StatusForbidden)
		return
	}

	purchaseOrder.GetItems(ctx)
	purchaseOrder.GetTags(ctx)
	purchaseOrder.GetHistory(ctx)
	purchaseOrder.GetStatus(ctx)
	purchaseOrder.GetClient(ctx)

	response := responses.GenerateGetPurchaseOrderResponse(ctx, *purchaseOrder)
	util.JSONResponse(w, response, http.StatusOK)
}

func ListPurchaseOrders(w http.ResponseWriter, r *http.Request) {

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

	request := models.PurchaseOrderListRequest{}
	errors := request.ParseAndValidateRequest(r)
	request.ClientID = user.Client.ID
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	if user.Client.ID != request.ClientID {
		util.ErrResponse(w, models.ErrUserDoesNotBelongToClient, http.StatusForbidden)
		return
	}

	purchaseOrders, count, total, err := user.Client.GetPurchaseOrders(ctx, request)
	if err != nil {
		util.ErrResponse(w, err, http.StatusInternalServerError)
		return
	}

	for i := range purchaseOrders {
		purchaseOrders[i].GetStatus(ctx)
	}

	response := responses.GenerateListPurchaseOrdersResponse(purchaseOrders, count, total)
	util.JSONResponse(w, response, http.StatusOK)
}

func CreatePurchaseOrder(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	_, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	request := models.PurchaseOrderCreateRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	purchaseOrder := &models.PurchaseOrder{
		ClientID:       request.ClientID,
		PONumber:       request.PONumber,
		StatusID:       request.StatusID,
		ExpectedDate:   request.ExpectedDate.Time,
		ShipDate:       request.ShipDate.Time,
		ClosedDate:     request.ClosedDate.Time,
		VendorID:       request.VendorID,
		WarehouseID:    request.WarehouseID,
		TrackingNumber: request.TrackingNumber,
		TrackingURL:    request.TrackingURL,
	}

	err = purchaseOrder.Create(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusInternalServerError)
		return
	}

	purchaseOrder.GetStatus(ctx)

	response := responses.GenerateCreatePurchaseOrderResponse(*purchaseOrder)
	util.JSONResponse(w, response, http.StatusOK)
}

func UpdatePurchaseOrder(w http.ResponseWriter, r *http.Request) {
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

	request := models.PurchaseOrderUpdateRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	purchaseOrderID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	purchaseOrder, err := models.GetPurchaseOrderByID(ctx, purchaseOrderID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if user.OwnerID != purchaseOrder.ClientID {
		util.ErrResponse(w, models.ErrUserDoesNotBelongToClient, http.StatusForbidden)
		return
	}

	err = purchaseOrder.UpdateWithRequest(ctx, request)
	if err != nil {
		util.ErrorResponse(w, "failed to update purchase order", http.StatusInternalServerError)
		return
	}

	err = purchaseOrder.GetItems(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	err = purchaseOrder.GetTags(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	err = purchaseOrder.GetHistory(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	err = purchaseOrder.GetStatus(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	response := responses.GenerateUpdatePurchaseOrderResponse(ctx, *purchaseOrder)
	util.JSONResponse(w, response, http.StatusOK)

}

func DeletePurchaseOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	purchaseOrderID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	purchaseOrder, err := models.GetPurchaseOrderByID(ctx, purchaseOrderID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if user.OwnerID != purchaseOrder.ClientID {
		util.ErrResponse(w, models.ErrUserDoesNotBelongToClient, http.StatusForbidden)
		return
	}

	err = purchaseOrder.Delete(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusInternalServerError)
		return
	}

	util.SuccessResponse(w, http.StatusOK)

}

func CreatePurchaseOrderHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
		return
	}

	request := &models.PurchaseOrderHistoryCreateRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	purchaseOrder, err := models.GetPurchaseOrderByID(ctx, request.PurchaseOrderId)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if user.OwnerID != purchaseOrder.ClientID {
		util.ErrResponse(w, models.ErrUserDoesNotBelongToClient, http.StatusForbidden)
		return
	}

	request.CreatedBy = user.ID
	PurchaseOrderHistory := request.ConvertToPurchaseOrderHistory()
	err = PurchaseOrderHistory.Create(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	util.SuccessResponse(w, http.StatusOK)

}
