package OrganizationHandlers

import (
	"net/http"

	"github.com/shipply-io/shipply-io-backend/api/shopify"
	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/responses"
	"github.com/shipply-io/shipply-io-backend/util"
)

func GetStore(w http.ResponseWriter, r *http.Request) {

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

	storeID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "Invalid store id", http.StatusBadRequest)
		return
	}

	store, err := models.GetStoreByID(ctx, storeID)
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

	storeID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "Invalid store id", http.StatusBadRequest)
		return
	}

	store, err := models.GetStoreByID(ctx, storeID)
	if err != nil {
		util.ErrorResponse(w, "failed to get store", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsClientOwner(ctx, store.ClientID) {
		util.ErrorResponse(w, "user does not have access to store", http.StatusForbidden)
		return
	}

	err = store.Activate(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to activate store", http.StatusBadRequest)
		return
	}

	util.SuccessResponse(w, http.StatusOK)
}

func DeactivateStore(w http.ResponseWriter, r *http.Request) {
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

	storeID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "Invalid store id", http.StatusBadRequest)
		return
	}

	store, err := models.GetStoreByID(ctx, storeID)
	if err != nil {
		util.ErrorResponse(w, "failed to get store", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsClientOwner(ctx, store.ClientID) {
		util.ErrorResponse(w, "user does not have access to store", http.StatusForbidden)
		return
	}

	err = store.Deactivate(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to deactivate store", http.StatusBadRequest)
		return
	}

	util.SuccessResponse(w, http.StatusOK)
}

func UpdateStore(w http.ResponseWriter, r *http.Request) {
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

	storeID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "invalid store id", http.StatusBadRequest)
		return
	}

	store, err := models.GetStoreByID(ctx, storeID)
	if err != nil {
		util.ErrorResponse(w, "failed to find store", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsClientOwner(ctx, store.ClientID) {
		util.ErrorResponse(w, "user does not have access to this store", http.StatusForbidden)
		return
	}

	request := &models.StoreUpdateRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	err = store.UpdateWithRequest(ctx, request)
	if err != nil {
		util.ErrorResponse(w, "failed to update store", http.StatusInternalServerError)
		return
	}

	store.GetMarketplace()

	response := responses.GenerateUpdateStoreResponse(*store)
	util.JSONResponse(w, response, http.StatusOK)

}

func DeleteStore(w http.ResponseWriter, r *http.Request) {
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

	storeID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "Invalid store id", http.StatusBadRequest)
		return
	}

	store, err := models.GetStoreByID(ctx, storeID)
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
