package responses

import "github.com/shipply-io/shipply-io-backend/models"

/* ------------------------------ ListLocations ----------------------------- */

// ListLocationsResponse represents the expected response body for the ListLocations endpoint
type ListLocationsResponse struct {
	TotalCount    int                                `json:"total_count"`
	FilteredCount int                                `json:"filtered_count"`
	Data          []LocationResponseForListLocations `json:"data"`
}

// LocationResponseForListLocations represents the expected response body for the ListLocations endpoint
type LocationResponseForListLocations struct {
	ID           int                                  `json:"id"`
	Name         string                               `json:"name"`
	WarehouseID  int                                  `json:"warehouse_id"`
	Pickable     bool                                 `json:"pickable"`
	Sellable     bool                                 `json:"sellable"`
	IsTote       bool                                 `json:"is_tote"`
	LocationType LocationTypeResponseForListLocations `json:"location_type"`
}

// LocationTypeResponseForListLocations represents the expected response body for the ListLocations endpoint
type LocationTypeResponseForListLocations struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GenerateListLocationsResponse generates the response body for the ListLocations endpoint
func GenerateListLocationsResponse(locations []models.Location, totalCount int, filteredCount int) *ListLocationsResponse {

	var locationResponses []LocationResponseForListLocations
	for _, location := range locations {
		locationResponses = append(locationResponses, LocationResponseForListLocations{
			ID:          location.ID,
			Name:        location.Name,
			WarehouseID: location.WarehouseID,
			Pickable:    location.Pickable,
			Sellable:    location.Sellable,
			IsTote:      location.IsTote,
			LocationType: LocationTypeResponseForListLocations{
				ID:   location.LocationType.ID,
				Name: location.LocationType.Name,
			},
		})
	}

	return &ListLocationsResponse{
		TotalCount:    totalCount,
		FilteredCount: filteredCount,
		Data:          locationResponses,
	}

}

/* ------------------------------- GetLocation ------------------------------- */

// GetLocationResponse represents the expected response body for the GetLocation endpoint
type GetLocationResponse struct {
	ID           int                                `json:"id"`
	Name         string                             `json:"name"`
	WarehouseID  int                                `json:"warehouse_id"`
	Pickable     bool                               `json:"pickable"`
	Sellable     bool                               `json:"sellable"`
	IsTote       bool                               `json:"is_tote"`
	LocationType LocationTypeResponseForGetLocation `json:"location_type"`
}

// LocationTypeResponseForGetLocation represents the expected response body for the GetLocation endpoint
type LocationTypeResponseForGetLocation struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GenerateGetLocationResponse generates the response body for the GetLocation endpoint
func GenerateGetLocationResponse(location models.Location) *GetLocationResponse {

	return &GetLocationResponse{
		ID:          location.ID,
		Name:        location.Name,
		WarehouseID: location.WarehouseID,
		Pickable:    location.Pickable,
		Sellable:    location.Sellable,
		IsTote:      location.IsTote,
		LocationType: LocationTypeResponseForGetLocation{
			ID:   location.LocationType.ID,
			Name: location.LocationType.Name,
		},
	}

}

/* ------------------------------- CreateLocation ------------------------------- */

// CreateLocationResponse represents the expected response body for the CreateLocation endpoint
type CreateLocationResponse struct {
	ID           int                                   `json:"id"`
	Name         string                                `json:"name"`
	WarehouseID  int                                   `json:"warehouse_id"`
	Pickable     bool                                  `json:"pickable"`
	Sellable     bool                                  `json:"sellable"`
	IsTote       bool                                  `json:"is_tote"`
	LocationType LocationTypeResponseForCreateLocation `json:"location_type"`
}

// LocationTypeResponseForCreateLocation represents the expected response body for the CreateLocation endpoint
type LocationTypeResponseForCreateLocation struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GenerateCreateLocationResponse generates the response body for the CreateLocation endpoint
func GenerateCreateLocationResponse(location models.Location) *CreateLocationResponse {

	return &CreateLocationResponse{
		ID:          location.ID,
		Name:        location.Name,
		WarehouseID: location.WarehouseID,
		Pickable:    location.Pickable,
		Sellable:    location.Sellable,
		IsTote:      location.IsTote,
		LocationType: LocationTypeResponseForCreateLocation{
			ID:   location.LocationType.ID,
			Name: location.LocationType.Name,
		},
	}

}

/* ------------------------------- UpdateLocation ------------------------------- */

// UpdateLocationResponse represents the expected response body for the UpdateLocation endpoint
type UpdateLocationResponse struct {
	ID           int                                   `json:"id"`
	Name         string                                `json:"name"`
	WarehouseID  int                                   `json:"warehouse_id"`
	Pickable     bool                                  `json:"pickable"`
	Sellable     bool                                  `json:"sellable"`
	IsTote       bool                                  `json:"is_tote"`
	LocationType LocationTypeResponseForUpdateLocation `json:"location_type"`
}

// LocationTypeResponseForUpdateLocation represents the expected response body for the UpdateLocation endpoint
type LocationTypeResponseForUpdateLocation struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GenerateUpdateLocationResponse generates the response body for the UpdateLocation endpoint
func GenerateUpdateLocationResponse(location models.Location) *UpdateLocationResponse {

	return &UpdateLocationResponse{
		ID:          location.ID,
		Name:        location.Name,
		WarehouseID: location.WarehouseID,
		Pickable:    location.Pickable,
		Sellable:    location.Sellable,
		IsTote:      location.IsTote,
		LocationType: LocationTypeResponseForUpdateLocation{
			ID:   location.LocationType.ID,
			Name: location.LocationType.Name,
		},
	}

}
