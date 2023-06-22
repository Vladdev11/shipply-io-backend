package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/tasks"
	"github.com/shipply-io/shipply-io-backend/util"

	"github.com/shipply-io/shipply-io-backend/api/shopify"
	Shopify "github.com/shipply-io/shipply-io-backend/api/shopify"
)

// ** START INSTALLATION HANDLERS ** //

func ShopifyOAuth(w http.ResponseWriter, r *http.Request) {

	// TODO add error messages that are friendly to the user

	err := models.LogShopifyInstallRequest(r)
	if err != nil {
		models.CreateSystemError(err.Error())
	}

	shop, err := util.GetStringQueryParam(r, "shop")
	if err != nil {
		models.CreateSystemError("shopify install failed: missing shop parameter")
		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify?error=%s", util.FrontEndBaseURL, url.QueryEscape("shop parameter is required")), http.StatusFound)
		return
	}

	clientID, err := util.GetIntQueryParam(r, "client_id")
	if err != nil {
		models.CreateSystemError("shopify install failed: missing client_id parameter")
		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify?error=%s&shop_name=%s", util.FrontEndBaseURL, url.QueryEscape("client_id parameter is required"), shop), http.StatusFound)
		return
	}

	nonce := uuid.New()

	existingStore, _ := models.GetShopifyStoreByShopName(shop)
	if existingStore != nil {

		shopifyCredentials, err := existingStore.GetShopifyCredentials()
		if err != nil {

			models.SystemError{
				Message:  err.Error(),
				ClientID: clientID,
			}.Create()

			util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify?error=%s&shop_name=%s", util.FrontEndBaseURL, url.QueryEscape("failed to get shopify credentials"), shop), http.StatusFound)
			return
		}

		if shopifyCredentials.AccessToken != "" {
			util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify?error=%s&shop_name=%s", util.FrontEndBaseURL, url.QueryEscape("shopify store "+shop+" already associated with a client"), shop), http.StatusFound)
			return
		}

		err = existingStore.UpdateAPICredentials(&models.ShopifyCredentials{
			ShopName:            shop,
			ShopifyInstallNonce: nonce.String(),
		})
		if err != nil {

			models.SystemError{
				Message:  err.Error(),
				ClientID: clientID,
			}.Create()

			util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify?error=%s&shop_name=%s", util.FrontEndBaseURL, url.QueryEscape("failed to update API Credentials"), shop), http.StatusFound)
			return
		}

		err = existingStore.Activate()
		if err != nil {

			models.SystemError{
				Message:  err.Error(),
				ClientID: clientID,
			}.Create()

			util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify?error=%s&shop_name=%s", util.FrontEndBaseURL, url.QueryEscape("failed to activate store"), shop), http.StatusFound)
			return
		}

	} else {

		credentials := models.ShopifyCredentials{
			ShopName:            shop,
			ShopifyInstallNonce: nonce.String(),
		}

		credentialsJSON, err := json.Marshal(credentials)
		if err != nil {

			models.SystemError{
				Message:  err.Error(),
				ClientID: clientID,
			}.Create()

			util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify?error=%s&shop_name=%s", util.FrontEndBaseURL, url.QueryEscape("failed to marshal credentials"), shop), http.StatusFound)
			return
		}

		storeName, err := util.GetStringQueryParam(r, "store_name")
		if err != nil && err != util.ErrMissingQueryParam {

			models.SystemError{
				Message:  "Invalid store_name parameter",
				ClientID: clientID,
			}.Create()

		}

		if storeName == "" {
			storeName = strings.Split(shop, ".")[0]
		}

		store := models.Store{
			MarketplaceID:  util.ShopifyMarketplaceID,
			ClientID:       clientID,
			Name:           storeName,
			APICredentials: json.RawMessage(credentialsJSON),
			Active:         true,
		}

		err = store.Create()
		if err != nil {

			models.SystemError{
				Message:  err.Error(),
				ClientID: clientID,
			}.Create()

			util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify?error=%s&shop_name=%s", util.FrontEndBaseURL, url.QueryEscape("failed to create store"), shop), http.StatusFound)
			return
		}

	}

	url := Shopify.AuthorizeUrl(shop, nonce.String())

	util.RedirectResponse(w, r, url, http.StatusFound)

}

