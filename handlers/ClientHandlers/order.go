package ClientHandlers

import (
	"net/http"

	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/responses"
	"github.com/shipply-io/shipply-io-backend/util"
)

func ListOrders(w http.ResponseWriter, r *http.Request) {
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

	request := models.OrdersListRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	orders, count, total, err := user.Client.GetOrders(ctx, request)
	if err != nil {
		util.ErrorResponse(w, "failed to get orders", http.StatusBadRequest)
		return
	}

	response := responses.GenerateListOrdersResponse(orders, count, total)

	util.JSONResponse(w, response, http.StatusOK)
}

func GetOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	orderID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "invalid order id", http.StatusBadRequest)
		return
	}

	order, err := models.GetOrderByID(ctx, orderID)
	if err != nil {
		util.ErrorResponse(w, "failed to get order", http.StatusBadRequest)
		return
	}

	err = order.GetStore(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get store", http.StatusBadRequest)
		return
	}

	if user.OwnerID != order.Store.ClientID {
		util.ErrorResponse(w, "user does not have access to this order", http.StatusForbidden)
		return
	}

	err = order.GetOrderItems(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get order items", http.StatusBadRequest)
		return
	}

	err = order.GetBillToAddress(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get bill to address", http.StatusBadRequest)
		return
	}

	err = order.GetShipToAddress(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get ship to address", http.StatusBadRequest)
		return
	}

	err = order.GetWarehouse(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get warehouse", http.StatusBadRequest)
		return
	}

	err = order.GetShippingMethod(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get shipping method", http.StatusBadRequest)
		return
	}

	err = order.GetTags(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get tags", http.StatusBadRequest)
		return
	}

	err = order.GetStatus(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get status", http.StatusBadRequest)
		return
	}

	response := responses.GenerateGetOrderResponse(r.Context(), *order)

	util.JSONResponse(w, response, http.StatusOK)

}
