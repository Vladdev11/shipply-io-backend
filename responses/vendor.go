package responses

import "github.com/shipply-io/shipply-io-backend/models"

/* -------------------------------- ListVendors -------------------------------- */

// ListVendorsResponse represents the expected response body for the ListVendors endpoint
type ListVendorsResponse struct {
	TotalCount    int                            `json:"total_count"`
	FilteredCount int                            `json:"filtered_count"`
	Data          []VendorResponseForListVendors `json:"data"`
}

// VendorResponseForListVendors represents the expected response body for the ListVendors endpoint
type VendorResponseForListVendors struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GenerateListVendorsResponse generates the response body for the ListVendors endpoint
func GenerateListVendorsResponse(vendors []models.Vendor, totalCount int, filteredCount int) *ListVendorsResponse {

	var vendorResponses []VendorResponseForListVendors
	for _, vendor := range vendors {
		vendorResponses = append(vendorResponses, VendorResponseForListVendors{
			ID:   vendor.ID,
			Name: vendor.Name,
		})
	}

	return &ListVendorsResponse{
		TotalCount:    totalCount,
		FilteredCount: filteredCount,
		Data:          vendorResponses,
	}

}

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

/* -------------------------------- CreateVendor -------------------------------- */

// CreateVendorResponse represents the expected response body for the CreateVendor endpoint
type CreateVendorResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GenerateCreateVendorResponse generates the response body for the CreateVendor endpoint
func GenerateCreateVendorResponse(vendor models.Vendor) *CreateVendorResponse {
	return &CreateVendorResponse{
		ID:   vendor.ID,
		Name: vendor.Name,
	}
}

/* -------------------------------- UpdateVendor -------------------------------- */

// UpdateVendorResponse represents the expected response body for the UpdateVendor endpoint
type UpdateVendorResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GenerateUpdateVendorResponse generates the response body for the UpdateVendor endpoint
func GenerateUpdateVendorResponse(vendor models.Vendor) *UpdateVendorResponse {
	return &UpdateVendorResponse{
		ID:   vendor.ID,
		Name: vendor.Name,
	}
}
