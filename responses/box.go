package responses

import "github.com/shipply-io/shipply-io-backend/models"

/* -------------------------------- ListBoxes ------------------------------- */

// ListBoxesResponse represents the expected response body for the ListBoxes endpoint
type ListBoxesResponse struct {
	Boxes []BoxResponseForListBoxes `json:"boxes"`
}

// BoxResponseForListBoxes represents the expected response body for the ListBoxes endpoint
type BoxResponseForListBoxes struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Barcode string `json:"barcode"`
}

// GenerateListBoxesResponse generates the response body for the ListBoxes endpoint
func GenerateListBoxesResponse(boxes []models.Box) *ListBoxesResponse {

	boxResponses := []BoxResponseForListBoxes{}

	for _, box := range boxes {
		boxResponses = append(boxResponses, BoxResponseForListBoxes{
			ID:      box.ID,
			Name:    box.Name,
			Barcode: box.Barcode,
		})
	}

	return &ListBoxesResponse{
		Boxes: boxResponses,
	}

}

/* --------------------------------- GetBox --------------------------------- */

// GetBoxResponse represents the expected response body for the GetBox endpoint
type GetBoxResponse struct {
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	Barcode    string     `json:"barcode"`
	Type       string     `json:"type"`
	Active     bool       `json:"active"`
	Cost       float64    `json:"cost"`
	Dimensions Dimensions `json:"dimensions"`
}

// GenerateGetBoxResponse generates the response body for the GetBox endpoint
func GenerateGetBoxResponse(box models.Box) *GetBoxResponse {

	return &GetBoxResponse{
		ID:      box.ID,
		Name:    box.Name,
		Barcode: box.Barcode,
		Type:    box.Type,
		Active:  box.Active,
		Cost:    box.Cost,
		Dimensions: Dimensions{
			Length: box.Length,
			Width:  box.Width,
			Height: box.Height,
		},
	}

}

/* -------------------------------- CreateBox -------------------------------- */

// CreateBoxResponse represents the expected response body for the CreateBox endpoint
type CreateBoxResponse struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Barcode string `json:"barcode"`
}

// GenerateCreateBoxResponse generates the response body for the CreateBox endpoint
func GenerateCreateBoxResponse(box models.Box) *CreateBoxResponse {

	return &CreateBoxResponse{
		ID:      box.ID,
		Name:    box.Name,
		Barcode: box.Barcode,
	}

}

/* -------------------------------- UpdateBox -------------------------------- */

// UpdateBoxResponse represents the expected response body for the UpdateBox endpoint
type UpdateBoxResponse struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Barcode string `json:"barcode"`
}

// GenerateUpdateBoxResponse generates the response body for the UpdateBox endpoint
func GenerateUpdateBoxResponse(box models.Box) *UpdateBoxResponse {

	return &UpdateBoxResponse{
		ID:      box.ID,
		Name:    box.Name,
		Barcode: box.Barcode,
	}

}
