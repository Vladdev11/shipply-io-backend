package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/robfig/cron"

	"github.com/shipply-io/shipply-io-backend/api"
	SendgridAPI "github.com/shipply-io/shipply-io-backend/api/sendgrid"
	ShipengineAPI "github.com/shipply-io/shipply-io-backend/api/shipengine/api"
	"github.com/shipply-io/shipply-io-backend/handlers"
	"github.com/shipply-io/shipply-io-backend/middlewares"
	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/tasks"
	"github.com/shipply-io/shipply-io-backend/util"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var k = koanf.New(".")

func main() {

	var err error

	dev := flag.Bool("dev", false, "development mode")
	printSQL := flag.Bool("print-sql", false, "log database queries")

	flag.Parse()

	var configPath string
	if *dev {
		configPath = "config.dev.yml"
	} else {
		configPath = "config.prod.yml"
	}

	// TODO: Just gonna do yaml parsing from local file for now
	// We should probably discuss what we want to do for this
	// later with regards to proper CI/CD.
	if err := k.Load(file.Provider(configPath), yaml.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	util.DevelopmentMode = dev
	util.PrintSQL = printSQL

	util.LoadConfig(k)
	// TODO: Get rid of this as much as possible later...
	// Especially the part with the database connection, they should be using context values instead...
	models.PostgresInit()
	api.InitAWSS3()
	SendgridAPI.Init()
	ShipengineAPI.Init()

	// ** START UP FUNCTIONS ** //

	tasks.EnsureSeedData()
	tasks.InitalizeTaskProcesser(1)

	// ** END START UP FUNCTIONS ** //

	//** CRON JOBS **//

	c := cron.New()

	if !*util.DevelopmentMode {
		//TODO add recurring task to pick up tasks that did not start processing
		c.AddFunc("@every 1h", tasks.SyncCarrierConnectionOptions)
		c.AddFunc("@every 1h", tasks.SyncCarrierConnectionPackageTypes)
		c.AddFunc("@every 1h", tasks.SyncCarrierConnectionServices)
	}

	//** END CRON JOBS **//

	router := mux.NewRouter()

	v1 := router.PathPrefix("/v1").Subrouter()

	protected := v1.PathPrefix("/").Subrouter()
	protected.Use(middlewares.AuthMiddleware)

	shopifyRouter := v1.PathPrefix("/shopify").Subrouter()
	shopifyWebhookRouter := shopifyRouter.PathPrefix("/webhooks").Subrouter()
	shopifyWebhookRouter.Use(middlewares.VerifyShopifyWebhook)

	//** SHOPIFY OAUTH **//
	shopifyRouter.HandleFunc("/oauth", handlers.ShopifyOAuth).Methods(http.MethodGet)
	shopifyRouter.HandleFunc("/oauth/confirmation", handlers.ShopifyOAuthConfirmation).Methods(http.MethodGet)
	//** END SHOPIFY OAUTH **//

	//** SHOPIFY WEBHOOKS **//
	shopifyWebhookRouter.HandleFunc("/app-uninstalled", handlers.ShopifyAppUninstalledWebhook).Methods(http.MethodPost)
	shopifyWebhookRouter.HandleFunc("/orders-create", handlers.ShopifyOrderCreatedWebhook).Methods(http.MethodPost)
	shopifyWebhookRouter.HandleFunc("/orders-update", handlers.ShopifyOrderUpdatedWebhook).Methods(http.MethodPost)
	shopifyWebhookRouter.HandleFunc("/products-create", handlers.ShopifyProductCreatedWebhook).Methods(http.MethodPost)
	shopifyWebhookRouter.HandleFunc("/products-update", handlers.ShopifyProductUpdatedWebhook).Methods(http.MethodPost)
	// ** END SHOPIFY WEBHOOKS **//

	// ** SHOPIFY MANDATORY WEBHOOKS **//
	shopifyWebhookRouter.HandleFunc("/customers/data-request", handlers.ShopifyCustomerDataRequestWebhook).Methods(http.MethodPost)
	shopifyWebhookRouter.HandleFunc("/customers/redact", handlers.ShopifyCustomerRedactWebhook).Methods(http.MethodPost)
	shopifyWebhookRouter.HandleFunc("/shop/redact", handlers.ShopifyShopRedactWebhook).Methods(http.MethodPost)
	// ** END SHOPIFY MANDATORY WEBHOOKS **//

	//**  AUTH ROUTES **//
	v1.HandleFunc("/auth/login", handlers.AuthLogin).Methods(http.MethodPost)
	v1.HandleFunc("/auth/login-no-expiration", handlers.AuthNoExpirationToken).Methods(http.MethodPost)
	v1.HandleFunc("/auth/reset-password", handlers.ResetPassword).Methods(http.MethodPost)
	v1.HandleFunc("/auth/validate-reset-password-token", handlers.ValidateResetPasswordToken).Methods(http.MethodGet)
	v1.HandleFunc("/auth/update-password", handlers.UpdatePassword).Methods(http.MethodPost)
	//**  END AUTH ROUTES **//

	//**  CLIENT ROUTES **//
	protected.HandleFunc("/client/list", handlers.ListClients).Methods(http.MethodGet)
	protected.HandleFunc("/client/create", handlers.ClientCreate).Methods(http.MethodPost)
	protected.HandleFunc("/client/{client_id}", handlers.GetClient).Methods(http.MethodGet)
	protected.HandleFunc("/client/{client_id}/update", handlers.UpdateClient).Methods(http.MethodPatch)
	protected.HandleFunc("/client/{client_id}/update-avatar", handlers.UpdateClientAvatar).Methods(http.MethodPost)
	//**  END CLIENT ROUTES **//

	// ** STORE ROUTES **//
	protected.HandleFunc("/store/list", handlers.ListStores).Methods(http.MethodGet)
	protected.HandleFunc("/store/{id}", handlers.GetStore).Methods(http.MethodGet)
	protected.HandleFunc("/store/{id}/activate", handlers.ActivateStore).Methods(http.MethodPatch)
	protected.HandleFunc("/store/{id}/deactivate", handlers.DeactivateStore).Methods(http.MethodPatch)
	protected.HandleFunc("/store/{id}/update", handlers.UpdateStore).Methods(http.MethodPatch)
	protected.HandleFunc("/store/{id}/delete", handlers.DeleteStore).Methods(http.MethodDelete)
	// ** END STORE ROUTES **//

	//**  USER ROUTES **//
	protected.HandleFunc("/user/self", handlers.UserGet).Methods(http.MethodGet)
	protected.HandleFunc("/user/self/update-password", handlers.UserUpdatePassword).Methods(http.MethodPatch)
	protected.HandleFunc("/user/{id}", handlers.UserGetByID).Methods(http.MethodGet)
	protected.HandleFunc("/user/create", handlers.UserCreate).Methods(http.MethodPost)
	protected.HandleFunc("/user/{user_id}/update", handlers.UserUpdate).Methods(http.MethodPatch)
	protected.HandleFunc("/user/{user_id}/update-avatar", handlers.UserUpdateAvatar).Methods(http.MethodPost)
	protected.HandleFunc("/user/{user_id}/delete", handlers.UserDelete).Methods(http.MethodDelete)
	//**  END USER ROUTES **//

	//**  PURCHASE ORDER ROUTES **//
	protected.HandleFunc("/purchase-order/list", handlers.PurchaseOrderList).Methods(http.MethodGet)
	protected.HandleFunc("/purchase-order/create", handlers.PurchaseOrderCreate).Methods(http.MethodPost)
	protected.HandleFunc("/purchase-order/{id}", handlers.PurchaseOrderGet).Methods(http.MethodGet)
	protected.HandleFunc("/purchase-order/{id}/update", handlers.PurchaseOrderUpdate).Methods(http.MethodPatch)
	protected.HandleFunc("/purchase-order/{id}/delete", handlers.PurchaseOrderDelete).Methods(http.MethodDelete)
	//**  END PURCHASE ORDER ROUTES **//

	//**  PURCHASE ORDER ITEM ROUTES **//
	protected.HandleFunc("/purchase-order/{id}/items/bulk-update", handlers.PurchaseOrderItemUpdateBulk).Methods(http.MethodPatch)
	protected.HandleFunc("/purchase-order/{id}/item/create", handlers.PurchaseOrderItemCreate).Methods(http.MethodPost)
	protected.HandleFunc("/purchase-order/{id}/item/{item_id}", handlers.PurchaseOrderItemGet).Methods(http.MethodGet)
	protected.HandleFunc("/purchase-order/{id}/item/{item_id}/update", handlers.PurchaseOrderItemUpdate).Methods(http.MethodPatch)
	protected.HandleFunc("/purchase-order/{id}/item/{item_id}/delete", handlers.PurchaseOrderItemDelete).Methods(http.MethodDelete)
	//** END PURCHASE ORDER ITEM ROUTES **//

	//**  PURCHASE ORDER STATUS ROUTES **//
	protected.HandleFunc("/purchase-order-status/create", handlers.PurchaseOrderStatusCreate).Methods(http.MethodPost)
	protected.HandleFunc("/purchase-order-status/list", handlers.PurchaseOrderStatusList).Methods(http.MethodGet)
	protected.HandleFunc("/purchase-order-status/{id}/update", handlers.PurchaseOrderStatusUpdate).Methods(http.MethodPatch)
	protected.HandleFunc("/purchase-order-status/{id}/delete", handlers.PurchaseOrderStatusDelete).Methods(http.MethodDelete)
	//** END PURCHASE ORDER STATUS ROUTES **//

	//**  PURCHASE ORDER NOTE ROUTES **//
	protected.HandleFunc("/purchase-order/{id}/notes/create", handlers.PurchaseOrderHistoryCreate).Methods(http.MethodPost)
	//** END PURCHASE ORDER NOTE ROUTES **//

	//** PURCHASE ORDER ATTACHMENT ROUTES **//
	protected.HandleFunc("/purchase-order/{id}/attachment/create", handlers.PurchaseOrderAttachmentCreate).Methods(http.MethodPost)
	protected.HandleFunc("/purchase-order/{id}/attachment/{purchase_order_attachment_id}/delete", handlers.PurchaseOrderAttachmentDelete).Methods(http.MethodDelete)
	protected.HandleFunc("/purchase-order/{id}/attachment/list", handlers.PurchaseOrderAttachmentList).Methods(http.MethodGet)
	//** END PURCHASE ORDER ATTACHMENT ROUTES **//

	//** RECEIVING ROUTES **//
	protected.HandleFunc("/purchase-order/{purchase_order_id}/receiving/list-items", handlers.ReceivingListItems).Methods(http.MethodGet)
	protected.HandleFunc("/purchase-order/{purchase_order_id}/receiving/batch-receive", handlers.PurchaseOrderItemReceive).Methods(http.MethodPost)
	protected.HandleFunc("/purchase-order/{purchase_order_id}/receiving/{item_id}", handlers.PurchaseOrderItemGetReceivingDetails).Methods(http.MethodGet)
	protected.HandleFunc("/purchase-order/{purchase_order_id}/receiving/{item_id}/reject", handlers.PurchaseOrderItemReject).Methods(http.MethodPost)
	protected.HandleFunc("/purchase-order/{purchase_order_id}/receiving/{item_id}/update-ipa", handlers.PurchaseOrderItemUpdateIPAInfo).Methods(http.MethodPatch)
	protected.HandleFunc("/purchase-order/{purchase_order_id}/receiving/scan-input", handlers.PurchaseOrderItemScanInput).Methods(http.MethodPost)
	//** END RECEIVING ROUTES **//

	//** PRODUCT LOT ROUTES **//
	// protected.HandleFunc("/product-lots/list", handlers.ProductLotList).Methods(http.MethodGet)
	protected.HandleFunc("/product-lots/create", handlers.ProductLotCreate).Methods(http.MethodPost)
	protected.HandleFunc("/product-lots/{product_id}/list", handlers.ProductLotListByProduct).Methods(http.MethodGet)
	// protected.HandleFunc("/product-lots/{product_id)/{product_lot_id}/update", handlers.ProductLotUpdate).Methods(http.MethodPatch)
	//** END PRODUCT LOT ROUTES **//

	//** PRODUCT ROUTES **//
	protected.HandleFunc("/product-search", handlers.ProductSearch).Methods(http.MethodGet)
	//** END PRODUCT ROUTES **//

	//** PRODUCT ALIAS ROUTES **//
	protected.HandleFunc("/product-aliases", handlers.ProductAliasCreate).Methods(http.MethodPost)
	protected.HandleFunc("/product-aliases/{barcode}", handlers.ProductAliasGetByBarcode).Methods(http.MethodGet)
	protected.HandleFunc("/product-aliases/{barcode}", handlers.ProductAliasUpdateByBarcode).Methods(http.MethodPatch)
	protected.HandleFunc("/product-aliases/{barcode}", handlers.ProductAliasDeleteByBarcode).Methods(http.MethodDelete)
	//** END PRODUCT ALIAS ROUTES **//

	//** VENDOR ROUTES **//
	protected.HandleFunc("/vendor/list", handlers.VendorList).Methods(http.MethodGet)
	protected.HandleFunc("/vendor/create", handlers.VendorCreate).Methods(http.MethodPost)
	protected.HandleFunc("/vendor/{id}", handlers.VendorGet).Methods(http.MethodGet)
	protected.HandleFunc("/vendor/{id}/update", handlers.VendorUpdate).Methods(http.MethodPatch)
	protected.HandleFunc("/vendor/{id}/delete", handlers.VendorDelete).Methods(http.MethodDelete)
	//** END VENDOR ROUTES **//

	// ** WAREHOUSE ROUTES **//
	protected.HandleFunc("/warehouse/list", handlers.WarehouseList).Methods(http.MethodGet)
	protected.HandleFunc("/warehouse/create", handlers.WarehouseCreate).Methods(http.MethodPost)
	protected.HandleFunc("/warehouse/{id}", handlers.WarehouseGet).Methods(http.MethodGet)
	protected.HandleFunc("/warehouse/{id}/update", handlers.WarehouseUpdate).Methods(http.MethodPatch)
	protected.HandleFunc("/warehouse/{id}/delete", handlers.WarehouseDelete).Methods(http.MethodDelete)
	// ** END WAREHOUSE ROUTES **//

	// ** LOCATION ROUTES **//
	protected.HandleFunc("/location/list", handlers.LocationList).Methods(http.MethodGet)
	protected.HandleFunc("/location/create", handlers.LocationCreate).Methods(http.MethodPost)
	protected.HandleFunc("/location/{id}", handlers.LocationGet).Methods(http.MethodGet)
	protected.HandleFunc("/location/{id}/update", handlers.LocationUpdate).Methods(http.MethodPatch)
	protected.HandleFunc("/location/{id}/delete", handlers.LocationDelete).Methods(http.MethodDelete)
	// ** END LOCATION ROUTES **//

	// ** LOCATION TYPE ROUTES **//
	protected.HandleFunc("/location-type/list", handlers.LocationTypeList).Methods(http.MethodGet)
	protected.HandleFunc("/location-type/create", handlers.LocationTypeCreate).Methods(http.MethodPost)
	protected.HandleFunc("/location-type/{id}", handlers.LocationTypeGet).Methods(http.MethodGet)
	protected.HandleFunc("/location-type/{id}/update", handlers.LocationTypeUpdate).Methods(http.MethodPatch)
	protected.HandleFunc("/location-type/{id}/delete", handlers.LocationTypeDelete).Methods(http.MethodDelete)
	// ** END LOCATION TYPE ROUTES **//

	//** ORDER ROUTES **//
	protected.HandleFunc("/orders/list", handlers.ListOrders).Methods(http.MethodGet)
	protected.HandleFunc("/order/{id}", handlers.GetOrder).Methods(http.MethodGet)
	//** END ORDER ROUTES **//

	//** ORDER ITEM ROUTES **//
	//** END ORDER ITEM ROUTES **//

	//**  BOX ROUTES **//
	protected.HandleFunc("/boxes/list", handlers.BoxList).Methods(http.MethodGet)
	protected.HandleFunc("/box/create", handlers.BoxCreate).Methods(http.MethodPost)
	protected.HandleFunc("/box/{id}", handlers.GetBox).Methods(http.MethodGet)
	protected.HandleFunc("/box/{id}/delete", handlers.DeleteBox).Methods(http.MethodDelete)
	protected.HandleFunc("/box/{id}/update", handlers.UpdateBox).Methods(http.MethodPatch)
	//** END BOX ROUTES **//

	//** CARRIER ROUTES **//
	protected.HandleFunc("/carriers/list", handlers.ListCarriers).Methods(http.MethodGet)
	protected.HandleFunc("/carrier/{carrier_id}/connect", handlers.CreateCarrierConnection).Methods(http.MethodPost)
	//** END CARRIER ROUTES **//

	//** CARRIER CONNECTION ROUTES **//
	protected.HandleFunc("/carrier-connections/list", handlers.ListCarrierConnections).Methods(http.MethodGet)
	protected.HandleFunc("/carrier-connection/{carrier_connection_id}", handlers.GetCarrierConnection).Methods(http.MethodGet)
	// protected.HandleFunc("/carrier-connection/{carrier_connection_id}/update", handlers.UpdateCarrierConnection).Methods(http.MethodPatch)
	protected.HandleFunc("/carrier-connection/{carrier_connection_id}/disconnect", handlers.DisconnectCarrierConnection).Methods(http.MethodDelete)
	//** END CARRIER CONNECTION ROUTES **//

	//** SHIPPING METHOD ROUTES **//
	protected.HandleFunc("/shipping-methods/list", handlers.ListShippingMethods).Methods(http.MethodGet)
	protected.HandleFunc("/shipping-method/{shipping_method_id}", handlers.GetShippingMethod).Methods(http.MethodGet)
	protected.HandleFunc("/shipping-method/{shipping_method_id}", handlers.UpdateShippingMethod).Methods(http.MethodPatch)
	//** END SHIPPING METHOD ROUTES **//

	// ** PICK SESSION ROUTES **//
	protected.HandleFunc("/pick-session/create", handlers.CreatePickSession).Methods(http.MethodPost)
	protected.HandleFunc("/pick-session/active", handlers.GetActivePickSession).Methods(http.MethodGet)
	protected.HandleFunc("/pick-session/select-item", handlers.PickSessionSelectItem).Methods(http.MethodPost)
	protected.HandleFunc("/pick-session/assign-tote", handlers.PickSessionAssignTote).Methods(http.MethodPost)
	protected.HandleFunc("/pick-session/confirm-tote", handlers.PickSessionConfirmTote).Methods(http.MethodPost)
	protected.HandleFunc("/pick-session/pick", handlers.PickSessionPick).Methods(http.MethodPost)
	protected.HandleFunc("/pick-session/complete", handlers.PickSessionComplete).Methods(http.MethodPost)
	// ** END PICK SESSION ROUTES **//

	// ** SHIPPING ROUTES **//
	protected.HandleFunc("/shipping/scan-tote", handlers.ShippingScanTote).Methods(http.MethodPost)
	protected.HandleFunc("/shipping/pick-session-order/{id}", handlers.ShippingGetPickSessionOrder).Methods(http.MethodGet)
	protected.HandleFunc("/shipping/pick-session-order/{id}/shop-rates", handlers.ShippingShopRates).Methods(http.MethodPost)
	protected.HandleFunc("/shipping/pick-session-order/{id}/select-rate", handlers.ShippingSelectRate).Methods(http.MethodPost)
	protected.HandleFunc("/shipping/pick-session-order/{id}/purchase-label", handlers.ShippingPurchaseLabel).Methods(http.MethodPost)
	// protected.HandleFunc("/shipping/pick-session-order/{id}/void-label", handlers.ShippingVoidLabel).Methods(http.MethodPost)
	// ** END SHIPPING ROUTES **//

	fmt.Printf("Server starting on port %s", util.ConfigLocalPort)
	err = http.ListenAndServe(fmt.Sprintf(":%s", util.ConfigLocalPort), gorillaHandlers.CORS(
		gorillaHandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}),
		gorillaHandlers.AllowedHeaders([]string{"Access-Control-Allow-Headers", "Content-Type", "Authorization", "Accept", "Accept-Language", "X-Authorization", "X-API", "X-REAL-IP"}),
		gorillaHandlers.AllowedOrigins([]string{"*"}),
		gorillaHandlers.AllowCredentials(),
	)(router))
	if err != nil {
		panic(err)
	}

}
