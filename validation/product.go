package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/shipply-io/shipply-io-backend/models"
)

/* --------------------------------- Errors --------------------------------- */
var (
	//ErrMainImageRequiredIfAdditionalImagesProvided is returned when a main image is not provided but additional images are provided
	ErrMainImageRequiredIfAdditionalImagesProvided = errors.New("main_image is required if additional images are provided")
)

/* ------------------------------- GetProduct ------------------------------- */

// GetProductRequestData represents the formatted and validated data for the GetProduct endpoint
type GetProductRequestData struct {
	ProductID int
}

// ParseRequestToGetProductRequestData parses the request body to GetProductRequestData
func ParseRequestToGetProductRequestData(r *http.Request) (*GetProductRequestData, error) {

	var err error
	var getProductRequestData GetProductRequestData

	getProductRequestData.ProductID, err = validateIntPathParameter(r, "id", 1)
	if err != nil {
		return nil, err
	}

	return &getProductRequestData, nil

}

/* ---------------------------- GetProductOrders ---------------------------- */

// GetProductOrdersRequestData represents the formatted and validated data for the GetProductOrders endpoint
type GetProductOrdersRequestData struct {
	ProductID int
}

// ParseRequestToGetProductOrdersRequestData parses the request body to GetProductOrdersRequestData
func ParseRequestToGetProductOrdersRequestData(r *http.Request) (*GetProductOrdersRequestData, error) {

	var err error
	var getProductOrdersRequestData GetProductOrdersRequestData

	getProductOrdersRequestData.ProductID, err = validateIntPathParameter(r, "id", 1)
	if err != nil {
		return nil, err
	}

	return &getProductOrdersRequestData, nil
}

/* --------------------------- GetProductInventory -------------------------- */

// GetProductInventoryRequestData represents the formatted and validated data for the GetProductInventory endpoint
type GetProductInventoryRequestData struct {
	ProductID int
}

// ParseRequestToGetProductInventoryRequestData parses the request body to GetProductInventoryRequestData
func ParseRequestToGetProductInventoryRequestData(r *http.Request) (*GetProductInventoryRequestData, error) {

	var err error
	var getProductInventoryRequestData GetProductInventoryRequestData

	getProductInventoryRequestData.ProductID, err = validateIntPathParameter(r, "id", 1)
	if err != nil {
		return nil, err
	}

	return &getProductInventoryRequestData, nil
}

/* ---------------------------- GetProductBundles --------------------------- */

// GetProductBundlesRequestData represents the formatted and validated data for the GetProductBundles endpoint
type GetProductBundlesRequestData struct {
	ProductID int
}

// ParseRequestToGetProductBundlesRequestData parses the request body to GetProductBundlesRequestData
func ParseRequestToGetProductBundlesRequestData(r *http.Request) (*GetProductBundlesRequestData, error) {

	var err error
	var getProductBundlesRequestData GetProductBundlesRequestData

	getProductBundlesRequestData.ProductID, err = validateIntPathParameter(r, "id", 1)
	if err != nil {
		return nil, err
	}

	return &getProductBundlesRequestData, nil

}

/* ----------------------- GetProductBundleComponents ----------------------- */

// GetProductBundleComponentsRequestData represents the formatted and validated data for the GetProductBundleComponents endpoint
type GetProductBundleComponentsRequestData struct {
	ProductID int
}

// ParseRequestToGetProductBundleComponentsRequestData parses the request body to GetProductBundleComponentsRequestData
func ParseRequestToGetProductBundleComponentsRequestData(r *http.Request) (*GetProductBundleComponentsRequestData, error) {

	var err error
	var getProductBundleComponentsRequestData GetProductBundleComponentsRequestData

	getProductBundleComponentsRequestData.ProductID, err = validateIntPathParameter(r, "id", 1)
	if err != nil {
		return nil, err
	}

	return &getProductBundleComponentsRequestData, nil

}

/* ---------------------------- GetProductStores ---------------------------- */

// GetProductStoresRequestData represents the formatted and validated data for the GetProductStores endpoint
type GetProductStoresRequestData struct {
	ProductID int
}

// ParseRequestToGetProductStoresRequestData parses the request body to GetProductStoresRequestData
func ParseRequestToGetProductStoresRequestData(r *http.Request) (*GetProductStoresRequestData, error) {

	var err error
	var getProductStoresRequestData GetProductStoresRequestData

	getProductStoresRequestData.ProductID, err = validateIntPathParameter(r, "id", 1)
	if err != nil {
		return nil, err
	}

	return &getProductStoresRequestData, nil

}

/* ------------------------------ CreateProduct ----------------------------- */

