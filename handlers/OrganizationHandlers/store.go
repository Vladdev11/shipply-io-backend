package OrganizationHandlers

import (
	"net/http"

	"github.com/shipply-io/shipply-io-backend/api/shopify"
	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/responses"
	"github.com/shipply-io/shipply-io-backend/util"
	"github.com/shipply-io/shipply-io-backend/validation"
)

func GetStore(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	user := models.UserFromContext(ctx)

	getStoreRequestData, err := validation.ParseRequestToGetStoreRequestData(r)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	store, err := models.GetStoreByID(ctx, getStoreRequestData.StoreID)
	if err != nil {
		util.ErrorResponse(w, "failed to get store", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsClientOwner(ctx, store.ClientID) {
		util.ErrorResponse(w, "user does not have access to store", http.StatusForbidden)
		return
	}

	store.GetMarketplace()

	response := responses.GenerateGetStoreResponse(*store)
	util.JSONResponse(w, response, http.StatusOK)
}

func ListStores(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	user := models.UserFromContext(ctx)

	// TODO clean and refactor request and query
	request := models.StoreListRequest{}
	request.OrganizationID = user.Organization.ID
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	stores, err := user.Organization.GetStores(ctx, request)
	if err != nil {
		util.ErrorResponse(w, "failed to find shipping methods", http.StatusBadRequest)
		return
	}

	for i := range stores {
		stores[i].GetMarketplace()
	}

	response := responses.GenerateListStoresResponse(stores)
	util.JSONResponse(w, response, http.StatusOK)

}

func ActivateStore(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := models.UserFromContext(ctx)

	activateStoreRequestData, err := validation.ParseRequestToActivateStoreRequestData(r)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	store, err := models.GetStoreByID(ctx, activateStoreRequestData.StoreID)
	if err != nil {
		util.ErrorResponse(w, "failed to get store", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsClientOwner(ctx, store.ClientID) {
		util.ErrorResponse(w, "user does not have access to store", http.StatusForbidden)
		return
	}

	_, err = models.UpdateStore(ctx, models.UpdateStoreInput{
		ID:     store.ID,
		Active: util.BoolPointer(true),
	})
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	util.SuccessResponse(w, http.StatusOK)
}

func DeactivateStore(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := models.UserFromContext(ctx)

	deactivateStoreRequestData, err := validation.ParseRequestToDeactivateStoreRequestData(r)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	store, err := models.GetStoreByID(ctx, deactivateStoreRequestData.StoreID)
	if err != nil {
		util.ErrorResponse(w, "failed to get store", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsClientOwner(ctx, store.ClientID) {
		util.ErrorResponse(w, "user does not have access to store", http.StatusForbidden)
		return
	}

	_, err = models.UpdateStore(ctx, models.UpdateStoreInput{
		ID:     store.ID,
		Active: util.BoolPointer(false),
	})
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	util.SuccessResponse(w, http.StatusOK)
}

func UpdateStore(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	user := models.UserFromContext(ctx)

	updateStoreRequestData, err := validation.ParseRequestToUpdateStoreRequestData(r)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	store, err := models.GetStoreByID(ctx, updateStoreRequestData.StoreID)
	if err != nil {
		util.ErrorResponse(w, "failed to find store", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsClientOwner(ctx, store.ClientID) {
		util.ErrorResponse(w, "user does not have access to this store", http.StatusForbidden)
		return
	}

	store, err = models.UpdateStore(ctx, models.UpdateStoreInput{
		ID:   store.ID,
		Name: updateStoreRequestData.Name,
	})
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	err = store.GetMarketplace()
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	response := responses.GenerateUpdateStoreResponse(*store)
	util.JSONResponse(w, response, http.StatusOK)

}

func DeleteStore(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := models.UserFromContext(ctx)

	deleteStoreRequestData, err := validation.ParseRequestToDeleteStoreRequestData(r)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	store, err := models.GetStoreByID(ctx, deleteStoreRequestData.StoreID)
	if err != nil {
		util.ErrorResponse(w, "failed to get store", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsClientOwner(ctx, store.ClientID) {
		util.ErrorResponse(w, "user does not have access to store", http.StatusForbidden)
		return
	}

	if store.MarketplaceID == util.ShopifyMarketplaceID {

		shopifyCredentials, err := store.GetShopifyCredentials(ctx)
		if err != nil {
			util.ErrorResponse(w, "failed to get shopify credentials", http.StatusBadRequest)
			return
		}

		// Delete ShopifyLocations in DB
		store.DeleteShopifyLocations(ctx)

		err = shopify.UninstallApp(shopifyCredentials.ShopName, shopifyCredentials.AccessToken)
		if err != nil {
			util.ErrorResponse(w, "failed to uninstall app", http.StatusBadRequest)
			return
		}

	}

	err = store.Delete(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to delete store", http.StatusBadRequest)
		return
	}

	util.SuccessResponse(w, http.StatusOK)
}
