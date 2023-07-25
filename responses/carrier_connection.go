package responses

import (
	"github.com/shipply-io/shipply-io-backend/models"
)

/* ------------------------- ListCarrierConnections ------------------------- */
// ListCarrierConnectionsResponse represents the response body for the ListCarrierConnections endpoint
type ListCarrierConnectionsResponse struct {
	CarrierConnections []CarrierConnectionResponseForListCarrierConnections `json:"carrier_connections"`
}

// CarrierResponseForListCarrierConnections represents the response body for each carrier in the ListCarrierConnections endpoint
type CarrierConnectionResponseForListCarrierConnections struct {
	ID                  int    `json:"id"`
	Nickname            string `json:"nickname"`
	Active              bool   `json:"active"`
	CarrierThumbnailURL string `json:"carrier_thumbnail_url"`
}

// GenerateListCarrierConnectionsResponse generates the response body for the ListCarrierConnections endpoint
func GenerateListCarrierConnectionsResponse(carrierConnections []models.CarrierConnection) *ListCarrierConnectionsResponse {

	carrierConnectionResponses := make([]CarrierConnectionResponseForListCarrierConnections, len(carrierConnections))

	for i, carrierConnection := range carrierConnections {
		carrierConnectionResponses[i] = CarrierConnectionResponseForListCarrierConnections{
			ID:                  carrierConnection.ID,
			Nickname:            carrierConnection.ShipengineNickname,
			Active:              carrierConnection.Active,
			CarrierThumbnailURL: carrierConnection.Carrier.ThumbnailURL,
		}
	}

	return &ListCarrierConnectionsResponse{
		CarrierConnections: carrierConnectionResponses,
	}

}

/* --------------------------- GetCarrierConnection ------------------------- */
// GetCarrierConnectionResponse represents the response body for the GetCarrierConnection endpoint
type GetCarrierConnectionResponse struct {
	ID       int                                     `json:"id"`
	Nickname string                                  `json:"nickname"`
	Active   bool                                    `json:"active"`
	Carrier  CarrierResponseForGetCarrierConnection  `json:"carrier"`
	Settings SettingsResponseForGetCarrierConnection `json:"settings"`
}

// SettingsResponseForGetCarrierConnection represents the response body for the settings in the GetCarrierConnection endpoint
type SettingsResponseForGetCarrierConnection struct {
	GeneralSettings         struct{}                                        `json:"general_settings"`
	CarrierSpecificSettings struct{}                                        `json:"carrier_specific_settings"`
	EnabledServices         []EnabledServiceResponseForGetCarrierConnection `json:"enabled_services"`
}

// EnabledServiceResponseForGetCarrierConnection represents the response body for each enabled service in the GetCarrierConnection endpoint
type EnabledServiceResponseForGetCarrierConnection struct {
	ServiceCode string `json:"service_code"`
	ServiceName string `json:"service_name"`
	Enabled     bool   `json:"enabled"`
}

// CarrierResponseForGetCarrierConnection represents the response body for the carrier in the GetCarrierConnection endpoint
type CarrierResponseForGetCarrierConnection struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	ThumbnailURL      string `json:"thumbnail_url"`
	SmallThumbnailURL string `json:"small_thumbnail_url"`
}

// GenerateGetCarrierConnectionResponse generates the response body for the GetCarrierConnection endpoint
func GenerateGetCarrierConnectionResponse(carrierConnection models.CarrierConnection) *GetCarrierConnectionResponse {

	enabledServices := make([]EnabledServiceResponseForGetCarrierConnection, len(carrierConnection.ParsedSettings.EnabledServices))

	for i, enabledService := range carrierConnection.ParsedSettings.EnabledServices {
		enabledServices[i] = EnabledServiceResponseForGetCarrierConnection{
			ServiceCode: enabledService.ServiceCode,
			ServiceName: enabledService.ServiceName,
			Enabled:     enabledService.Enabled,
		}
	}

	return &GetCarrierConnectionResponse{
		ID:       carrierConnection.ID,
		Nickname: carrierConnection.ShipengineNickname,
		Active:   carrierConnection.Active,
		Carrier: CarrierResponseForGetCarrierConnection{
			ID:                carrierConnection.Carrier.ID,
			Name:              carrierConnection.Carrier.Name,
			ThumbnailURL:      carrierConnection.Carrier.ThumbnailURL,
			SmallThumbnailURL: carrierConnection.Carrier.SmallThumbnailURL,
		},
		Settings: SettingsResponseForGetCarrierConnection{
			GeneralSettings:         carrierConnection.ParsedSettings.GeneralSettings,
			CarrierSpecificSettings: carrierConnection.ParsedSettings.CarrierSpecificSettings,
			EnabledServices:         enabledServices,
		},
	}

}

/* ------------------------- CreateCarrierConnection ------------------------ */

// CreateCarrierConnectionResponse represents response body for the CreateCarrierConnection endpoint
type CreateCarrierConnectionResponse struct {
	ID int `json:"id"`
}

// GenerateCreateCarrierConnectionResponse generates the response body for the CreateCarrierConnection endpoint
func GenerateCreateCarrierConnectionResponse(carrierConnection models.CarrierConnection) *CreateCarrierConnectionResponse {

	return &CreateCarrierConnectionResponse{
		ID: carrierConnection.ID,
	}

}
