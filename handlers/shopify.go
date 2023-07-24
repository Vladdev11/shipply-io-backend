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
)

// ** START INSTALLATION HANDLERS ** //

func ShopifyOAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// TODO add error messages that are friendly to the user

	err := models.LogShopifyInstallRequest(r)
	if err != nil {
		models.CreateSystemError(ctx, err.Error())
	}

	shop, err := util.GetStringQueryParam(r, "shop")
	if err != nil {
		models.CreateSystemError(ctx, "shopify install failed: missing shop parameter")
		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify?error=%s", util.FrontendBaseURLFromContext(ctx), url.QueryEscape("shop parameter is required")), http.StatusFound)
		return
	}

	clientID, err := util.GetIntQueryParam(r, "client_id")
	if err != nil {
		models.CreateSystemError(ctx, "shopify install failed: missing client_id parameter")
		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify?error=%s&shop_name=%s", util.FrontendBaseURLFromContext(ctx), url.QueryEscape("client_id parameter is required"), shop), http.StatusFound)
		return
	}

	nonce := uuid.New()

	existingStore, _ := models.GetShopifyStoreByShopName(ctx, shop)
	if existingStore != nil {

		shopifyCredentials, err := existingStore.GetShopifyCredentials(ctx)
		if err != nil {

			models.SystemError{
				Message:  err.Error(),
				ClientID: clientID,
			}.Create(ctx)

			util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify?error=%s&shop_name=%s", util.FrontendBaseURLFromContext(ctx), url.QueryEscape("failed to get shopify credentials"), shop), http.StatusFound)
			return
		}

		if shopifyCredentials.AccessToken != "" {
			util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify?error=%s&shop_name=%s", util.FrontendBaseURLFromContext(ctx), url.QueryEscape("shopify store "+shop+" already associated with a client"), shop), http.StatusFound)
			return
		}

		err = existingStore.UpdateAPICredentials(ctx, &models.ShopifyCredentials{
			ShopName:            shop,
			ShopifyInstallNonce: nonce.String(),
		})
		if err != nil {

			models.SystemError{
				Message:  err.Error(),
				ClientID: clientID,
			}.Create(ctx)

			util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify?error=%s&shop_name=%s", util.FrontendBaseURLFromContext(ctx), url.QueryEscape("failed to update API Credentials"), shop), http.StatusFound)
			return
		}

		err = existingStore.Activate(ctx)
		if err != nil {

			models.SystemError{
				Message:  err.Error(),
				ClientID: clientID,
			}.Create(ctx)

			util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify?error=%s&shop_name=%s", util.FrontendBaseURLFromContext(ctx), url.QueryEscape("failed to activate store"), shop), http.StatusFound)
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
			}.Create(ctx)

			util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify?error=%s&shop_name=%s", util.FrontendBaseURLFromContext(ctx), url.QueryEscape("failed to marshal credentials"), shop), http.StatusFound)
			return
		}

		storeName, err := util.GetStringQueryParam(r, "store_name")
		if err != nil && err != util.ErrMissingQueryParam {

			models.SystemError{
				Message:  "Invalid store_name parameter",
				ClientID: clientID,
			}.Create(ctx)

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

		err = store.Create(ctx)
		if err != nil {

			models.SystemError{
				Message:  err.Error(),
				ClientID: clientID,
			}.Create(ctx)

			util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify?error=%s&shop_name=%s", util.FrontendBaseURLFromContext(ctx), url.QueryEscape("failed to create store"), shop), http.StatusFound)
			return
		}

	}

	url := shopify.FromContext(ctx).AuthorizeUrl(shop, nonce.String())

	util.RedirectResponse(w, r, url, http.StatusFound)

}

