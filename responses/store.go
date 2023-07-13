package responses

import (
	"encoding/json"

	"github.com/shipply-io/shipply-io-backend/models"
)

/* -------------------------------- GetStore -------------------------------- */

// GetStoreResponse represents the expected response body for the GetStore endpoint
type GetStoreResponse struct {
	ID          int                            `json:"id"`
	Name        string                         `json:"name"`
	Active      bool                           `json:"active"`
	Marketplace MarketplaceResponseForGetStore `json:"marketplace"`
	Settings    json.RawMessage                `json:"settings"`
}

// MarketplaceResponseForGetStore represents the expected response body for an individual marketplace in the GetStore endpoint
type MarketplaceResponseForGetStore struct {
	Name              string `json:"name"`
	SmallThumbnailURL string `json:"small_thumbnail_url"`
}

// GenerateGetStoreResponse generates the response body for the GetStore endpoint
func GenerateGetStoreResponse(store models.Store) *GetStoreResponse {

	marketplaceResponse := MarketplaceResponseForGetStore{
		Name:              store.Marketplace.Name,
		SmallThumbnailURL: store.Marketplace.SmallThumbnailURL,
	}

	return &GetStoreResponse{
		ID:          store.ID,
		Name:        store.Name,
		Active:      store.Active,
		Marketplace: marketplaceResponse,
		Settings:    store.Settings,
	}
}

/* ------------------------------- ListStores ------------------------------- */

// ListStoresResponse represents the expected response body for the ListStores endpoint
type ListStoresResponse struct {
	Stores []StoreResponseForListStores `json:"stores"`
}

// StoreResponseForListStores represents the expected response body for an individual store in the ListStores endpoint
type StoreResponseForListStores struct {
	ID          int                              `json:"id"`
	Name        string                           `json:"name"`
	Active      bool                             `json:"active"`
	Marketplace MarketplaceResponseForListStores `json:"marketplace"`
}

// MarketplaceResponseForListStores represents the expected response body for an individual marketplace in the ListStores endpoint
type MarketplaceResponseForListStores struct {
	Name              string `json:"name"`
	SmallThumbnailURL string `json:"small_thumbnail_url"`
}

// GenerateListStoresResponse generates the response body for the ListStores endpoint
func GenerateListStoresResponse(stores []models.Store) *ListStoresResponse {

	storeResponses := []StoreResponseForListStores{}
	for _, store := range stores {

		marketplaceResponse := MarketplaceResponseForListStores{
			Name:              store.Marketplace.Name,
			SmallThumbnailURL: store.Marketplace.SmallThumbnailURL,
		}

		storeResponses = append(storeResponses, StoreResponseForListStores{
			ID:          store.ID,
			Name:        store.Name,
			Active:      store.Active,
			Marketplace: marketplaceResponse,
		})
	}

	return &ListStoresResponse{
		Stores: storeResponses,
	}
}

/* ------------------------------- UpdateStore ------------------------------ */

// UpdateStoreResponse represents the expected response body for the UpdateStore endpoint
type UpdateStoreResponse struct {
	ID          int                            `json:"id"`
	Name        string                         `json:"name"`
	Active      bool                           `json:"active"`
	Marketplace MarketplaceResponseForGetStore `json:"marketplace"`
	Settings    json.RawMessage                `json:"settings"`
}

// GenerateUpdateStoreResponse generates the response body for the UpdateStore endpoint
func GenerateUpdateStoreResponse(store models.Store) *UpdateStoreResponse {

	marketplaceResponse := MarketplaceResponseForGetStore{
		Name:              store.Marketplace.Name,
		SmallThumbnailURL: store.Marketplace.SmallThumbnailURL,
	}

	return &UpdateStoreResponse{
		ID:          store.ID,
		Name:        store.Name,
		Active:      store.Active,
		Marketplace: marketplaceResponse,
		Settings:    store.Settings,
	}
}
