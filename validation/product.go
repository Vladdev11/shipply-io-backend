package validation

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/shipply-io/shipply-io-backend/models"
)

/* ------------------------------ CreateProduct ----------------------------- */

// CreateProductInput represents the formatted and validated input for the CreateProduct endpoint
type CreateProductInput struct {
	Name        string `json:"name"`
	ClientID    int    `json:"client_id"`
	Description string `json:"description"`

	Value       float64 `json:"value"`
	Sku         string  `json:"sku"`
	Barcode     string  `json:"barcode"`
	WarehouseID int     `json:"warehouse_id"`

	Images []ImageForCreateProductInput `json:"images"`

	Dimensions *Dimensions `json:"dimensions"`
	Weight     *Weight     `json:"weight"`
}

// ImageForCreateProductInput represents the formatted and validated input for an individual image in the CreateProduct endpoint
type ImageForCreateProductInput struct {
	ImageData multipart.File `json:"image_data"`
	FileType  string         `json:"file_type"`
	FileName  string         `json:"file_name"`
}

// ParseRequestToCreateProductInput parses and validates the request body for the CreateProduct endpoint
func ParseRequestToCreateProductInput(r *http.Request) (*CreateProductInput, []string) {

	createProductInput := CreateProductInput{}
	errors := []string{}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		errors = append(errors, err.Error())
		return nil, errors
	}

	var images []ImageForCreateProductInput
	i := 1
	for {
		imageKey := fmt.Sprintf("file%d", i)
		file, fileHeader, err := r.FormFile(imageKey)
		if err == nil {
			defer file.Close()
			if fileHeader.Header.Get("Content-Type") != "image/jpeg" && fileHeader.Header.Get("Content-Type") != "image/png" {
				errors = append(errors, "file must be a valid image")
				i++
				continue
			}

			images = append(images, ImageForCreateProductInput{
				ImageData: file,
				FileType:  fileHeader.Header.Get("Content-Type"),
				FileName:  fileHeader.Filename,
			})
		} else {
			break
		}
		i++
	}

	createProductInput.Images = images

	// Retrieve non-image JSON data (stored within the data field)
	dataField := r.FormValue("data")

	// Parse the JSON data into a map
	var rawData map[string]json.RawMessage
	if err := json.Unmarshal([]byte(dataField), &rawData); err != nil {
		errors = append(errors, "invalid JSON data")
		return nil, errors
	}

	// Validate each field using the new validation functions
	var err error

	user, _ := models.GetRequestingUser(r)
	if !user.IsClientRole() {
		if createProductInput.ClientID, err = validateRequiredIntField(rawData["client_id"], "client_id", 1); err != nil {
			errors = append(errors, err.Error())
		}
	} else {
		createProductInput.ClientID = user.OwnerID
	}

	if createProductInput.Name, err = validateRequiredStringField(rawData["name"], "name", 1); err != nil {
		errors = append(errors, err.Error())
	}

	if createProductInput.Description, err = validateOptionalStringField(rawData["description"], "description", 1); err != nil {
		errors = append(errors, err.Error())
	}

	if createProductInput.Value, err = validateRequiredFloatField(rawData["value"], "value", 0); err != nil {
		errors = append(errors, err.Error())
	}

	if createProductInput.Sku, err = validateRequiredStringField(rawData["sku"], "sku", 1, 255); err != nil {
		errors = append(errors, err.Error())
	}

	if createProductInput.Barcode, err = validateOptionalStringField(rawData["barcode"], "barcode", 1, 255); err != nil {
		errors = append(errors, err.Error())
	}

	if createProductInput.WarehouseID, err = validateOptionalIntField(rawData["warehouse_id"], "warehouse_id", 1); err != nil {
		errors = append(errors, err.Error())
	}

	if rawData["dimensions"] != nil {
		if createProductInput.Dimensions, err = validateDimensionsField(rawData["dimensions"]); err != nil {
			errors = append(errors, err.Error())
		}
	}

	if rawData["weight"] != nil {
		if createProductInput.Weight, err = validateWeightField(rawData["weight"]); err != nil {
			errors = append(errors, err.Error())
		}
	}

	if len(errors) > 0 {
		return nil, errors
	}

	return &createProductInput, nil

}