func ShopifyOAuthConfirmation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := models.LogShopifyInstallRequest(r)
	if err != nil {
		models.CreateSystemError(ctx, err.Error())
	}

	// TODO add error messages that are friendly to the user

	URL, err := url.ParseRequestURI(r.RequestURI)
	if err != nil {
		models.CreateSystemError(ctx, err.Error())
		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s", util.FrontendBaseURLFromContext(ctx), err.Error()), http.StatusPermanentRedirect)
		return
	}

	if ok, err := shopify.FromContext(ctx).VerifyAuthorizationURL(URL); !ok {
		models.CreateSystemError(ctx, err.Error())
		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s", util.FrontendBaseURLFromContext(ctx), err.Error()), http.StatusPermanentRedirect)
		return
	}

	shopName, err := util.GetStringQueryParam(r, "shop")
	if err != nil {
		models.CreateSystemError(ctx, err.Error())
		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s", util.FrontendBaseURLFromContext(ctx), "invalid shop"), http.StatusPermanentRedirect)
		return
	}

	code, err := util.GetStringQueryParam(r, "code")
	if err != nil {
		models.CreateSystemError(ctx, err.Error())
		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontendBaseURLFromContext(ctx), "invalid code", shopName), http.StatusPermanentRedirect)
		return
	}

	store, err := models.GetShopifyStoreByShopName(ctx, shopName)
	if err != nil {
		models.CreateSystemError(ctx, err.Error())
		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontendBaseURLFromContext(ctx), "failed to find store with that shop name", shopName), http.StatusPermanentRedirect)
		return
	}

	shopifyCredentials, err := store.GetShopifyCredentials(ctx)
	if err != nil {

		models.SystemError{
			Message:  err.Error(),
			ClientID: store.ClientID,
		}.Create(ctx)

		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontendBaseURLFromContext(ctx), "failed to get API credentials", shopName), http.StatusPermanentRedirect)
		return
	}

	// Verify the security checks
	err = shopify.FromContext(ctx).VerifyOAuthCallback(r.URL, shopifyCredentials.ShopifyInstallNonce)
	if err != nil {

		models.SystemError{
			Message:  err.Error(),
			ClientID: store.ClientID,
		}.Create(ctx)

		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontendBaseURLFromContext(ctx), err.Error(), shopName), http.StatusPermanentRedirect)
		return
	}

	accessToken, err := shopify.FromContext(ctx).GetAccessToken(ctx, shopName, code)
	if err != nil {

		models.SystemError{
			Message:  err.Error(),
			ClientID: store.ClientID,
		}.Create(ctx)

		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontendBaseURLFromContext(ctx), "failed to get access token", shopName), http.StatusPermanentRedirect)
		return
	}

	shopifyCredentials.AccessToken = accessToken
	err = store.UpdateAPICredentials(ctx, shopifyCredentials)
	if err != nil {

		models.SystemError{
			Message:  err.Error(),
			ClientID: store.ClientID,
		}.Create(ctx)

		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontendBaseURLFromContext(ctx), "failed to update store with access token", shopName), http.StatusPermanentRedirect)
		return
	}

	client, err := models.GetClientByID(ctx, store.ClientID)
	if err != nil {

		models.SystemError{
			Message:  err.Error(),
			ClientID: store.ClientID,
		}.Create(ctx)

		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontendBaseURLFromContext(ctx), "failed to get client", shopName), http.StatusPermanentRedirect)
		return
	}
	warehouses, err := models.GetWarehousesByOrganizationID(ctx, client.OrganizationID)
	if err != nil {

		models.SystemError{
			Message:  err.Error(),
			ClientID: client.ID,
		}.Create(ctx)

		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontendBaseURLFromContext(ctx), "failed to get warehouses", shopName), http.StatusPermanentRedirect)
		return
	}

	// Create Locations in Shopify (Can't do in task because they must complete before we can activate inventory items)
	for _, warehouse := range warehouses {
		err = warehouse.GetShipFromAddress(ctx)
		if err != nil {

			models.SystemError{
				Message:  err.Error(),
				ClientID: client.ID,
			}.Create(ctx)

			util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontendBaseURLFromContext(ctx), "failed to get warehouse address", shopName), http.StatusPermanentRedirect)
			return
		}

		locationGraphqlID, err := shopify.AddLocation(ctx, shopName, accessToken, warehouse)
		if err != nil {

			models.SystemError{
				Message:  "failed to create location: " + err.Error(),
				ClientID: client.ID,
			}.Create(ctx)

			util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontendBaseURLFromContext(ctx), "failed to create location", shopName), http.StatusPermanentRedirect)
			return
		}

		shopifyLocation := models.ShopifyLocation{
			StoreID:           store.ID,
			ShopifyLocationID: locationGraphqlID,
			WarehouseID:       warehouse.ID,
		}
		err = shopifyLocation.Create(ctx)
		if err != nil {

			models.SystemError{
				Message:  "failed to create shopify location: " + err.Error(),
				ClientID: client.ID,
			}.Create(ctx)

			util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontendBaseURLFromContext(ctx), "failed to create shopify location", shopName), http.StatusPermanentRedirect)
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
		}.Create(ctx)

		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontendBaseURLFromContext(ctx), "An error occured while importing products", shopName), http.StatusPermanentRedirect)
		return
	}

	retrieveProductsTask := &models.SystemTask{
		TaskType: util.PullShopifyProductsTaskType,
		Payload:  json.RawMessage(productsTaskParameters),
	}
	err = retrieveProductsTask.Create(ctx)
	if err != nil {

		models.SystemError{
			Message:  err.Error(),
			ClientID: client.ID,
		}.Create(ctx)

		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontendBaseURLFromContext(ctx), "An error occured while importing products", shopName), http.StatusPermanentRedirect)
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
		}.Create(ctx)

		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontendBaseURLFromContext(ctx), "An error occured while importing orders", shopName), http.StatusPermanentRedirect)
		return
	}

	retrieveOrdersTask := &models.SystemTask{
		TaskType: util.PullShopifyOrdersTaskType,
		Payload:  json.RawMessage(ordersTaskParameters),
	}
	err = retrieveOrdersTask.Create(ctx)
	if err != nil {

		models.SystemError{
			Message:  err.Error(),
			ClientID: client.ID,
		}.Create(ctx)

		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontendBaseURLFromContext(ctx), "An error occured while importing orders", shopName), http.StatusPermanentRedirect)
		return
	}
	tasks.TaskChannel <- retrieveOrdersTask

	err = shopify.SubscribeToAppUninstalledWebhook(ctx, shopifyCredentials.ShopName, shopifyCredentials.AccessToken)
	if err != nil {

		models.SystemError{
			Message:  "failed to subscribe to app uninstalled webhook: " + err.Error(),
			ClientID: client.ID,
		}.Create(ctx)

		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontendBaseURLFromContext(ctx), "failed to subscribe to app uninstalled webhook", shopName), http.StatusPermanentRedirect)
		return
	}

	err = shopify.SubscribeToOrdersCreateWebhook(ctx, shopifyCredentials.ShopName, shopifyCredentials.AccessToken)
	if err != nil {

		models.SystemError{
			Message:  "failed to subscribe to orders create webhook: " + err.Error(),
			ClientID: client.ID,
		}.Create(ctx)

		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontendBaseURLFromContext(ctx), "failed to subscribe to orders create webhook", shopName), http.StatusPermanentRedirect)
		return
	}

	err = shopify.SubscribeToOrdersUpdateWebhook(ctx, shopifyCredentials.ShopName, shopifyCredentials.AccessToken)
	if err != nil {

		models.SystemError{
			Message:  "failed to subscribe to orders update webhook: " + err.Error(),
			ClientID: client.ID,
		}.Create(ctx)

		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontendBaseURLFromContext(ctx), "failed to subscribe to orders update webhook", shopName), http.StatusPermanentRedirect)
		return
	}

	// Product create webhook is for a new product with new variants
	err = shopify.SubscribeToProductsCreateWebhook(ctx, shopifyCredentials.ShopName, shopifyCredentials.AccessToken)
	if err != nil {

		models.SystemError{
			Message:  "failed to subscribe to products create webhook: " + err.Error(),
			ClientID: client.ID,
		}.Create(ctx)

		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontendBaseURLFromContext(ctx), "failed to subscribe to products create webhook", shopName), http.StatusPermanentRedirect)
		return
	}

	// Product update webhook is for an existing product with new variants
	err = shopify.SubscribeToProductsUpdateWebhook(ctx, shopifyCredentials.ShopName, shopifyCredentials.AccessToken)
	if err != nil {

		models.SystemError{
			Message:  "failed to subscribe to products update webhook: " + err.Error(),
			ClientID: client.ID,
		}.Create(ctx)

		util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/connect/shopify/error?error=%s&shop_name=%s", util.FrontendBaseURLFromContext(ctx), "failed to subscribe to products update webhook", shopName), http.StatusPermanentRedirect)
		return
	}

	util.RedirectResponse(w, r, fmt.Sprintf("%s/stores/%d", util.FrontendBaseURLFromContext(ctx), store.ID), http.StatusPermanentRedirect)
}