func ShopifyOAuthConfirmation(w http.ResponseWriter, r *http.Request) {

	err := models.LogShopifyInstallRequest(r)
	if err != nil {
		models.CreateSystemError(err.Error())
	}

	// TODO add error messages that are friendly to the user

	URL, err := url.ParseRequestURI(r.RequestURI)
	if err != nil {
		models.CreateSystemError(err.Error())
		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s", util.FrontEndBaseURL, err.Error()), http.StatusPermanentRedirect)
		return
	}

	if ok, err := Shopify.VerifyAuthorizationURL(URL); !ok {
		models.CreateSystemError(err.Error())
		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s", util.FrontEndBaseURL, err.Error()), http.StatusPermanentRedirect)
		return
	}

	shopName, err := util.GetStringQueryParam(r, "shop")
	if err != nil {
		models.CreateSystemError(err.Error())
		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s", util.FrontEndBaseURL, "invalid shop"), http.StatusPermanentRedirect)
		return
	}

	code, err := util.GetStringQueryParam(r, "code")
	if err != nil {
		models.CreateSystemError(err.Error())
		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontEndBaseURL, "invalid code", shopName), http.StatusPermanentRedirect)
		return
	}

	store, err := models.GetShopifyStoreByShopName(shopName)
	if err != nil {
		models.CreateSystemError(err.Error())
		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontEndBaseURL, "failed to find store with that shop name", shopName), http.StatusPermanentRedirect)
		return
	}

	shopifyCredentials, err := store.GetShopifyCredentials()
	if err != nil {

		models.SystemError{
			Message:  err.Error(),
			ClientID: store.ClientID,
		}.Create()

		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontEndBaseURL, "failed to get API credentials", shopName), http.StatusPermanentRedirect)
		return
	}

	// Verify the security checks
	err = Shopify.VerifyOAuthCallback(r.URL, shopifyCredentials.ShopifyInstallNonce)
	if err != nil {

		models.SystemError{
			Message:  err.Error(),
			ClientID: store.ClientID,
		}.Create()

		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontEndBaseURL, err.Error(), shopName), http.StatusPermanentRedirect)
		return
	}

	accessToken, err := Shopify.GetAccessToken(shopName, code)
	if err != nil {

		models.SystemError{
			Message:  err.Error(),
			ClientID: store.ClientID,
		}.Create()

		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontEndBaseURL, "failed to get access token", shopName), http.StatusPermanentRedirect)
		return
	}

	shopifyCredentials.AccessToken = accessToken
	err = store.UpdateAPICredentials(shopifyCredentials)
	if err != nil {

		models.SystemError{
			Message:  err.Error(),
			ClientID: store.ClientID,
		}.Create()

		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontEndBaseURL, "failed to update store with access token", shopName), http.StatusPermanentRedirect)
		return
	}

	client, err := models.GetClientByID(store.ClientID)
	if err != nil {

		models.SystemError{
			Message:  err.Error(),
			ClientID: store.ClientID,
		}.Create()

		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontEndBaseURL, "failed to get client", shopName), http.StatusPermanentRedirect)
		return
	}
	warehouses, err := models.GetWarehousesByOrganizationID(client.OrganizationID)
	if err != nil {

		models.SystemError{
			Message:  err.Error(),
			ClientID: client.ID,
		}.Create()

		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontEndBaseURL, "failed to get warehouses", shopName), http.StatusPermanentRedirect)
		return
	}

	// Create Locations in Shopify (Can't do in task because they must complete before we can activate inventory items)
	for _, warehouse := range warehouses {
		err = warehouse.GetShipFromAddress()
		if err != nil {

			models.SystemError{
				Message:  err.Error(),
				ClientID: client.ID,
			}.Create()

			util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontEndBaseURL, "failed to get warehouse address", shopName), http.StatusPermanentRedirect)
			return
		}

		locationGraphqlID, err := shopify.AddLocation(shopName, accessToken, warehouse)
		if err != nil {

			models.SystemError{
				Message:  "failed to create location: " + err.Error(),
				ClientID: client.ID,
			}.Create()

			util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontEndBaseURL, "failed to create location", shopName), http.StatusPermanentRedirect)
			return
		}

		shopifyLocation := models.ShopifyLocation{
			StoreID:           store.ID,
			ShopifyLocationID: locationGraphqlID,
			WarehouseID:       warehouse.ID,
		}
		err = shopifyLocation.Create()
		if err != nil {

			models.SystemError{
				Message:  "failed to create shopify location: " + err.Error(),
				ClientID: client.ID,
			}.Create()

			util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontEndBaseURL, "failed to create shopify location", shopName), http.StatusPermanentRedirect)
			return
		}
	}

	// Create json.RawMessage for the parameters needed for products task
	productsTaskParameters, err := json.Marshal(map[string]interface{}{
		"shop_name":    shopifyCredentials.ShopName,
		"access_token": shopifyCredentials.AccessToken,
		"store_id":     store.ID,
	})
	if err != nil {

		models.SystemError{
			Message:  err.Error(),
			ClientID: client.ID,
		}.Create()

		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontEndBaseURL, "An error occured while importing products", shopName), http.StatusPermanentRedirect)
		return
	}

	retrieveProductsTask := &models.SystemTask{
		TaskType: util.PullShopifyProductsTaskType,
		Payload:  json.RawMessage(productsTaskParameters),
	}
	err = retrieveProductsTask.Create()
	if err != nil {

		models.SystemError{
			Message:  err.Error(),
			ClientID: client.ID,
		}.Create()

		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontEndBaseURL, "An error occured while importing products", shopName), http.StatusPermanentRedirect)
		return
	}
	tasks.TaskChannel <- retrieveProductsTask

	// Create json.RawMessage for the parameters needed for orders task
	ordersTaskParameters, err := json.Marshal(map[string]interface{}{
		"shop_name":    shopifyCredentials.ShopName,
		"access_token": shopifyCredentials.AccessToken,
		"store_id":     store.ID,
	})
	if err != nil {

		models.SystemError{
			Message:  err.Error(),
			ClientID: client.ID,
		}.Create()

		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontEndBaseURL, "An error occured while importing orders", shopName), http.StatusPermanentRedirect)
		return
	}

	retrieveOrdersTask := &models.SystemTask{
		TaskType: util.PullShopifyOrdersTaskType,
		Payload:  json.RawMessage(ordersTaskParameters),
	}
	err = retrieveOrdersTask.Create()
	if err != nil {

		models.SystemError{
			Message:  err.Error(),
			ClientID: client.ID,
		}.Create()

		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontEndBaseURL, "An error occured while importing orders", shopName), http.StatusPermanentRedirect)
		return
	}
	tasks.TaskChannel <- retrieveOrdersTask

	err = shopify.SubscribeToAppUninstalledWebhook(shopifyCredentials.ShopName, shopifyCredentials.AccessToken)
	if err != nil {

		models.SystemError{
			Message:  "failed to subscribe to app uninstalled webhook: " + err.Error(),
			ClientID: client.ID,
		}.Create()

		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontEndBaseURL, "failed to subscribe to app uninstalled webhook", shopName), http.StatusPermanentRedirect)
		return
	}

	err = shopify.SubscribeToOrdersCreateWebhook(shopifyCredentials.ShopName, shopifyCredentials.AccessToken)
	if err != nil {

		models.SystemError{
			Message:  "failed to subscribe to orders create webhook: " + err.Error(),
			ClientID: client.ID,
		}.Create()

		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontEndBaseURL, "failed to subscribe to orders create webhook", shopName), http.StatusPermanentRedirect)
		return
	}

	err = shopify.SubscribeToOrdersUpdateWebhook(shopifyCredentials.ShopName, shopifyCredentials.AccessToken)
	if err != nil {

		models.SystemError{
			Message:  "failed to subscribe to orders update webhook: " + err.Error(),
			ClientID: client.ID,
		}.Create()

		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontEndBaseURL, "failed to subscribe to orders update webhook", shopName), http.StatusPermanentRedirect)
		return
	}

	// Product create webhook is for a new product with new variants
	err = shopify.SubscribeToProductsCreateWebhook(shopifyCredentials.ShopName, shopifyCredentials.AccessToken)
	if err != nil {

		models.SystemError{
			Message:  "failed to subscribe to products create webhook: " + err.Error(),
			ClientID: client.ID,
		}.Create()

		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontEndBaseURL, "failed to subscribe to products create webhook", shopName), http.StatusPermanentRedirect)
		return
	}

	// Product update webhook is for an existing product with new variants
	err = shopify.SubscribeToProductsUpdateWebhook(shopifyCredentials.ShopName, shopifyCredentials.AccessToken)
	if err != nil {

		models.SystemError{
			Message:  "failed to subscribe to products update webhook: " + err.Error(),
			ClientID: client.ID,
		}.Create()

		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontEndBaseURL, "failed to subscribe to products update webhook", shopName), http.StatusPermanentRedirect)
		return
	}

	util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/%d", util.FrontEndBaseURL, store.ID), http.StatusPermanentRedirect)
}

