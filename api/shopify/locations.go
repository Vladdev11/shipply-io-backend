package shopify

import (
	"context"
	"errors"

	"github.com/shipply-io/shipply-io-backend/api/shopify/gen"
	"github.com/shipply-io/shipply-io-backend/models"
)

func AddLocation(ctx context.Context, shopName string, accessToken string, warehouse models.Warehouse) (string, error) {

	client := NewClient(shopName, accessToken)
	fulfillsOnlineOrders := true

	shopifyLocation, err := client.AddLocation(context.Background(), gen.LocationAddInput{
		Name:                 warehouse.Name,
		FulfillsOnlineOrders: &fulfillsOnlineOrders,
		Address: gen.LocationAddAddressInput{
			Address1:     &warehouse.ShipFromAddress.Street1,
			Address2:     &warehouse.ShipFromAddress.Street2,
			City:         &warehouse.ShipFromAddress.City,
			ProvinceCode: &warehouse.ShipFromAddress.State,
			Zip:          &warehouse.ShipFromAddress.PostalCode,
			CountryCode:  gen.CountryCode(warehouse.ShipFromAddress.Country),
		},
	})
	if err != nil {
		return "", err
	}

	if len(shopifyLocation.LocationAdd.UserErrors) > 0 {
		if shopifyLocation.LocationAdd.UserErrors[0].Message != "You already have a location with this name" {
			return "", errors.New(shopifyLocation.LocationAdd.UserErrors[0].Message)
		}

		location, err := client.GetLocationsByName(context.Background(), warehouse.Name)
		if err != nil {
			return "", err
		}

		if location.Locations.Edges[0] != nil {
			return location.Locations.Edges[0].Node.ID, nil
		} else {
			return "", errors.New("Could not find location")
		}
	}

	models.LogShopifyAPIEvent(ctx, &models.ShopifyApiLog{
		ShopDomain: shopName,
		EventType:  "api_response",
		Endpoint:   "locations",
		Method:     "POST",
	})

	return shopifyLocation.LocationAdd.Location.ID, nil

}
