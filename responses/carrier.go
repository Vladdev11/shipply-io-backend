package responses

import "github.com/shipply-io/shipply-io-backend/models"

/* ------------------------------ ListCarriers ------------------------------ */

// ListCarriersResponse represents the expected response body for the ListCarriers endpoint
type ListCarriersResponse struct {
	Carriers []CarrierResponseForListCarriers `json:"carriers"`
}

// CarrierResponseForListCarriers represents the response body for each carrier in the ListCarriers endpoint
type CarrierResponseForListCarriers struct {
	ID                int               `json:"id"`
	Name              string            `json:"name"`
	ThumbnailURL      string            `json:"thumbnail_url"`
	SmallThumbnailURL string            `json:"small_thumbnail_url"`
	RequiredFields    []FieldProperties `json:"required_fields"`
}

// GenerateListCarriersResponse generates the response body for the ListCarriers endpoint
func GenerateListCarriersResponse(carriers []models.Carrier) *ListCarriersResponse {

	carrierResponses := make([]CarrierResponseForListCarriers, len(carriers))

	for i, carrier := range carriers {

		requiredFields := make([]FieldProperties, len(carrier.RequiredFields))
		for i, field := range carrier.RequiredFields {
			requiredFields[i] = FieldProperties{
				Name: field.Name,
				Type: field.Type,
			}
		}

		carrierResponses[i] = CarrierResponseForListCarriers{
			ID:                carrier.ID,
			Name:              carrier.Name,
			ThumbnailURL:      carrier.ThumbnailURL,
			SmallThumbnailURL: carrier.SmallThumbnailURL,
			RequiredFields:    requiredFields,
		}
	}

	return &ListCarriersResponse{
		Carriers: carrierResponses,
	}

}