// ** END INSTALLATION HANDLERS ** //

// ** BEGIN WEBHOOK HANDLERS ** //
func ShopifyAppUninstalledWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	shopName := r.Header.Get("X-Shopify-Shop-Domain")
	if shopName == "" {
		util.ErrResponse(w, ErrInvalidShopifyShopName, http.StatusBadRequest)
		return
	}

	store, err := models.GetShopifyStoreByShopName(ctx, shopName)
	if err != nil {
		util.ErrResponse(w, ErrNoStoreWithShopifyShopName, http.StatusInternalServerError)
		return
	}

	err = store.DeleteAccessToken(ctx)
	if err != nil {
		util.ErrResponse(w, ErrDeleteShopifyAccessToken, http.StatusInternalServerError)
		return
	}

	err = store.DeleteShopifyLocations(ctx)
	if err != nil {
		util.ErrResponse(w, ErrDeleteShopifyLocations, http.StatusInternalServerError)
		return
	}

	err = store.Deactivate(ctx)
	if err != nil {
		util.ErrResponse(w, ErrShopifyDeactivateStore, http.StatusInternalServerError)
		return
	}

	util.SuccessResponse(w, http.StatusOK)
}

func ShopifyOrderCreatedWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	shopName := r.Header.Get("X-Shopify-Shop-Domain")
	if shopName == "" {
		util.ErrResponse(w, ErrInvalidShopifyShopName, http.StatusBadRequest)
		return
	}

	store, err := models.GetShopifyStoreByShopName(ctx, shopName)
	if err != nil {
		util.ErrResponse(w, ErrNoStoreWithShopifyShopName, http.StatusInternalServerError)
		return
	}

	orderID, errors := shopify.ParseAndValidateOrderCreateRequest(r)
	if errors != nil {
		util.ErrResponse(w, errors, http.StatusBadRequest)
		return
	}

	order, err := shopify.RetrieveOrder(store.ShopifyShopName(ctx), store.ShopifyAccessToken(ctx), *orderID)
	if err != nil {
		util.ErrResponse(w, ErrRetrieveShopifyGraphqlOrder, http.StatusInternalServerError)
		return
	}

	gqlOrder, err := shopify.ConvertGetOrderById_OrderToShopifyModelOrder(*order)
	if err != nil {
		util.ErrResponse(w, ErrConvertShopifyGraphqlOrder, http.StatusInternalServerError)
		return
	}

	err = tasks.SyncShopifyGraphqlOrder(ctx, store.ID, store.ShopifyShopName(ctx), store.ShopifyAccessToken(ctx), gqlOrder)
	if err != nil {
		util.ErrResponse(w, ErrShopifySyncGraphqlOrder, http.StatusInternalServerError)
		return
	}

	util.SuccessResponse(w, http.StatusOK)
}

func ShopifyOrderUpdatedWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	shopName := r.Header.Get("X-Shopify-Shop-Domain")
	if shopName == "" {
		util.ErrResponse(w, ErrInvalidShopifyShopName, http.StatusBadRequest)
		return
	}

	store, err := models.GetShopifyStoreByShopName(ctx, shopName)
	if err != nil {
		util.ErrResponse(w, ErrNoStoreWithShopifyShopName, http.StatusInternalServerError)
		return
	}

	orderID, errors := shopify.ParseAndValidateOrderCreateRequest(r)
	if errors != nil {
		util.ErrResponse(w, errors, http.StatusBadRequest)
		return
	}

	order, err := shopify.RetrieveOrder(store.ShopifyShopName(ctx), store.ShopifyAccessToken(ctx), *orderID)
	if err != nil {
		util.ErrResponse(w, ErrRetrieveShopifyGraphqlOrder, http.StatusInternalServerError)
		return
	}

	gqlOrder, err := shopify.ConvertGetOrderById_OrderToShopifyModelOrder(*order)
	if err != nil {
		util.ErrResponse(w, ErrConvertShopifyGraphqlOrder, http.StatusInternalServerError)
		return
	}

	err = tasks.SyncShopifyGraphqlOrder(ctx, store.ID, store.ShopifyShopName(ctx), store.ShopifyAccessToken(ctx), gqlOrder)
	if err != nil {
		util.ErrResponse(w, ErrShopifySyncGraphqlOrder, http.StatusInternalServerError)
		return
	}

	util.SuccessResponse(w, http.StatusOK)
}

func ShopifyProductCreatedWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	shopName := r.Header.Get("X-Shopify-Shop-Domain")
	if shopName == "" {
		util.ErrResponse(w, ErrInvalidShopifyShopName, http.StatusBadRequest)
		return
	}

	store, err := models.GetShopifyStoreByShopName(ctx, shopName)
	if err != nil {
		util.ErrResponse(w, ErrInvalidShopifyShopName, http.StatusInternalServerError)
		return
	}

	if !store.Active {
		util.SuccessResponse(w, http.StatusOK)
		return
	}

	productID, errors := shopify.ParseAndValidateProductCreateRequest(r)
	if errors != nil {
		util.ErrResponse(w, errors, http.StatusBadRequest)
		return
	}

	productVariants, err := shopify.RetrieveProductVariantsByProductID(store.ShopifyShopName(ctx), store.ShopifyAccessToken(ctx), *productID)
	if err != nil {
		util.ErrResponse(w, ErrRetrieveShopifyGraphqlProduct, http.StatusInternalServerError)
		return
	}

	for _, productVariant := range productVariants {
		gqlProduct, err := shopify.ConvertGetProductById_ProductToShopifyModelProduct(productVariant)
		if err != nil {
			models.SystemError{
				Message:  "failed to convert shopify product to model product: " + err.Error(),
				ClientID: store.ClientID,
			}.Create(ctx)
			continue
		}
		err = tasks.SyncShopifyGraphqlProduct(ctx, store.ID, store.ShopifyShopName(ctx), store.ShopifyAccessToken(ctx), *gqlProduct)
		if err != nil {
			models.SystemError{
				Message:  "failed to sync shopify product: " + err.Error(),
				ClientID: store.ClientID,
			}.Create(ctx)
			continue
		}

	}

	util.SuccessResponse(w, http.StatusOK)

}

func ShopifyProductUpdatedWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	shopName := r.Header.Get("X-Shopify-Shop-Domain")
	if shopName == "" {
		util.ErrResponse(w, ErrInvalidShopifyShopName, http.StatusBadRequest)
		return
	}

	store, err := models.GetShopifyStoreByShopName(ctx, shopName)
	if err != nil {
		util.ErrResponse(w, ErrNoStoreWithShopifyShopName, http.StatusInternalServerError)
		return
	}

	if !store.Active {
		util.SuccessResponse(w, http.StatusOK)
		return
	}

	productID, errors := shopify.ParseAndValidateProductCreateRequest(r)
	if errors != nil {
		util.ErrResponse(w, errors, http.StatusBadRequest)
		return
	}

	productVariants, err := shopify.RetrieveProductVariantsByProductID(store.ShopifyShopName(ctx), store.ShopifyAccessToken(ctx), *productID)
	if err != nil {
		util.ErrResponse(w, ErrRetrieveShopifyGraphqlProduct, http.StatusInternalServerError)
		return
	}

	for _, productVariant := range productVariants {
		gqlProduct, err := shopify.ConvertGetProductById_ProductToShopifyModelProduct(productVariant)
		if err != nil {
			models.SystemError{
				Message:  "failed to convert shopify product to model product: " + err.Error(),
				ClientID: store.ClientID,
			}.Create(ctx)
			continue
		}
		err = tasks.SyncShopifyGraphqlProduct(ctx, store.ID, store.ShopifyShopName(ctx), store.ShopifyAccessToken(ctx), *gqlProduct)
		if err != nil {
			models.SystemError{
				Message:  "failed to sync shopify product: " + err.Error(),
				ClientID: store.ClientID,
			}.Create(ctx)
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