// ** END INSTALLATION HANDLERS ** //

// ** BEGIN WEBHOOK HANDLERS ** //
func ShopifyAppUninstalledWebhook(w http.ResponseWriter, r *http.Request) {

	shopName := r.Header.Get("X-Shopify-Shop-Domain")
	if shopName == "" {
		util.ErrorResponse(w, "invalid shop name", http.StatusBadRequest)
		return
	}

	store, err := models.GetShopifyStoreByShopName(shopName)
	if err != nil {
		util.ErrorResponse(w, fmt.Sprintf("No store with the ShopName %s exists", shopName), http.StatusInternalServerError)
		return
	}

	err = store.DeleteAccessToken()
	if err != nil {
		util.ErrorResponse(w, "failed to delete access token", http.StatusInternalServerError)
		return
	}

	err = store.DeleteShopifyLocations()
	if err != nil {
		util.ErrorResponse(w, "failed to delete shopify locations", http.StatusInternalServerError)
		return
	}

	err = store.Deactivate()
	if err != nil {
		util.ErrorResponse(w, "failed to deactivate store", http.StatusInternalServerError)
		return
	}

	util.SuccessResponse(w, http.StatusOK)
}

func ShopifyOrderCreatedWebhook(w http.ResponseWriter, r *http.Request) {

	shopName := r.Header.Get("X-Shopify-Shop-Domain")
	if shopName == "" {
		util.ErrorResponse(w, "invalid shop name", http.StatusBadRequest)
		return
	}

	store, err := models.GetShopifyStoreByShopName(shopName)
	if err != nil {
		util.ErrorResponse(w, fmt.Sprintf("No store with the ShopName %s exists", shopName), http.StatusInternalServerError)
		return
	}

	orderID, errors := shopify.ParseAndValidateOrderCreateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	order, err := shopify.RetrieveOrder(store.ShopifyShopName(), store.ShopifyAccessToken(), *orderID)
	if err != nil {
		util.ErrorResponse(w, "failed to retrieve order", http.StatusInternalServerError)
		return
	}

	gqlOrder, err := shopify.ConvertGetOrderById_OrderToShopifyModelOrder(*order)
	if err != nil {
		util.ErrorResponse(w, "failed to convert order", http.StatusInternalServerError)
		return
	}

	err = tasks.SyncShopifyGraphqlOrder(store.ID, store.ShopifyShopName(), store.ShopifyAccessToken(), gqlOrder)
	if err != nil {
		util.ErrorResponse(w, "failed to sync order", http.StatusInternalServerError)
		return
	}

	util.SuccessResponse(w, http.StatusOK)
}

