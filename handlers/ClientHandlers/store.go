package ClientHandlers

import (
	"errors"
	"net/http"

	"github.com/shipply-io/shipply-io-backend/api/shopify"
	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/responses"
	"github.com/shipply-io/shipply-io-backend/util"
	"github.com/shipply-io/shipply-io-backend/validation"
)

var (
	//ErrClientDoesNotHaveAccessToStore is returned when a client does not have access to a store
	ErrClientDoesNotHaveAccessToStore = errors.New("client does not have access to store")
)

var (
	//ErrClientDoesNotHaveAccessToStore is returned when a client does not have access to a store
	ErrClientDoesNotHaveAccessToStore = errors.New("client does not have access to store")
)

func GetStore(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := models.UserFromContext(ctx)

	getStoreRequestData, err := validation.ParseRequestToGetStoreRequestData(r)

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

	store, err := models.GetStoreByID(ctx, getStoreRequestData.StoreID)
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
	user := models.UserFromContext(ctx)

	err := user.Client.GetStores(ctx)
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
	user := models.UserFromContext(ctx)

	activateStoreRequestData, err := validation.ParseRequestToActivateStoreRequestData(r)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	store, err := models.GetStoreByID(ctx, activateStoreRequestData.StoreID)
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

	store, err := models.GetStoreByID(ctx, updateStoreRequestData.StoreID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if user.OwnerID != store.ClientID {
		util.ErrResponse(w, ErrClientDoesNotHaveAccessToStore, http.StatusForbidden)
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

	store, err := models.GetStoreByID(ctx, deleteStoreRequestData.StoreID)
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
