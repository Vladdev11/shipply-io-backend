package responses

import "github.com/shipply-io/shipply-io-backend/models"

/* ---------------------------- ListLocationTypes --------------------------- */

// ListLocationTypesResponse represents the expected response body for the ListLocationTypes endpoint
type ListLocationTypesResponse struct {
	LocationTypes []LocationTypeResponseForListLocationTypes `json:"location_types"`
}

// LocationTypeResponseForListLocationTypes represents the expected response body for the ListLocationTypes endpoint
type LocationTypeResponseForListLocationTypes struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GenerateListLocationTypesResponse generates the response body for the ListLocationTypes endpoint
func GenerateListLocationTypesResponse(locationTypes []models.LocationType) *ListLocationTypesResponse {

	var locationTypeResponses []LocationTypeResponseForListLocationTypes
	for _, locationType := range locationTypes {
		locationTypeResponses = append(locationTypeResponses, LocationTypeResponseForListLocationTypes{
			ID:   locationType.ID,
			Name: locationType.Name,
		})
	}

	return &ListLocationTypesResponse{
		LocationTypes: locationTypeResponses,
	}

}

/* ----------------------------- GetLocationType ----------------------------- */

// GetLocationTypeResponse represents the expected response body for the GetLocationType endpoint
type GetLocationTypeResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GenerateGetLocationTypeResponse generates the response body for the GetLocationType endpoint
func GenerateGetLocationTypeResponse(locationType models.LocationType) *GetLocationTypeResponse {

	return &GetLocationTypeResponse{
		ID:   locationType.ID,
		Name: locationType.Name,
	}

}

/* ----------------------------- CreateLocationType ----------------------------- */

// CreateLocationTypeResponse represents the expected response body for the CreateLocationType endpoint
type CreateLocationTypeResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GenerateCreateLocationTypeResponse generates the response body for the CreateLocationType endpoint
func GenerateCreateLocationTypeResponse(locationType models.LocationType) *CreateLocationTypeResponse {

	return &CreateLocationTypeResponse{
		ID:   locationType.ID,
		Name: locationType.Name,
	}

}

/* ----------------------------- UpdateLocationType ----------------------------- */

// UpdateLocationTypeResponse represents the expected response body for the UpdateLocationType endpoint
type UpdateLocationTypeResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GenerateUpdateLocationTypeResponse generates the response body for the UpdateLocationType endpoint
func GenerateUpdateLocationTypeResponse(locationType models.LocationType) *UpdateLocationTypeResponse {

	return &UpdateLocationTypeResponse{
		ID:   locationType.ID,
		Name: locationType.Name,
	}

}
