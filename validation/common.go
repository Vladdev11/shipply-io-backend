package validation

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

/* ------------------------ Path Parameter Validation ----------------------- */

func validateIntPathParameter(r *http.Request, key string, min int) (int, error) {
	vars := mux.Vars(r)
	stringValue := vars[key]
	if stringValue == "" {
		return 0, ErrMissingPathParameter{Parameter: key}
	}

	value, err := strconv.Atoi(stringValue)
	if err != nil {
		return 0, ErrInvalidPathParameterType{Parameter: key, ExpectedType: "integer"}
	}
	return value, nil
}

func validateStringPathParameter(r *http.Request, key string, minLen int, maxLen ...int) (string, error) {
	vars := mux.Vars(r)
	value := vars[key]
	if value == "" {
		return "", ErrMissingPathParameter{Parameter: key}
	} else if len(value) < minLen {
		return "", ErrInvalidPathParameterLength{Parameter: key, Min: minLen}
	} else if len(maxLen) > 0 && len(value) > maxLen[0] {
		return "", ErrInvalidPathParameterLength{Parameter: key, Min: minLen, Max: maxLen[0]}
	} else {
		return value, nil
	}
}

/* --------------------------------- Custom Type Structs and Validation Functions (Body) -------------------------------- */

// Dimensions represents the expected response body for the Dimensions used by multiple endpoints
type Dimensions struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// validateDimensionsField validates the Dimensions field
func validateDimensionsField(rawData json.RawMessage) (*Dimensions, error) {
	dimensions := &Dimensions{}

	if err := json.Unmarshal(rawData, dimensions); err != nil {
		return nil, ErrInvalidField{Field: "dimensions"}
	}

	if dimensions.Length < 0 {
		return nil, ErrInvalidIntFieldRange{Field: "length", Min: 0}
	} else if dimensions.Width < 0 {
		return nil, ErrInvalidIntFieldRange{Field: "width", Min: 0}
	} else if dimensions.Height < 0 {
		return nil, ErrInvalidIntFieldRange{Field: "height", Min: 0}
	}

	return dimensions, nil
}

// Weight represents the expected response body for the Weight used by multiple endpoints
type Weight struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

// validateWeightField validates the Weight field
func validateWeightField(rawData json.RawMessage) (*Weight, error) {

	weight := &Weight{}

	if err := json.Unmarshal(rawData, weight); err != nil {
		return nil, ErrInvalidField{Field: "weight"}
	}

	if weight.Value < 0 {
		return nil, ErrInvalidIntFieldRange{Field: "value", Min: 0}
	} else if weight.Unit == "" {
		return nil, ErrMissingField{Field: "unit"}
	}

	return weight, nil
}

/* -------------------------------- Native Type Validation Functions (Body) ------------------------------- */

// validateRequiredIntField validates a required int field
func validateRequiredIntField(raw json.RawMessage, fieldName string, min int) (int, error) {
	var value int
	if raw == nil {
		return 0, ErrMissingField{Field: fieldName}
	} else if err := json.Unmarshal(raw, &value); err != nil {
		return 0, ErrInvalidFieldType{Field: fieldName, ExpectedType: "integer"}
	} else if value < min {
		return 0, ErrInvalidIntFieldRange{Field: fieldName, Min: min}
	} else {
		return value, nil
	}
}

// validateOptionalIntField validates an optional int field
func validateOptionalIntField(raw json.RawMessage, fieldName string, min int) (*int, error) {
	var value int
	if raw == nil {
		return nil, nil
	} else if err := json.Unmarshal(raw, &value); err != nil {
		return nil, ErrInvalidFieldType{Field: fieldName, ExpectedType: "integer"}
	} else if value < min {
		return nil, ErrInvalidIntFieldRange{Field: fieldName, Min: min}
	} else {
		return &value, nil
	}
}

// validateRequiredFloatField validates a required float field
func validateRequiredFloatField(raw json.RawMessage, fieldName string, min float64, max ...float64) (float64, error) {
	var value float64
	if raw == nil {
		return 0, ErrMissingField{Field: fieldName}
	} else if err := json.Unmarshal(raw, &value); err != nil {
		return 0, ErrInvalidFieldType{Field: fieldName, ExpectedType: "number"}
	} else if value < min {
		return 0, ErrInvalidFloatFieldRange{Field: fieldName, Min: min}
	} else if len(max) > 0 && value > max[0] {
		return 0, ErrInvalidFloatFieldRange{Field: fieldName, Min: min, Max: max[0]}
	} else {
		return value, nil
	}
}

