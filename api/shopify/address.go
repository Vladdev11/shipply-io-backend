package shopify

import (
	"github.com/shipply-io/shipply-io-backend/api/shopify/gen"
	"github.com/shipply-io/shipply-io-backend/models"
)

func mapShopifyAddressToModelAddress(shopifyAddress *gen.AddressInfo) models.Address {
	modelAddress := models.Address{}

	if shopifyAddress.FirstName != nil {
		modelAddress.FirstName = *shopifyAddress.FirstName
	}

	if shopifyAddress.LastName != nil {
		modelAddress.LastName = *shopifyAddress.LastName
	}

	if shopifyAddress.Address1 != nil {
		modelAddress.Street1 = *shopifyAddress.Address1
	}

	if shopifyAddress.Address2 != nil {
		modelAddress.Street2 = *shopifyAddress.Address2
	}

	if shopifyAddress.City != nil {
		modelAddress.City = *shopifyAddress.City
	}

	if shopifyAddress.Province != nil {
		modelAddress.State = *shopifyAddress.Province
	}

	if shopifyAddress.Zip != nil {
		modelAddress.PostalCode = *shopifyAddress.Zip
	}

	if shopifyAddress.Country != nil {
		modelAddress.Country = *shopifyAddress.Country
	}

	if shopifyAddress.Phone != nil {
		modelAddress.Phone = *shopifyAddress.Phone
	}

	return modelAddress
}
