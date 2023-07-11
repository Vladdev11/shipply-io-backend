package responses

import (
	"fmt"
	"time"

	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/util"
)

/* ----------------------------- SearchProducts ----------------------------- */

// SearchProductsResponse represents the response body for the SearchProducts endpoint
type SearchProductsResponse struct {
	Products []ProductResponseForSearchProducts `json:"products"`
}

// ProductResponseForSearchProducts represents the response body for an individual product for SearchProducts endpoint
type ProductResponseForSearchProducts struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Sku         string  `json:"sku"`
	Available   int     `json:"available"`
	OnHand      int     `json:"on_hand"`
	Allocated   int     `json:"allocated"`
	Backordered int     `json:"backordered"`
	Price       float64 `json:"price"`
}

// GenerateSearchProductsResponse converts a product to a SearchProductsResponse
func GenerateSearchProductsResponse(products []models.Product) *SearchProductsResponse {
	// Initialize a slice to hold the product responses
	productResponses := make([]ProductResponseForSearchProducts, len(products))

	// Iterate over the products and convert each to a ProductResponseForSearchProducts
	for i, product := range products {
		productResponses[i] = ProductResponseForSearchProducts{
			ID:          product.ID,
			Name:        product.Name,
			Sku:         product.Sku,
			Available:   product.Available,
			OnHand:      product.OnHand,
			Allocated:   product.Allocated,
			Backordered: product.Backordered,
			Price:       product.Price,
		}
	}

	// Return the final response
	return &SearchProductsResponse{
		Products: productResponses,
	}
}

/* ------------------------------ ListProducts ------------------------------ */

// ListProductsResponse represents the response body for the ListProducts endpoint
type ListProductsResponse struct {
	TotalCount    int                              `json:"total_count"`
	FilteredCount int                              `json:"filtered_count"`
	Data          []ProductResponseForListProducts `json:"data"`
}

// ProductResponseForListProducts represents the response body for an individual product for ListProducts endpoint
type ProductResponseForListProducts struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Sku         string  `json:"sku"`
	Available   int     `json:"available"`
	OnHand      int     `json:"on_hand"`
	Allocated   int     `json:"allocated"`
	Backordered int     `json:"backordered"`
	Price       float64 `json:"price"`
}

// GenerateListProductsResponse converts a product to a ListProductsResponse
func GenerateListProductsResponse(products []models.Product, count int, total int) *ListProductsResponse {
	// Initialize a slice to hold the product responses
	productResponses := make([]ProductResponseForListProducts, len(products))

	// Iterate over the products and convert each to a ProductResponseForListProducts
	for i, product := range products {
		productResponses[i] = ProductResponseForListProducts{
			ID:          product.ID,
			Name:        product.Name,
			Sku:         product.Sku,
			Available:   product.Available,
			OnHand:      product.OnHand,
			Allocated:   product.Allocated,
			Backordered: product.Backordered,
			Price:       product.Price,
		}
	}

	// Return the final response
	return &ListProductsResponse{
		TotalCount:    total,
		FilteredCount: count,
		Data:          productResponses,
	}
}

/* ------------------------------ GetProduct ------------------------------- */