// validateOptionalFloatField validates an optional float field
func validateOptionalFloatField(raw json.RawMessage, fieldName string, min float64, max ...float64) (*float64, error) {
	var value float64
	if raw == nil {
		return nil, nil
	} else if err := json.Unmarshal(raw, &value); err != nil {
		return nil, ErrInvalidFieldType{Field: fieldName, ExpectedType: "number"}
	} else if value < min {
		return nil, ErrInvalidFloatFieldRange{Field: fieldName, Min: min}
	} else if len(max) > 0 && value > max[0] {
		return nil, ErrInvalidFloatFieldRange{Field: fieldName, Min: min, Max: max[0]}
	} else {
		return &value, nil
	}
}

// validateRequiredStringField validates a required string field
func validateRequiredStringField(raw json.RawMessage, fieldName string, minLen int, maxLen ...int) (string, error) {
	var value string
	if raw == nil {
		return "", ErrMissingField{Field: fieldName}
	} else if err := json.Unmarshal(raw, &value); err != nil {
		return "", ErrInvalidFieldType{Field: fieldName, ExpectedType: "string"}
	} else if len(value) < minLen {
		return "", ErrInvalidStringFieldLength{Field: fieldName, Min: minLen}
	} else if len(maxLen) > 0 && len(value) > maxLen[0] {
		return "", ErrInvalidStringFieldLength{Field: fieldName, Min: minLen, Max: maxLen[0]}
	} else {
		return value, nil
	}
}

// validateOptionalStringField validates an optional string field
func validateOptionalStringField(raw json.RawMessage, fieldName string, minLen int, maxLen ...int) (*string, error) {
	var value string
	if raw == nil {
		return nil, nil
	} else if err := json.Unmarshal(raw, &value); err != nil {
		return nil, ErrInvalidFieldType{Field: fieldName, ExpectedType: "string"}
	} else if len(value) < minLen {
		return nil, ErrInvalidStringFieldLength{Field: fieldName, Min: minLen}
	} else if len(maxLen) > 0 && len(value) > maxLen[0] {
		return nil, ErrInvalidStringFieldLength{Field: fieldName, Min: minLen, Max: maxLen[0]}
	} else {
		return &value, nil
	}
}

// validateRequiredIntArrayField validates a required int array field
func validateRequiredIntArrayField(raw json.RawMessage, fieldName string, minLen int, maxLen ...int) ([]int, error) {
	var value []int
	if raw == nil {
		return nil, ErrMissingField{Field: fieldName}
	} else if err := json.Unmarshal(raw, &value); err != nil {
		return nil, ErrInvalidFieldType{Field: fieldName, ExpectedType: "array of integers"}
	} else if len(value) < minLen {
		return nil, ErrInvalidArrayFieldLength{Field: fieldName, Min: minLen}
	} else if len(maxLen) > 0 && len(value) > maxLen[0] {
		return nil, ErrInvalidArrayFieldLength{Field: fieldName, Min: minLen, Max: maxLen[0]}
	} else {
		return value, nil
	}
}

/* ---------------------------------- Files --------------------------------- */

// MultipartFileData represents a custom struct that is returned when parsing a multipart file upload
type MultipartFileData struct {
	FileData multipart.File
	FileName string
	FileType string
}

// ParseMultipartFile parses a multipart file upload
func ParseMultipartFile(r *http.Request, key string) (*MultipartFileData, error) {
	file, fileHeader, err := r.FormFile(key)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a buffer to store a portion of the file
	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil {
		return nil, err
	}

	// Detect the content type of the file
	contentType := http.DetectContentType(buf)

	// Make sure to reset the read pointer to the beginning of the file after reading
	file.Seek(0, 0)

	fileData := &MultipartFileData{
		FileData: file,
		FileName: fileHeader.Filename,
		FileType: contentType,
	}

	return fileData, nil
}

// ParseMultipartImage parses a multipart image upload
func ParseMultipartImage(r *http.Request, key string) (*MultipartFileData, error) {

	fileData, err := ParseMultipartFile(r, key)
	if err != nil {
		if err == http.ErrMissingFile {
			return nil, ErrMissingField{Field: key}
		}

		return nil, err
	}

	if fileData.FileType != "image/jpeg" && fileData.FileType != "image/png" {
		return nil, ErrInvalidFieldType{Field: key, ExpectedType: "image/jpeg or image/png"}
	}

	return fileData, nil
}

/* ---------------------------- Common Functions ---------------------------- */
func ParseJSONRequestBody(r *http.Request) (map[string]json.RawMessage, error) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, ErrUnableToReadBody
	}

	// Parse the JSON data into a map
	var rawData map[string]json.RawMessage
	if err := json.Unmarshal(body, &rawData); err != nil {
		return nil, ErrInvalidJSON
	}

	return rawData, nil
}
