package shopify

import (
	"context"

	"github.com/shipply-io/shipply-io-backend/api/shopify/gen"
	"github.com/shipply-io/shipply-io-backend/models"
)

func ActivateInventoryItem(shopName string, accessToken string, inventoryItemID string, locationID string) error {

	client := NewClient(shopName, accessToken)

	_, err := client.InventoryActivate(context.Background(), inventoryItemID, locationID)
	if err != nil {
		return err
	}

	models.LogShopifyAPIEvent(&models.ShopifyApiLog{
		ShopDomain: shopName,
		EventType:  "api_response",
		Endpoint:   "activate_inventory_item",
		Method:     "POST",
	})

	return nil

}

func SetOnHandInventoryInShopify(shopName string, accessToken string, shopifyLocationID string, shopifyInventoryItemID string, quantity int) error {

	client := NewClient(shopName, accessToken)

	_, err := client.InventorySetOnHandQuantities(context.Background(), gen.InventorySetOnHandQuantitiesInput{
		SetQuantities: []*gen.InventorySetQuantityInput{
			{
				LocationID:      shopifyLocationID,
				InventoryItemID: shopifyInventoryItemID,
				Quantity:        quantity,
			},
		},
		Reason: "correction",
	})

	if err != nil {
		return err
	}

	models.LogShopifyAPIEvent(&models.ShopifyApiLog{
		ShopDomain: shopName,
		EventType:  "api_response",
		Endpoint:   "set_on_hand_inventory",
		Method:     "POST",
	})

	return nil

}