// CreateProductRequestData represents the formatted and validated data for the CreateProduct endpoint
type CreateProductRequestData struct {
	Name        string  `json:"name"`
	ClientID    int     `json:"client_id"`
	Description *string `json:"description"`

	Value   float64 `json:"value"`
	Sku     string  `json:"sku"`
	Barcode *string `json:"barcode"`

	MainImage        *ImageForCreateProductRequestData  `json:"main_image"`
	AdditionalImages []ImageForCreateProductRequestData `json:"images"`

	Dimensions *Dimensions `json:"dimensions"`
	Weight     *Weight     `json:"weight"`
}

// ImageForCreateProductRequestData represents the formatted and validated input for an individual image in the CreateProduct endpoint
type ImageForCreateProductRequestData struct {
	ImageData multipart.File `json:"image_data"`
	FileType  string         `json:"file_type"`
	FileName  string         `json:"file_name"`
}

// ParseRequestToCreateProductRequestData parses and validates the request body for the CreateProduct endpoint
func ParseRequestToCreateProductRequestData(r *http.Request) (*CreateProductRequestData, error) {

	createProductRequestData := CreateProductRequestData{}
	var errs []error

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		return nil, err
	}

	// Main Image
	mainImageFile, err := ValidateOptionalMultipartImage(r, "main_image")
	if err != nil {
		return nil, err
	}

	if mainImageFile != nil {
		createProductRequestData.MainImage = &ImageForCreateProductRequestData{
			ImageData: mainImageFile.FileData,
			FileType:  mainImageFile.FileType,
			FileName:  mainImageFile.FileName,
		}
	}

	// Additional Images
	var additionalImages []ImageForCreateProductRequestData
	i := 1
	for {
		imageKey := fmt.Sprintf("additional_image_%d", i)
		imageFile, err := ValidateOptionalMultipartImage(r, imageKey)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		if imageFile == nil {
			break
		}

		additionalImages = append(additionalImages, ImageForCreateProductRequestData{
			ImageData: imageFile.FileData,
			FileType:  imageFile.FileType,
			FileName:  imageFile.FileName,
		})

		i++
	}
	createProductRequestData.AdditionalImages = additionalImages

	// Validate that a main image is provided if additional images are provided
	if mainImageFile == nil && len(additionalImages) > 0 {
		errs = append(errs, ErrMainImageRequiredIfAdditionalImagesProvided)
	}

	// Retrieve non-image JSON data (stored within the data field)
	dataField := r.FormValue("data")

	// Parse the JSON data into a map
	var rawData map[string]json.RawMessage
	if err := json.Unmarshal([]byte(dataField), &rawData); err != nil {
		errs = append(errs, ErrInvalidJSON)
		return nil, errors.Join(errs...)
	}

	// Validate each field using the new validation functions
	user, _ := models.GetRequestingUser(r)
	if !user.IsClientRole() {
		if createProductRequestData.ClientID, err = validateRequiredIntField(rawData["client_id"], "client_id", 1); err != nil {
			errs = append(errs, err)
		}
	} else {
		createProductRequestData.ClientID = user.OwnerID
	}

	if createProductRequestData.Name, err = validateRequiredStringField(rawData["name"], "name", 1); err != nil {
		errs = append(errs, err)
	}

	if createProductRequestData.Description, err = validateOptionalStringField(rawData["description"], "description", 1); err != nil {
		errs = append(errs, err)
	}

	if createProductRequestData.Value, err = validateRequiredFloatField(rawData["value"], "value", 0); err != nil {
		errs = append(errs, err)
	}

	if createProductRequestData.Sku, err = validateRequiredStringField(rawData["sku"], "sku", 1, 255); err != nil {
		errs = append(errs, err)
	}

	if createProductRequestData.Barcode, err = validateOptionalStringField(rawData["barcode"], "barcode", 1, 255); err != nil {
		errs = append(errs, err)
	}

	if rawData["dimensions"] != nil {
		if createProductRequestData.Dimensions, err = validateDimensionsField(rawData["dimensions"]); err != nil {
			errs = append(errs, err)
		}
	}

	if rawData["weight"] != nil {
		if createProductRequestData.Weight, err = validateWeightField(rawData["weight"]); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return &createProductRequestData, nil

}

/* ------------------------------ UpdateProduct ----------------------------- */
type UpdateProductRequestData struct {
	ProductID int

	Name        *string `json:"name"`
	Description *string `json:"description"`

	Value   *float64 `json:"value"`
	Sku     *string  `json:"sku"`
	Barcode *string  `json:"barcode"`

	Dimensions *Dimensions `json:"dimensions"`
	Weight     *Weight     `json:"weight"`
}

