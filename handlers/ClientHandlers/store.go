package ClientHandlers

import (
	"errors"
	"net/http"

	"github.com/shipply-io/shipply-io-backend/api/shopify"
	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/responses"
	"github.com/shipply-io/shipply-io-backend/util"
)

var (
	//ErrClientDoesNotHaveAccessToStore is returned when a client does not have access to a store
	ErrClientDoesNotHaveAccessToStore = errors.New("client does not have access to store")
)

func GetStore(w http.ResponseWriter, r *http.Request) {
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

	storeID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	store, err := models.GetStoreByID(ctx, storeID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if user.OwnerID != store.ClientID {
		util.ErrResponse(w, ErrClientDoesNotHaveAccessToStore, http.StatusForbidden)
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

	err = user.GetClient(ctx)
	if err != nil {
		util.ErrResponse(w, ErrGetClient, http.StatusUnauthorized)
		return
	}

	err = user.Client.GetStores(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	for i := range user.Client.Stores {
		user.Client.Stores[i].GetMarketplace()
	}

	response := responses.GenerateListStoresResponse(user.Client.Stores)
	util.JSONResponse(w, response, http.StatusOK)

}

func ActivateStore(w http.ResponseWriter, r *http.Request) {
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

	storeID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	store, err := models.GetStoreByID(ctx, storeID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if user.OwnerID != store.ClientID {
		util.ErrResponse(w, ErrClientDoesNotHaveAccessToStore, http.StatusForbidden)
		return
	}

	err = store.Activate(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
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

	err = user.GetClient(ctx)
	if err != nil {
		util.ErrResponse(w, ErrGetClient, http.StatusUnauthorized)
		return
	}

	storeID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	store, err := models.GetStoreByID(ctx, storeID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if user.OwnerID != store.ClientID {
		util.ErrResponse(w, ErrClientDoesNotHaveAccessToStore, http.StatusForbidden)
		return
	}

	err = store.Deactivate(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
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

	err = user.GetClient(ctx)
	if err != nil {
		util.ErrResponse(w, ErrGetClient, http.StatusUnauthorized)
		return
	}

	storeID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	store, err := models.GetStoreByID(ctx, storeID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if user.OwnerID != store.ClientID {
		util.ErrResponse(w, ErrClientDoesNotHaveAccessToStore, http.StatusForbidden)
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

	err = user.GetClient(ctx)
	if err != nil {
		util.ErrResponse(w, ErrGetClient, http.StatusUnauthorized)
		return
	}

	storeID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	store, err := models.GetStoreByID(ctx, storeID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if user.OwnerID != store.ClientID {
		util.ErrResponse(w, ErrClientDoesNotHaveAccessToStore, http.StatusForbidden)
		return
	}

	if store.MarketplaceID == util.ShopifyMarketplaceID {

		shopifyCredentials, err := store.GetShopifyCredentials(ctx)
		if err != nil {
			util.ErrResponse(w, err, http.StatusBadRequest)
			return
		}

		//TODO add error handling here
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
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	util.SuccessResponse(w, http.StatusOK)
}
