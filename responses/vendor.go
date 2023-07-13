package responses

import "github.com/shipply-io/shipply-io-backend/models"

/* -------------------------------- GetVendor ------------------------------- */

// GetVendorResponse represents the expected response body for the GetVendor endpoint
type GetVendorResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ClientID int    `json:"client_id"`
}

// GenerateGetVendorResponse generates the response body for the GetVendor endpoint
func GenerateGetVendorResponse(vendor models.Vendor) *GetVendorResponse {

	return &GetVendorResponse{
		ID:       vendor.ID,
		Name:     vendor.Name,
		ClientID: vendor.ClientID,
	}

}
