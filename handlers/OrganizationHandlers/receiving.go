package OrganizationHandlers

import (
	"net/http"

	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/util"
)

func ReceivingListItems(w http.ResponseWriter, r *http.Request) {
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

	purchaseOrderID, err := util.GetIntFromPath(r, "purchase_order_id")
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

	err = purchaseOrder.GetItems(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get purchase order items", http.StatusBadRequest)
		return
	}

	purchaseOrderItems := []models.PurchaseOrderItemReturnJSON{}
	for _, item := range purchaseOrder.Items {

		//get locations levels for item
		err = item.GetLocationsAndLevels(ctx)
		if err != nil {
			util.ErrorResponse(w, "failed to get locations and levels", http.StatusBadRequest)
			return
		}

		//set to 0 to omit from json response
		item.ProductID = 0
		item.PurchaseOrderID = 0
		purchaseOrderItems = append(purchaseOrderItems, item.ConvertToReturnJSON(ctx))
	}

	util.JSONResponse(w, purchaseOrderItems, http.StatusOK)

}

func PurchaseOrderItemGetReceivingDetails(w http.ResponseWriter, r *http.Request) {
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

	purchaseOrderID, err := util.GetIntFromPath(r, "purchase_order_id")
	if err != nil {
		util.ErrorResponse(w, "invalid purchase order id", http.StatusBadRequest)
		return
	}

	purchaseOrder, err := models.GetPurchaseOrderByID(ctx, purchaseOrderID)
	if err != nil {
		util.ErrorResponse(w, "failed to find purchase order", http.StatusBadRequest)
		return
	}

	purchaseOrderItemID, err := util.GetIntFromPath(r, "item_id")
	if err != nil {
		util.ErrorResponse(w, "invalid purchase order id", http.StatusBadRequest)
		return
	}

	purchaseOrderItem, err := models.GetPurchaseOrderItemByID(ctx, purchaseOrderItemID)
	if err != nil {
		util.ErrorResponse(w, "failed to find purchase order item", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsClientOwner(ctx, purchaseOrder.ClientID) {
		util.ErrorResponse(w, "user does not have access to this client", http.StatusForbidden)
		return
	}

	//get product for item
	err = purchaseOrderItem.GetProduct(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get product", http.StatusBadRequest)
		return
	}

	//get locations levels for item
	err = purchaseOrderItem.GetLocationsAndLevels(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get locations and levels", http.StatusBadRequest)
		return
	}

	//get rejections for item
	err = purchaseOrderItem.GetRejections(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get rejections", http.StatusBadRequest)
		return
	}

	//get item history
	err = purchaseOrderItem.GetHistory(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get item history", http.StatusBadRequest)
		return
	}

	//set to 0 to omit from json response
	purchaseOrderItem.ProductID = 0
	purchaseOrderItem.PurchaseOrderID = 0

	util.JSONResponse(w, purchaseOrderItem.ConvertToReturnJSON(ctx), http.StatusOK)

}