func ShopifyOrderUpdatedWebhook(w http.ResponseWriter, r *http.Request) {

	shopName := r.Header.Get("X-Shopify-Shop-Domain")
	if shopName == "" {
		util.ErrorResponse(w, "invalid shop name", http.StatusBadRequest)
		return
	}

	store, err := models.GetShopifyStoreByShopName(shopName)
	if err != nil {
		util.ErrorResponse(w, fmt.Sprintf("No store with the ShopName %s exists", shopName), http.StatusInternalServerError)
		return
	}

	orderID, errors := shopify.ParseAndValidateOrderCreateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	order, err := shopify.RetrieveOrder(store.ShopifyShopName(), store.ShopifyAccessToken(), *orderID)
	if err != nil {
		util.ErrorResponse(w, "failed to retrieve order", http.StatusInternalServerError)
		return
	}

	gqlOrder, err := shopify.ConvertGetOrderById_OrderToShopifyModelOrder(*order)
	if err != nil {
		util.ErrorResponse(w, "failed to convert order", http.StatusInternalServerError)
		return
	}

	err = tasks.SyncShopifyGraphqlOrder(store.ID, store.ShopifyShopName(), store.ShopifyAccessToken(), gqlOrder)
	if err != nil {
		util.ErrorResponse(w, "failed to sync order", http.StatusInternalServerError)
		return
	}

	util.SuccessResponse(w, http.StatusOK)
}

