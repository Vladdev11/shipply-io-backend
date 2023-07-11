package ClientHandlers

import (
	"net/http"

	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/util"
)

func PurchaseOrderItemCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
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

	if user.OwnerID != purchaseOrder.ClientID {
		util.ErrorResponse(w, "user does not have access to this client", http.StatusForbidden)
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
		util.ErrorResponse(w, "failed to find product", http.StatusBadRequest)
		return
	}

	if product.ClientID != purchaseOrder.ClientID {
		util.ErrorResponse(w, "product does not belong to this client", http.StatusBadRequest)
		return
	}

	purchaseOrderItem := &models.PurchaseOrderItem{
		PurchaseOrderID: purchaseOrder.ID,
		ProductID:       product.ID,
		Product:         product,
	}

	err = purchaseOrderItem.Create(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to create purchase order item", http.StatusBadRequest)
		return
	}

	purchaseOrderItemJson := purchaseOrderItem.ConvertToReturnJSON(r.Context())
	util.JSONResponse(w, purchaseOrderItemJson, http.StatusCreated)
}

func PurchaseOrderItemGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
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

	if user.OwnerID != purchaseOrder.ClientID {
		util.ErrorResponse(w, "user does not have access to this client", http.StatusForbidden)
		return
	}

	purchaseOrderItemID, err := util.GetIntFromPath(r, "item_id")
	if err != nil {
		util.ErrorResponse(w, "invalid purchase order item id", http.StatusBadRequest)
		return
	}

	purchaseOrderItem, err := models.GetPurchaseOrderItemByID(ctx, purchaseOrderItemID)
	if err != nil {
		util.ErrorResponse(w, "failed to find purchase order item", http.StatusBadRequest)
		return
	}

	if purchaseOrderItem.PurchaseOrderID != purchaseOrder.ID {
		util.ErrorResponse(w, "purchase order item does not belong to this purchase order", http.StatusBadRequest)
		return
	}

	err = purchaseOrderItem.GetProduct(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to find product", http.StatusBadRequest)
		return
	}

	purchaseOrderItemJson := purchaseOrderItem.ConvertToReturnJSON(r.Context())
	util.JSONResponse(w, purchaseOrderItemJson, http.StatusOK)

}

func PurchaseOrderItemBulkUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
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

	if user.OwnerID != purchaseOrder.ClientID {
		util.ErrorResponse(w, "user does not have access to this client", http.StatusForbidden)
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
			util.ErrorResponse(w, "failed to find purchase order item", http.StatusBadRequest)
			return
		}

		if purchaseOrderItem.PurchaseOrderID != purchaseOrder.ID {
			util.ErrorResponse(w, "purchase order item does not belong to this purchase order", http.StatusBadRequest)
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

func PurchaseOrderItemUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
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

	if user.OwnerID != purchaseOrder.ClientID {
		util.ErrorResponse(w, "user does not have access to this client", http.StatusForbidden)
		return
	}

	purchaseOrderItemID, err := util.GetIntFromPath(r, "item_id")
	if err != nil {
		util.ErrorResponse(w, "invalid purchase order item id", http.StatusBadRequest)
		return
	}

	purchaseOrderItem, err := models.GetPurchaseOrderItemByID(ctx, purchaseOrderItemID)
	if err != nil {
		util.ErrorResponse(w, "failed to find purchase order item", http.StatusBadRequest)
		return
	}

	if purchaseOrderItem.PurchaseOrderID != purchaseOrder.ID {
		util.ErrorResponse(w, "purchase order item does not belong to this purchase order", http.StatusBadRequest)
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

	purchaseOrderItemJson := purchaseOrderItem.ConvertToReturnJSON(r.Context())
	util.JSONResponse(w, purchaseOrderItemJson, http.StatusOK)
}

func PurchaseOrderItemDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
		return
	}

	purchaseOrderID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "invalid purchase order ID", http.StatusBadRequest)
		return
	}

	purchaseOrder, err := models.GetPurchaseOrderByID(ctx, purchaseOrderID)
	if err != nil {
		util.ErrorResponse(w, "failed to find purchase order", http.StatusBadRequest)
		return
	}

	if user.OwnerID != purchaseOrder.ClientID {
		util.ErrorResponse(w, "user does not have access to this client", http.StatusForbidden)
		return
	}

	purchaseOrderItemID, err := util.GetIntFromPath(r, "item_id")
	if err != nil {
		util.ErrorResponse(w, "invalid purchase order item id", http.StatusBadRequest)
		return
	}

	purchaseOrderItem, err := models.GetPurchaseOrderItemByID(ctx, purchaseOrderItemID)
	if err != nil {
		util.ErrorResponse(w, "failed to find purchase order item", http.StatusBadRequest)
		return
	}

	if purchaseOrderItem.PurchaseOrderID != purchaseOrder.ID {
		util.ErrorResponse(w, "purchase order item does not belong to this purchase order", http.StatusBadRequest)
		return
	}

	err = purchaseOrderItem.Delete(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to delete purchase order item", http.StatusInternalServerError)
		return
	}

	util.SuccessResponse(w, http.StatusOK)

}
