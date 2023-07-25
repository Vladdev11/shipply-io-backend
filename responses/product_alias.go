package responses

import "github.com/shipply-io/shipply-io-backend/models"

/* ----------------------- GetProductAliasGetByBarcode ---------------------- */
// GetProductAliasGetByBarcodeResponse represents the response body for the GetProductAliasGetByBarcode endpoint
type GetProductAliasGetByBarcodeResponse struct {
	ProductID int    `json:"product_id"`
	Barcode   string `json:"barcode"`
	Quantity  int    `json:"quantity"`
}

// GenerateGetProductAliasGetByBarcodeResponse generates the response body for the GetProductAliasGetByBarcode endpoint
func GenerateGetProductAliasGetByBarcodeResponse(productAlias models.ProductAlias) *GetProductAliasGetByBarcodeResponse {
	return &GetProductAliasGetByBarcodeResponse{
		ProductID: productAlias.ProductID,
		Barcode:   productAlias.Barcode,
		Quantity:  productAlias.Quantity,
	}
}
