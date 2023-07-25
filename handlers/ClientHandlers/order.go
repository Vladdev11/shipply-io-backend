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
		util.ErrResponse(w, ErrGetClient, http.StatusUnauthorized)
		return
	}

	request := models.OrdersListRequest{}
	errs := request.ParseAndValidateRequest(r)
	if errs != nil {
		util.ErrResponse(w, errs, http.StatusBadRequest)
		return
	}

	orders, count, total, err := user.Client.GetOrders(ctx, request)
	if err != nil {
		util.ErrResponse(w, ErrGetOrders, http.StatusBadRequest)
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
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	order, err := models.GetOrderByID(ctx, orderID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	err = order.GetStore(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if user.OwnerID != order.Store.ClientID {
		util.ErrResponse(w, ErrStoreDoesNotBelongToClient, http.StatusForbidden)
		return
	}

	err = order.GetOrderItems(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	err = order.GetBillToAddress(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	err = order.GetShipToAddress(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	err = order.GetWarehouse(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	err = order.GetShippingMethod(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	err = order.GetTags(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	err = order.GetStatus(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	response := responses.GenerateGetOrderResponse(r.Context(), *order)

	util.JSONResponse(w, response, http.StatusOK)

}
