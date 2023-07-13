package validation

import (
	"encoding/json"
	"fmt"
)

/* --------------------------------- Custom Type Structs and Validation Functions -------------------------------- */

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
		return nil, fmt.Errorf("invalid dimensions data")
	}

	if dimensions.Length < 0 {
		return nil, fmt.Errorf("length must be greater than or equal to 0")
	} else if dimensions.Width < 0 {
		return nil, fmt.Errorf("width must be greater than or equal to 0")
	} else if dimensions.Height < 0 {
		return nil, fmt.Errorf("height must be greater than or equal to 0")
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
		return nil, fmt.Errorf("invalid weight data")
	}

	if weight.Value < 0 {
		return nil, fmt.Errorf("weight must be greater than or equal to 0")
	} else if weight.Unit == "" {
		return nil, fmt.Errorf("unit is required")
	}

	return weight, nil
}

/* -------------------------------- Native Type Validation Functions ------------------------------- */

// validateRequiredIntField validates a required int field
func validateRequiredIntField(raw json.RawMessage, fieldName string, min int) (int, error) {
	var value int
	if raw == nil {
		return 0, fmt.Errorf("%s is required", fieldName)
	} else if err := json.Unmarshal(raw, &value); err != nil {
		return 0, fmt.Errorf("%s must be an integer", fieldName)
	} else if value < min {
		return 0, fmt.Errorf("%s must be greater than %d", fieldName, min-1)
	} else {
		return value, nil
	}
}

// validateOptionalIntField validates an optional int field
func validateOptionalIntField(raw json.RawMessage, fieldName string, min int) (int, error) {
	var value int
	if raw == nil {
		return 0, nil
	} else if err := json.Unmarshal(raw, &value); err != nil {
		return 0, fmt.Errorf("%s must be an integer", fieldName)
	} else if value < min {
		return 0, fmt.Errorf("%s must be greater than %d", fieldName, min-1)
	} else {
		return value, nil
	}
}

// validateRequiredFloatField validates a required float field
func validateRequiredFloatField(raw json.RawMessage, fieldName string, min float64) (float64, error) {
	var value float64
	if raw == nil {
		return 0, fmt.Errorf("%s is required", fieldName)
	} else if err := json.Unmarshal(raw, &value); err != nil {
		return 0, fmt.Errorf("%s must be a number", fieldName)
	} else if value < min {
		return 0, fmt.Errorf("%s must be greater than or equal to %f", fieldName, min)
	} else {
		return value, nil
	}
}

// validateOptionalFloatField validates an optional float field
func validateOptionalFloatField(raw json.RawMessage, fieldName string, min float64) (float64, error) {
	var value float64
	if raw == nil {
		return 0, nil
	} else if err := json.Unmarshal(raw, &value); err != nil {
		return 0, fmt.Errorf("%s must be a number", fieldName)
	} else if value < min {
		return 0, fmt.Errorf("%s must be greater than or equal to %f", fieldName, min)
	} else {
		return value, nil
	}
}

// validateRequiredStringField validates a required string field
func validateRequiredStringField(raw json.RawMessage, fieldName string, minLen int, maxLen ...int) (string, error) {
	var value string
	if raw == nil {
		return "", fmt.Errorf("%s is required", fieldName)
	} else if err := json.Unmarshal(raw, &value); err != nil {
		return "", fmt.Errorf("%s must be a string", fieldName)
	} else if len(value) < minLen {
		return "", fmt.Errorf("%s must be at least %d characters", fieldName, minLen)
	} else if len(maxLen) > 0 && len(value) > maxLen[0] {
		return "", fmt.Errorf("%s must be less than or equal to %d characters", fieldName, maxLen[0])
	} else {
		return value, nil
	}
}

// validateOptionalStringField validates an optional string field
func validateOptionalStringField(raw json.RawMessage, fieldName string, minLen int, maxLen ...int) (string, error) {
	var value string
	if raw == nil {
		return "", nil
	} else if err := json.Unmarshal(raw, &value); err != nil {
		return "", fmt.Errorf("%s must be a string", fieldName)
	} else if len(value) < minLen {
		return "", fmt.Errorf("%s must be at least %d characters", fieldName, minLen)
	} else if len(maxLen) > 0 && len(value) > maxLen[0] {
		return "", fmt.Errorf("%s must be less than or equal to %d characters", fieldName, maxLen[0])
	} else {
		return value, nil
	}
}
