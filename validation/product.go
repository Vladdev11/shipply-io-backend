package validation

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/shipply-io/shipply-io-backend/models"
)

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
func ParseRequestToCreateProductRequestData(r *http.Request) (*CreateProductRequestData, []string) {

	createProductRequestData := CreateProductRequestData{}
	errors := []string{}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		errors = append(errors, err.Error())
		return nil, errors
	}

	// Main Image
	mainImageFile, err := ParseMultipartImage(r, "main_image")
	if err != nil && err != http.ErrMissingFile {
		errors = append(errors, err.Error())
		return nil, errors
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
		imageFile, err := ParseMultipartImage(r, imageKey)
		if err != nil {
			if err == http.ErrMissingFile {
				break
			}
			errors = append(errors, err.Error())
			continue
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
		errors = append(errors, "main_image is required if additional images are provided")
	}

	// Retrieve non-image JSON data (stored within the data field)
	dataField := r.FormValue("data")

	// Parse the JSON data into a map
	var rawData map[string]json.RawMessage
	if err := json.Unmarshal([]byte(dataField), &rawData); err != nil {
		errors = append(errors, "invalid JSON data")
		return nil, errors
	}

	// Validate each field using the new validation functions
	user, _ := models.GetRequestingUser(r)
	if !user.IsClientRole() {
		if createProductRequestData.ClientID, err = validateRequiredIntField(rawData["client_id"], "client_id", 1); err != nil {
			errors = append(errors, err.Error())
		}
	} else {
		createProductRequestData.ClientID = user.OwnerID
	}

	if createProductRequestData.Name, err = validateRequiredStringField(rawData["name"], "name", 1); err != nil {
		errors = append(errors, err.Error())
	}

	if createProductRequestData.Description, err = validateOptionalStringField(rawData["description"], "description", 1); err != nil {
		errors = append(errors, err.Error())
	}

	if createProductRequestData.Value, err = validateRequiredFloatField(rawData["value"], "value", 0); err != nil {
		errors = append(errors, err.Error())
	}

	if createProductRequestData.Sku, err = validateRequiredStringField(rawData["sku"], "sku", 1, 255); err != nil {
		errors = append(errors, err.Error())
	}

	if createProductRequestData.Barcode, err = validateOptionalStringField(rawData["barcode"], "barcode", 1, 255); err != nil {
		errors = append(errors, err.Error())
	}

	if rawData["dimensions"] != nil {
		if createProductRequestData.Dimensions, err = validateDimensionsField(rawData["dimensions"]); err != nil {
			errors = append(errors, err.Error())
		}
	}

	if rawData["weight"] != nil {
		if createProductRequestData.Weight, err = validateWeightField(rawData["weight"]); err != nil {
			errors = append(errors, err.Error())
		}
	}

	if len(errors) > 0 {
		return nil, errors
	}

	return &createProductRequestData, nil

}

/* ------------------------------ UpdateProduct ----------------------------- */
type UpdateProductRequestData struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`

	Value   *float64 `json:"value"`
	Sku     *string  `json:"sku"`
	Barcode *string  `json:"barcode"`

	Dimensions *Dimensions `json:"dimensions"`
	Weight     *Weight     `json:"weight"`
}

func ParseRequestToUpdateProductRequestData(r *http.Request) (*UpdateProductRequestData, []string) {

	updateProductRequestData := UpdateProductRequestData{}
	errors := []string{}

	rawData, err := ParseJSONRequestBody(r)

	if updateProductRequestData.Name, err = validateOptionalStringField(rawData["name"], "name", 1); err != nil {
		errors = append(errors, err.Error())
	}

	if updateProductRequestData.Description, err = validateOptionalStringField(rawData["description"], "description", 1); err != nil {
		errors = append(errors, err.Error())
	}

	if updateProductRequestData.Value, err = validateOptionalFloatField(rawData["value"], "value", 0); err != nil {
		errors = append(errors, err.Error())
	}

	if updateProductRequestData.Sku, err = validateOptionalStringField(rawData["sku"], "sku", 1, 255); err != nil {
		errors = append(errors, err.Error())
	}

	if updateProductRequestData.Barcode, err = validateOptionalStringField(rawData["barcode"], "barcode", 1, 255); err != nil {
		errors = append(errors, err.Error())
	}

	if rawData["dimensions"] != nil {
		if updateProductRequestData.Dimensions, err = validateDimensionsField(rawData["dimensions"]); err != nil {
			errors = append(errors, err.Error())
		}
	}

	if rawData["weight"] != nil {
		if updateProductRequestData.Weight, err = validateWeightField(rawData["weight"]); err != nil {
			errors = append(errors, err.Error())
		}
	}

	if len(errors) > 0 {
		return nil, errors
	}

	return &updateProductRequestData, nil

}

/* ----------------------------- AddProductImage ---------------------------- */

type AddProductImageRequestData struct {
	Image *ImageForAddProductImageRequestData `json:"image"`
}

type ImageForAddProductImageRequestData struct {
	ImageData multipart.File `json:"image_data"`
	FileType  string         `json:"file_type"`
	FileName  string         `json:"file_name"`
}

func ParseRequestToAddProductImageRequestData(r *http.Request) (*AddProductImageRequestData, []string) {

	addProductImageRequestData := AddProductImageRequestData{}
	errors := []string{}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		errors = append(errors, err.Error())
		return nil, errors
	}

	// Main Image
	imageFile, err := ParseMultipartImage(r, "image")
	if err != nil {
		errors = append(errors, err.Error())
		return nil, errors
	}

	addProductImageRequestData.Image = &ImageForAddProductImageRequestData{
		ImageData: imageFile.FileData,
		FileType:  imageFile.FileType,
		FileName:  imageFile.FileName,
	}

	return &addProductImageRequestData, nil

}

/* ------------------------- UpdateProductImageOrder ------------------------ */

// UpdateProductImageOrderRequestData represents the formatted and validated data for the UpdateProductImageOrder endpoint
type UpdateProductImageOrderRequestData struct {
	Order []int `json:"order"`
}

// ParseRequestToUpdateProductImageOrderRequestData parses and validates the request body for the UpdateProductImageOrder endpoint
func ParseRequestToUpdateProductImageOrderRequestData(r *http.Request) (*UpdateProductImageOrderRequestData, []string) {

	updateProductImageOrderRequestData := UpdateProductImageOrderRequestData{}
	errors := []string{}

	rawData, err := ParseJSONRequestBody(r)

	if updateProductImageOrderRequestData.Order, err = validateRequiredIntArrayField(rawData["order"], "order", 1); err != nil {
		errors = append(errors, err.Error())
	}

	if len(errors) > 0 {
		return nil, errors
	}

	return &updateProductImageOrderRequestData, nil

}