func ShopifyProductCreatedWebhook(w http.ResponseWriter, r *http.Request) {

	shopName := r.Header.Get("X-Shopify-Shop-Domain")
	if shopName == "" {
		util.ErrorResponse(w, "invalid shop name", http.StatusBadRequest)
		return
	}

	store, err := models.GetShopifyStoreByShopName(shopName)
	if err != nil {
		util.ErrorResponse(w, fmt.Sprintf("No store with the ShopName %s exists", shopName), http.StatusInternalServerError)
		return
	}

	if !store.Active {
		util.SuccessResponse(w, http.StatusOK)
		return
	}

	productID, errors := shopify.ParseAndValidateProductCreateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	productVariants, err := shopify.RetrieveProductVariantsByProductID(store.ShopifyShopName(), store.ShopifyAccessToken(), *productID)
	if err != nil {
		util.ErrorResponse(w, "failed to retrieve product", http.StatusInternalServerError)
		return
	}

	for _, productVariant := range productVariants {
		gqlProduct, err := shopify.ConvertGetProductById_ProductToShopifyModelProduct(productVariant)
		if err != nil {
			models.SystemError{
				Message:  "failed to convert shopify product to model product: " + err.Error(),
				ClientID: store.ClientID,
			}.Create()
			continue
		}
		err = tasks.SyncShopifyGraphqlProduct(store.ID, store.ShopifyShopName(), store.ShopifyAccessToken(), *gqlProduct)
		if err != nil {
			models.SystemError{
				Message:  "failed to sync shopify product: " + err.Error(),
				ClientID: store.ClientID,
			}.Create()
			continue
		}

	}

	util.SuccessResponse(w, http.StatusOK)

}

func ShopifyProductUpdatedWebhook(w http.ResponseWriter, r *http.Request) {

	shopName := r.Header.Get("X-Shopify-Shop-Domain")
	if shopName == "" {
		util.ErrorResponse(w, "invalid shop name", http.StatusBadRequest)
		return
	}

	store, err := models.GetShopifyStoreByShopName(shopName)
	if err != nil {
		util.ErrorResponse(w, fmt.Sprintf("No store with the ShopName %s exists", shopName), http.StatusInternalServerError)
		return
	}

	if !store.Active {
		util.SuccessResponse(w, http.StatusOK)
		return
	}

	productID, errors := shopify.ParseAndValidateProductCreateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	productVariants, err := shopify.RetrieveProductVariantsByProductID(store.ShopifyShopName(), store.ShopifyAccessToken(), *productID)
	if err != nil {
		util.ErrorResponse(w, "failed to retrieve product", http.StatusInternalServerError)
		return
	}

	for _, productVariant := range productVariants {
		gqlProduct, err := shopify.ConvertGetProductById_ProductToShopifyModelProduct(productVariant)
		if err != nil {
			models.SystemError{
				Message:  "failed to convert shopify product to model product: " + err.Error(),
				ClientID: store.ClientID,
			}.Create()
			continue
		}
		err = tasks.SyncShopifyGraphqlProduct(store.ID, store.ShopifyShopName(), store.ShopifyAccessToken(), *gqlProduct)
		if err != nil {
			models.SystemError{
				Message:  "failed to sync shopify product: " + err.Error(),
				ClientID: store.ClientID,
			}.Create()
			continue
		}

	}

	util.SuccessResponse(w, http.StatusOK)
}

// ** MANDATORY WEBHOOKS ** //
func ShopifyCustomerRedactWebhook(w http.ResponseWriter, r *http.Request) {
	util.SuccessResponse(w, http.StatusOK)
}

func ShopifyCustomerDataRequestWebhook(w http.ResponseWriter, r *http.Request) {
	util.SuccessResponse(w, http.StatusOK)
}

func ShopifyShopRedactWebhook(w http.ResponseWriter, r *http.Request) {
	util.SuccessResponse(w, http.StatusOK)
}

// ** END WEBHOOKS HANDLERS ** //