func ParseRequestToUpdateProductRequestData(r *http.Request) (*UpdateProductRequestData, error) {

	updateProductRequestData := UpdateProductRequestData{}
	var err error
	var errs []error

	updateProductRequestData.ProductID, err = validateIntPathParameter(r, "id", 1)
	if err != nil {
		return nil, err
	}

	rawData, err := ParseJSONRequestBody(r)

	if updateProductRequestData.Name, err = validateOptionalStringField(rawData["name"], "name", 1); err != nil {
		errs = append(errs, err)
	}

	if updateProductRequestData.Description, err = validateOptionalStringField(rawData["description"], "description", 1); err != nil {
		errs = append(errs, err)
	}

	if updateProductRequestData.Value, err = validateOptionalFloatField(rawData["value"], "value", 0); err != nil {
		errs = append(errs, err)
	}

	if updateProductRequestData.Sku, err = validateOptionalStringField(rawData["sku"], "sku", 1, 255); err != nil {
		errs = append(errs, err)
	}

	if updateProductRequestData.Barcode, err = validateOptionalStringField(rawData["barcode"], "barcode", 1, 255); err != nil {
		errs = append(errs, err)
	}

	if rawData["dimensions"] != nil {
		if updateProductRequestData.Dimensions, err = validateDimensionsField(rawData["dimensions"]); err != nil {
			errs = append(errs, err)
		}
	}

	if rawData["weight"] != nil {
		if updateProductRequestData.Weight, err = validateWeightField(rawData["weight"]); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return &updateProductRequestData, nil

}

/* ----------------------------- AddProductImage ---------------------------- */

type AddProductImageRequestData struct {
	ProductID int
	Image     *ImageForAddProductImageRequestData `json:"image"`
}

type ImageForAddProductImageRequestData struct {
	ImageData multipart.File `json:"image_data"`
	FileType  string         `json:"file_type"`
	FileName  string         `json:"file_name"`
}

func ParseRequestToAddProductImageRequestData(r *http.Request) (*AddProductImageRequestData, error) {

	addProductImageRequestData := AddProductImageRequestData{}
	var err error
	var errs []error

	addProductImageRequestData.ProductID, err = validateIntPathParameter(r, "id", 1)
	if err != nil {
		errs = append(errs, err)
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		return nil, err
	}

	// Main Image
	imageFile, err := ValidateRequiredMultipartImage(r, "image")
	if err != nil {
		errs = append(errs, err)
	}

	addProductImageRequestData.Image = &ImageForAddProductImageRequestData{
		ImageData: imageFile.FileData,
		FileType:  imageFile.FileType,
		FileName:  imageFile.FileName,
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return &addProductImageRequestData, nil

}

/* --------------------------- DeleteProductImage --------------------------- */

// DeleteProductImageRequestData represents the formatted and validated data for the DeleteProductImage endpoint
type DeleteProductImageRequestData struct {
	ProductID      int
	ProductImageID int
}

// ParseRequestToDeleteProductImageRequestData parses and validates the request body for the DeleteProductImage endpoint
func ParseRequestToDeleteProductImageRequestData(r *http.Request) (*DeleteProductImageRequestData, error) {

	deleteProductImageRequestData := DeleteProductImageRequestData{}
	var err error
	var errs []error

	deleteProductImageRequestData.ProductID, err = validateIntPathParameter(r, "id", 1)
	if err != nil {
		errs = append(errs, err)
	}

	deleteProductImageRequestData.ProductImageID, err = validateIntPathParameter(r, "product_image_id", 1)
	if err != nil {
		errs = append(errs, err)
	}

	return &deleteProductImageRequestData, nil

}

/* ------------------------- UpdateProductImageOrder ------------------------ */

// UpdateProductImageOrderRequestData represents the formatted and validated data for the UpdateProductImageOrder endpoint
type UpdateProductImageOrderRequestData struct {
	ProductID int
	Order     []int `json:"order"`
}

// ParseRequestToUpdateProductImageOrderRequestData parses and validates the request body for the UpdateProductImageOrder endpoint
func ParseRequestToUpdateProductImageOrderRequestData(r *http.Request) (*UpdateProductImageOrderRequestData, error) {

	updateProductImageOrderRequestData := UpdateProductImageOrderRequestData{}
	var err error
	var errs []error

	updateProductImageOrderRequestData.ProductID, err = validateIntPathParameter(r, "id", 1)
	if err != nil {
		errs = append(errs, err)
	}

	rawData, err := ParseJSONRequestBody(r)
	if err != nil {
		return nil, err
	}

	if updateProductImageOrderRequestData.Order, err = validateRequiredIntArrayField(rawData["order"], "order", 1); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return &updateProductImageOrderRequestData, nil

}
