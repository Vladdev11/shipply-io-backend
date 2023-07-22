package OrganizationHandlers

import (
	"net/http"

	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/responses"
	"github.com/shipply-io/shipply-io-backend/util"
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
		util.ErrorResponse(w, "failed to get organization", http.StatusUnauthorized)
		return
	}

	purchaseOrderID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "invalid purchase order id", http.StatusBadRequest)
		return
	}

	purchaseOrder, err := models.GetPurchaseOrderByID(ctx, purchaseOrderID)
	if err != nil {
		util.ErrorResponse(w, "failed to find purchase order", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsClientOwner(ctx, purchaseOrder.ClientID) {
		util.ErrorResponse(w, "user does not have access to this client", http.StatusForbidden)
		return
	}

	purchaseOrder.GetStatus(ctx)
	purchaseOrder.GetItems(ctx)
	purchaseOrder.GetTags(ctx)
	purchaseOrder.GetHistory(ctx)
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

	err = user.GetOrganization(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusUnauthorized)
		return
	}

	request := models.PurchaseOrderListRequest{}
	request.OrganizationID = user.Organization.ID
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	purchaseOrders, count, total, err := user.Organization.GetPurchaseOrders(ctx, request)
	if err != nil {
		util.ErrorResponse(w, "failed to get purchase orders", http.StatusInternalServerError)
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
		util.ErrorResponse(w, "failed to create purchase order", http.StatusInternalServerError)
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

	err = user.GetOrganization(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusUnauthorized)
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
		util.ErrorResponse(w, "invalid purchase order id", http.StatusBadRequest)
		return
	}

	purchaseOrder, err := models.GetPurchaseOrderByID(ctx, purchaseOrderID)
	if err != nil {
		util.ErrorResponse(w, "failed to find purchase order", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsClientOwner(ctx, purchaseOrder.ClientID) {
		util.ErrorResponse(w, "user does not have access to this client", http.StatusForbidden)
		return
	}

	err = purchaseOrder.UpdateWithRequest(ctx, request)
	if err != nil {
		util.ErrorResponse(w, "failed to update purchase order", http.StatusInternalServerError)
		return
	}

	err = purchaseOrder.GetItems(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get purchase order items", http.StatusBadRequest)
		return
	}

	err = purchaseOrder.GetTags(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get purchase order tags", http.StatusBadRequest)
		return
	}

	err = purchaseOrder.GetHistory(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get purchase order history", http.StatusBadRequest)
		return
	}

	err = purchaseOrder.GetStatus(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get purchase order status", http.StatusBadRequest)
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

	err = user.GetOrganization(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusUnauthorized)
		return
	}

	purchaseOrderID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "invalid purchase order id", http.StatusBadRequest)
		return
	}

	purchaseOrder, err := models.GetPurchaseOrderByID(ctx, purchaseOrderID)
	if err != nil {
		util.ErrorResponse(w, "failed to find purchase order", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsClientOwner(ctx, purchaseOrder.ClientID) {
		util.ErrorResponse(w, "user does not have access to this client", http.StatusForbidden)
		return
	}

	err = purchaseOrder.Delete(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to delete purchase order", http.StatusInternalServerError)
		return
	}

	util.SuccessResponse(w, http.StatusOK)

}

func PurchaseOrderHistoryCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
		return
	}

	if err = user.GetOrganization(ctx); err != nil {
		util.ErrResponse(w, err, http.StatusUnauthorized)
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
		util.ErrorResponse(w, "failed to find purchase order", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsClientOwner(ctx, purchaseOrder.ClientID) {
		util.ErrorResponse(w, "user does not have access to this purchase order", http.StatusForbidden)
		return
	}

	request.CreatedBy = user.ID
	PurchaseOrderHistory := request.ConvertToPurchaseOrderHistory()
	err = PurchaseOrderHistory.Create(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to create purchase order note", http.StatusBadRequest)
		return
	}

	util.SuccessResponse(w, http.StatusOK)

}
