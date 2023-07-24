package ClientHandlers

import (
	"errors"
	"net/http"

	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/responses"
	"github.com/shipply-io/shipply-io-backend/util"
)

var (
	//ErrPurchaseOrderDoesNotBelongToClient is returned when a user tries to access a purchase order that does not belong to their client
	ErrPurchaseOrderDoesNotBelongToClient = errors.New("purchase order does not belong to client")
	//ErrPurchaseOrderItemDoesNotBelongToPurchaseOrder is returned when a user tries to access a purchase order item that does not belong to the purchase order
	ErrPurchaseOrderItemDoesNotBelongToPurchaseOrder = errors.New("purchase order item does not belong to purchase order")
)

func CreatePurchaseOrderItem(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
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

	request := models.PurchaseOrderItemCreateRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	product, err := models.GetProductByID(ctx, request.ProductID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if product.ClientID != purchaseOrder.ClientID {
		util.ErrResponse(w, ErrProductDoesNotBelongToClient, http.StatusBadRequest)
		return
	}

	purchaseOrderItem := &models.PurchaseOrderItem{
		PurchaseOrderID: purchaseOrder.ID,
		ProductID:       product.ID,
		Product:         product,
	}

	err = purchaseOrderItem.Create(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	util.SuccessResponse(w, http.StatusOK)
}

func GetPurchaseOrderItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
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

	purchaseOrderItemID, err := util.GetIntFromPath(r, "item_id")
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	purchaseOrderItem, err := models.GetPurchaseOrderItemByID(ctx, purchaseOrderItemID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if purchaseOrderItem.PurchaseOrderID != purchaseOrder.ID {
		util.ErrResponse(w, ErrPurchaseOrderDoesNotBelongToClient, http.StatusBadRequest)
		return
	}

	err = purchaseOrderItem.GetProduct(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	response := responses.GenerateGetPurchaseOrderItemResponse(ctx, *purchaseOrderItem)
	util.JSONResponse(w, response, http.StatusOK)

}

func UpdatePurchaseOrderItemBulk(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
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

	purchaseOrderItemUpdateRequests, errs := models.ParseAndValidatePurchaseOrderItemUpdateRequests(r)
	if errs != nil {
		util.ErrorsResponse(w, errs, http.StatusBadRequest)
		return
	}

	itemsMap := make(map[int]struct {
		POItem  *models.PurchaseOrderItem
		Request models.PurchaseOrderItemUpdateRequest
	})

	for _, purchaseOrderItemUpdateRequest := range purchaseOrderItemUpdateRequests {

		request := purchaseOrderItemUpdateRequest

		purchaseOrderItem, err := models.GetPurchaseOrderItemByID(ctx, purchaseOrderItemUpdateRequest.ID)
		if err != nil {
			util.ErrResponse(w, err, http.StatusBadRequest)
			return
		}

		if purchaseOrderItem.PurchaseOrderID != purchaseOrder.ID {
			util.ErrResponse(w, ErrPurchaseOrderItemDoesNotBelongToPurchaseOrder, http.StatusBadRequest)
			return
		}

		itemsMap[purchaseOrderItem.ID] = struct {
			POItem  *models.PurchaseOrderItem
			Request models.PurchaseOrderItemUpdateRequest
		}{
			POItem:  purchaseOrderItem,
			Request: request,
		}
	}

	for _, item := range itemsMap {
		err := item.POItem.UpdateWithRequest(ctx, &item.Request)
		if err != nil {
			util.ErrorResponse(w, "failed to update purchase order item", http.StatusInternalServerError)
			return
		}
	}

	util.JSONResponse(w, util.JSONSuccess(), http.StatusOK)

}

func UpdatePurchaseOrderItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
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

	purchaseOrderItemID, err := util.GetIntFromPath(r, "item_id")
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	purchaseOrderItem, err := models.GetPurchaseOrderItemByID(ctx, purchaseOrderItemID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if purchaseOrderItem.PurchaseOrderID != purchaseOrder.ID {
		util.ErrResponse(w, ErrPurchaseOrderItemDoesNotBelongToPurchaseOrder, http.StatusBadRequest)
		return
	}

	request := &models.PurchaseOrderItemUpdateRequest{}
	errors := request.ParseAndValidateRequest(r)
	if len(errors) > 0 {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	err = purchaseOrderItem.UpdateWithRequest(ctx, request)
	if err != nil {
		util.ErrorResponse(w, "failed to update purchase order item", http.StatusInternalServerError)
		return
	}

	response := responses.GenerateUpdatePurchaseOrderItemResponse(ctx, *purchaseOrderItem)
	util.JSONResponse(w, response, http.StatusOK)
}

func DeletePurchaseOrderItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
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

	purchaseOrderItemID, err := util.GetIntFromPath(r, "item_id")
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	purchaseOrderItem, err := models.GetPurchaseOrderItemByID(ctx, purchaseOrderItemID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if purchaseOrderItem.PurchaseOrderID != purchaseOrder.ID {
		util.ErrResponse(w, ErrPurchaseOrderDoesNotBelongToClient, http.StatusBadRequest)
		return
	}

	err = purchaseOrderItem.Delete(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusInternalServerError)
		return
	}

	util.SuccessResponse(w, http.StatusOK)

}
