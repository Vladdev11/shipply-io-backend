package OrganizationHandlers

import (
	"net/http"

	"github.com/shipply-io/shipply-io-backend/api/shopify"
	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/util"
)

func ActivateStore(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetOrganization()
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusUnauthorized)
		return
	}

	storeID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "Invalid store id", http.StatusBadRequest)
		return
	}

	store, err := models.GetStoreByID(storeID)
	if err != nil {
		util.ErrorResponse(w, "failed to get store", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsClientOwner(store.ClientID) {
		util.ErrorResponse(w, "user does not have access to store", http.StatusForbidden)
		return
	}

	err = store.Activate()
	if err != nil {
		util.ErrorResponse(w, "failed to activate store", http.StatusBadRequest)
		return
	}

	util.SuccessResponse(w, http.StatusOK)
}

func DeactivateStore(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetOrganization()
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusUnauthorized)
		return
	}

	storeID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "Invalid store id", http.StatusBadRequest)
		return
	}

	store, err := models.GetStoreByID(storeID)
	if err != nil {
		util.ErrorResponse(w, "failed to get store", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsClientOwner(store.ClientID) {
		util.ErrorResponse(w, "user does not have access to store", http.StatusForbidden)
		return
	}

	err = store.Deactivate()
	if err != nil {
		util.ErrorResponse(w, "failed to deactivate store", http.StatusBadRequest)
		return
	}

	util.SuccessResponse(w, http.StatusOK)
}

func UpdateStore(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetOrganization()
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusUnauthorized)
		return
	}

	storeID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "invalid store id", http.StatusBadRequest)
		return
	}

	store, err := models.GetStoreByID(storeID)
	if err != nil {
		util.ErrorResponse(w, "failed to find store", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsClientOwner(store.ClientID) {
		util.ErrorResponse(w, "user does not have access to this store", http.StatusForbidden)
		return
	}

	request := &models.StoreUpdateRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	err = store.UpdateWithRequest(request)
	if err != nil {
		util.ErrorResponse(w, "failed to update store", http.StatusInternalServerError)
		return
	}

	storeJSON := store.ConvertToReturnJSON()
	util.JSONResponse(w, storeJSON, http.StatusOK)

}

func DeleteStore(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetOrganization()
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusUnauthorized)
		return
	}

	storeID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "Invalid store id", http.StatusBadRequest)
		return
	}

	store, err := models.GetStoreByID(storeID)
	if err != nil {
		util.ErrorResponse(w, "failed to get store", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsClientOwner(store.ClientID) {
		util.ErrorResponse(w, "user does not have access to store", http.StatusForbidden)
		return
	}

	if store.MarketplaceID == util.ShopifyMarketplaceID {

		shopifyCredentials, err := store.GetShopifyCredentials()
		if err != nil {
			util.ErrorResponse(w, "failed to get shopify credentials", http.StatusBadRequest)
			return
		}

		// Delete ShopifyLocations in DB
		store.DeleteShopifyLocations()

		err = shopify.UninstallApp(shopifyCredentials.ShopName, shopifyCredentials.AccessToken)
		if err != nil {
			util.ErrorResponse(w, "failed to uninstall app", http.StatusBadRequest)
			return
		}

	}

	err = store.Delete()
	if err != nil {
		util.ErrorResponse(w, "failed to delete store", http.StatusBadRequest)
		return
	}

	util.SuccessResponse(w, http.StatusOK)
}

func ListStores(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetOrganization()
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

	stores, err := user.Organization.GetStores(request)
	if err != nil {
		util.ErrorResponse(w, "failed to find shipping methods", http.StatusBadRequest)
		return
	}

	storesJSON := []models.StoreReturnJSON{}

	for _, store := range stores {
		storesJSON = append(storesJSON, *store.ConvertToReturnJSON())
	}

	util.JSONResponse(w, storesJSON, http.StatusOK)

}

func GetStore(w http.ResponseWriter, r *http.Request) {

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetOrganization()
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusUnauthorized)
		return
	}

	storeID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "Invalid store id", http.StatusBadRequest)
		return
	}

	store, err := models.GetStoreByID(storeID)
	if err != nil {
		util.ErrorResponse(w, "failed to get store", http.StatusBadRequest)
		return
	}

	if !user.Organization.IsClientOwner(store.ClientID) {
		util.ErrorResponse(w, "user does not have access to store", http.StatusForbidden)
		return
	}

	util.JSONResponse(w, store.ConvertToReturnJSON(), http.StatusOK)

}
