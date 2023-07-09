package OrganizationHandlers

import (
	"net/http"

	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/util"
)

func ListOrders(w http.ResponseWriter, r *http.Request) {
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

	request := models.OrdersListRequest{}
	request.OrganizationID = user.Organization.ID
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	orders, count, total, err := user.Organization.GetOrders(ctx, request)
	if err != nil {
		util.ErrorResponse(w, "failed to get orders", http.StatusBadRequest)
		return
	}

	searchResults, err := models.ConvertOrdersToSearchResults(ctx, orders, total, count)
	if err != nil {
		util.ErrorResponse(w, "failed to convert orders to search results", http.StatusBadRequest)
		return
	}

	util.JSONResponse(w, searchResults, http.StatusOK)
}

func GetOrder(w http.ResponseWriter, r *http.Request) {
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

	orderID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "Invalid order id", http.StatusBadRequest)
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

	if !user.Organization.IsClientOwner(ctx, order.Store.ClientID) {
		util.ErrorResponse(w, "user does not have access to order", http.StatusForbidden)
		return
	}

	util.JSONResponse(w, order.ConvertToReturnJSON(ctx), http.StatusOK)

}
