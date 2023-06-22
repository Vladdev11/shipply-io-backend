package ClientHandlers

import (
	"net/http"

	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/util"
)

func ListOrders(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetClient()
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

	orders, count, total, err := user.Client.GetOrders(request)
	if err != nil {
		util.ErrorResponse(w, "failed to get orders", http.StatusInternalServerError)
		return
	}

	searchResults, err := models.ConvertOrdersToSearchResults(orders, total, count)
	if err != nil {
		util.ErrorResponse(w, "failed to convert orders to search results", http.StatusBadRequest)
		return
	}

	util.JSONResponse(w, searchResults, http.StatusOK)

}

func GetOrder(w http.ResponseWriter, r *http.Request) {

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

	order, err := models.GetOrderByID(orderID)
	if err != nil {
		util.ErrorResponse(w, "failed to get order", http.StatusBadRequest)
		return
	}

	err = order.GetStore()
	if err != nil {
		util.ErrorResponse(w, "failed to get store", http.StatusBadRequest)
		return
	}

	if user.OwnerID != order.Store.ClientID {
		util.ErrorResponse(w, "user does not have access to this order", http.StatusForbidden)
		return
	}

	util.JSONResponse(w, order.ConvertToReturnJSON(), http.StatusOK)

}