// GetProductResponse represents the response body for the GetProduct endpoint
type GetProductResponse struct {
	ID                   int       `json:"id"`
	Name                 string    `json:"name"`
	Sku                  string    `json:"sku"`
	Barcode              string    `json:"barcode"`
	Price                float64   `json:"price"`
	Available            int       `json:"available"`
	OnHand               int       `json:"on_hand"`
	Allocated            int       `json:"allocated"`
	Backordered          int       `json:"backordered"`
	ImageURL             string    `json:"image_url"`
	AdditionalImageUrls  []string  `json:"additional_image_urls"`
	ProductNote          string    `json:"product_note"`
	PackNote             string    `json:"pack_note"`
	ReturnNote           string    `json:"return_note"`
	Width                float64   `json:"width"`
	Height               float64   `json:"height"`
	Length               float64   `json:"length"`
	Weight               float64   `json:"weight"`
	Value                float64   `json:"value"`
	CustomsValue         float64   `json:"customs_value"`
	CustomsDescription   string    `json:"customs_description"`
	LastCounted          time.Time `json:"last_counted"`
	CountryOfManufacture string    `json:"country_of_manufacture"`
	TariffCode           string    `json:"tariff_code"`
	NeedsSerialNumber    bool      `json:"needs_serial_number"`
	NeedsLotNumber       bool      `json:"needs_lot_number"`
	IgnoreOnInvoice      bool      `json:"ignore_on_invoice"`
	IgnoreOnCustoms      bool      `json:"ignore_on_customs"`
	IsComponent          bool      `json:"is_component"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// GenerateGetProductResponse converts a product to a GetProductResponse
func GenerateGetProductResponse(product *models.Product) *GetProductResponse {

	// Initialize a slice to hold the additional image urls
	additionalImageUrls := make([]string, len(product.AdditionalImageUrls))

	// Use copy to copy the values from product.AdditionalImageUrls to additionalImageUrls
	copy(additionalImageUrls, product.AdditionalImageUrls)

	//Determine if product is a component
	isComponent := false
	if product.ProductBundles != nil && len(product.ProductBundles) > 0 {
		isComponent = true
	}

	return &GetProductResponse{
		ID:                   product.ID,
		Name:                 product.Name,
		Sku:                  product.Sku,
		Barcode:              product.Barcode,
		Price:                product.Price,
		Available:            product.Available,
		OnHand:               product.OnHand,
		Allocated:            product.Allocated,
		Backordered:          product.Backordered,
		ImageURL:             product.ImageURL,
		AdditionalImageUrls:  additionalImageUrls,
		ProductNote:          product.ProductNote,
		PackNote:             product.PackNote,
		ReturnNote:           product.ReturnNote,
		Width:                product.Width,
		Height:               product.Height,
		Length:               product.Length,
		Weight:               product.Weight,
		Value:                product.Value,
		CustomsValue:         product.CustomsValue,
		CustomsDescription:   product.CustomsDescription,
		LastCounted:          product.LastCounted,
		CountryOfManufacture: product.CountryOfManufacture,
		TariffCode:           product.TariffCode,
		NeedsSerialNumber:    product.NeedsSerialNumber,
		NeedsLotNumber:       product.NeedsLotNumber,
		IgnoreOnInvoice:      product.IgnoreOnInvoice,
		IgnoreOnCustoms:      product.IgnoreOnCustoms,
		IsComponent:          isComponent,
		UpdatedAt:            product.UpdatedAt,
	}
}

/* ---------------------------- GetProductOrders ---------------------------- */

// GetProductOrdersResponse represents the response body for the GetProductOrders endpoint
type GetProductOrdersResponse struct {
	Orders []ProductOrderResponseForGetProductOrders `json:"orders"`
}

// ProductOrderResponseForGetProductOrders represents the response body for an individual product order for GetProductOrders endpoint
type ProductOrderResponseForGetProductOrders struct {
	ID          int       `json:"id"`
	OrderDate   time.Time `json:"order_date"`
	Status      string    `json:"status"`
	OrderNumber string    `json:"order_number"`
}

// GenerateGetProductOrdersResponse converts a product to a GetProductOrdersResponse
func GenerateGetProductOrdersResponse(orders []models.Order) *GetProductOrdersResponse {
	// Initialize a slice to hold the product responses
	productOrderResponses := make([]ProductOrderResponseForGetProductOrders, len(orders))

	// Iterate over the products and convert each to a ProductResponseForListProducts
	for i, order := range orders {
		productOrderResponses[i] = ProductOrderResponseForGetProductOrders{
			ID:          order.ID,
			OrderDate:   order.OrderDate,
			Status:      order.Status,
			OrderNumber: order.OrderNumber,
		}
	}

	// Return the final response
	return &GetProductOrdersResponse{
		Orders: productOrderResponses,
	}
}

/* --------------------------- GetProductInventory -------------------------- */

// GetProductInventoryResponse represents the response body for the GetProductInventory endpoint
type GetProductInventoryResponse struct {
	InventoryLocations []ProductInventoryLocationResponseForGetProductInventory `json:"inventory_locations"`
	Aliases            []ProductAliasResponseForGetProductInventory             `json:"aliases"`
	InventoryHistory   []ProductInventoryHistoryResponseForGetProductInventory  `json:"inventory_history"`
	ProductLots        []ProductLotResponseForGetProductInventory               `json:"lots"`
}

// ProductAliasResponseForGetProductInventory represents the response body for an individual product alias for GetProductInventory endpoint
type ProductInventoryLocationResponseForGetProductInventory struct {
	Quantity     int    `json:"quantity"`
	LocationName string `json:"location_name"`
	LocationID   int    `json:"location_id"`
	Pickable     bool   `json:"pickable"`
	Sellable     bool   `json:"sellable"`
}

// ProductAliasResponseForGetProductInventory represents the response body for an individual product alias for GetProductInventory endpoint
type ProductAliasResponseForGetProductInventory struct {
	ID       int    `json:"id"`
	Barcode  string `json:"barcode"`
	Quantity int    `json:"quantity"`
}

// ProductInventoryHistoryResponseForGetProductInventory represents the response body for an individual product inventory history for GetProductInventory endpoint
type ProductInventoryHistoryResponseForGetProductInventory struct {
	ID            int                  `json:"id"`
	LocationID    int                  `json:"location_id"`
	ProductID     int                  `json:"product_id"`
	Delta         int                  `json:"delta"`
	Note          string               `json:"note"`
	ChangedByUser ChangedByUserHistory `json:"changed_by_user"`
}

// ProductLotResponseForGetProductInventory represents the response body for an individual product lot for GetProductInventory endpoint
type ProductLotResponseForGetProductInventory struct {
	ID         int       `json:"id"`
	LotNumber  string    `json:"lot_number"`
	ExpiryDate time.Time `json:"expiry_date"`
	Active     bool      `json:"active"`
}

// GenerateGetProductInventoryResponse converts a product to a GetProductInventoryResponse
func GenerateGetProductInventoryResponse(product *models.Product) *GetProductInventoryResponse {
	// Converting each InventoryLocation
	inventoryLocations := make([]ProductInventoryLocationResponseForGetProductInventory, len(product.InventoryLocations))
	for i, loc := range product.InventoryLocations {
		inventoryLocations[i] = ProductInventoryLocationResponseForGetProductInventory{
			Quantity:     loc.Quantity,
			LocationName: loc.LocationName,
			LocationID:   loc.LocationID,
			Pickable:     loc.Pickable,
			Sellable:     loc.Sellable,
		}
	}

	// Converting each Alias
	aliases := make([]ProductAliasResponseForGetProductInventory, len(product.Aliases))
	for i, alias := range product.Aliases {
		aliases[i] = ProductAliasResponseForGetProductInventory{
			ID:       alias.ID,
			Barcode:  alias.Barcode,
			Quantity: alias.Quantity,
		}
	}

	// Converting each InventoryHistory
	inventoryHistories := make([]ProductInventoryHistoryResponseForGetProductInventory, len(product.InventoryHistory))
	for i, history := range product.InventoryHistory {
		inventoryHistories[i] = ProductInventoryHistoryResponseForGetProductInventory{
			ID:         history.ID,
			LocationID: history.LocationID,
			ProductID:  history.ProductID,
			Delta:      history.Delta,
			Note:       history.Note,
			ChangedByUser: ChangedByUserHistory{
				ID:        history.User.ID,
				FirstName: history.User.FirstName,
				LastName:  history.User.LastName,
				ImageURL:  fmt.Sprintf("%s/%s", util.ConfigCDNHost, history.User.AvatarFileName),
			},
		}
	}

	// Converting each ProductLot
	productLots := make([]ProductLotResponseForGetProductInventory, len(product.ProductLots))
	for i, lot := range product.ProductLots {
		productLots[i] = ProductLotResponseForGetProductInventory{
			ID:         lot.ID,
			LotNumber:  lot.LotNumber,
			ExpiryDate: lot.ExpiryDate,
			Active:     lot.Active,
		}
	}

	return &GetProductInventoryResponse{
		InventoryLocations: inventoryLocations,
		Aliases:            aliases,
		InventoryHistory:   inventoryHistories,
		ProductLots:        productLots,
	}
}

/* ---------------------------- GetProductBundles --------------------------- */

// GetProductBundlesResponse represents the response body for the GetProductBundles endpoint
type GetProductBundlesResponse struct {
	Bundles []ProductBundleResponseForGetProductBundles `json:"bundles"`
}

// ProductBundleResponseForGetProductBundles represents the response body for an individual product bundle
type ProductBundleResponseForGetProductBundles struct {
	ID       int                                `json:"id"`
	ClientID int                                `json:"client_id"`
	Product  ProductBundleParentProductResponse `json:"product"`
}

// ProductBundleParentProductResponse represents the response body for an individual product bundle's parent product
type ProductBundleParentProductResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Sku      string `json:"sku"`
	Barcode  string `json:"barcode"`
	ImageURL string `json:"image_url"`
}

// GenerateGetProductBundlesResponse converts a product to a GetProductBundlesResponse
func GenerateGetProductBundlesResponse(product *models.Product) *GetProductBundlesResponse {
	bundleResponses := make([]ProductBundleResponseForGetProductBundles, len(product.ProductBundles))

	for i, bundle := range product.ProductBundles {
		productResponse := ProductBundleParentProductResponse{
			ID:       bundle.Product.ID,
			Name:     bundle.Product.Name,
			Sku:      bundle.Product.Sku,
			Barcode:  bundle.Product.Barcode,
			ImageURL: bundle.Product.ImageURL,
		}

		bundleResponses[i] = ProductBundleResponseForGetProductBundles{
			ID:       bundle.ID,
			ClientID: bundle.ClientID,
			Product:  productResponse,
		}
	}

	return &GetProductBundlesResponse{
		Bundles: bundleResponses,
	}
}

/* ----------------------- GetProductBundleComponents ----------------------- */

// GetProductBundleComponentsResponse represents the response body for the GetProductBundleComponents endpoint
type GetProductBundleComponentsResponse struct {
	Components []ProductBundleComponentResponseForGetProductBundleComponents `json:"components"`
}

// ProductBundleComponentResponseForGetProductBundleComponents represents the response body for an individual product bundle component
type ProductBundleComponentResponseForGetProductBundleComponents struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	OnHand  int    `json:"on_hand"`
	Sku     string `json:"sku"`
	Barcode string `json:"barcode"`
}

// GenerateGetProductBundleComponentsResponse converts a product to a GetProductBundleComponentsResponse
func GenerateGetProductBundleComponentsResponse(components []models.Product) *GetProductBundleComponentsResponse {
	componentResponses := make([]ProductBundleComponentResponseForGetProductBundleComponents, len(components))

	for i, component := range components {
		componentResponses[i] = ProductBundleComponentResponseForGetProductBundleComponents{
			ID:      component.ID,
			Name:    component.Name,
			OnHand:  component.OnHand,
			Sku:     component.Sku,
			Barcode: component.Barcode,
		}
	}

	return &GetProductBundleComponentsResponse{
		Components: componentResponses,
	}
}

/* ---------------------------- GetProductStores ---------------------------- */

// GetProductStoresResponse represents the response body for the GetProductStores endpoint
type GetProductStoresResponse struct {
	Stores []ProductStoreResponseForGetProductStores `json:"stores"`
}

// ProductStoreResponseForGetProductStores represents the response body for an individual product store
type ProductStoreResponseForGetProductStores struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GenerateGetProductStoresResponse converts a product to a GetProductStoresResponse
func GenerateGetProductStoresResponse(product *models.Product) *GetProductStoresResponse {
	storeResponses := make([]ProductStoreResponseForGetProductStores, len(product.Stores))

	for i, store := range product.Stores {
		storeResponses[i] = ProductStoreResponseForGetProductStores{
			ID:   store.ID,
			Name: store.Name,
		}
	}

	return &GetProductStoresResponse{
		Stores: storeResponses,
	}
}
