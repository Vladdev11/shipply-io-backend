package ShipengineAPI

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	ShipengineModels "github.com/shipply-io/shipply-io-backend/api/shipengine/models"
	"github.com/shipply-io/shipply-io-backend/models"
)

// ValidateAddress validates an address
// Possible return values:
// Verified: The address is valid and verified.
// Unverified: Unable to verify the address
// Warning - The address is valid, but there is a warning about the address.
// Error - The address is invalid.
func ValidateAddress(addresses []ShipengineModels.Address) ([]*ShipengineModels.AddressResponse, *ShipengineModels.ShipengineError, error) {
	body, err := json.Marshal(addresses)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequest("POST", apiClient.GetApiHost()+"/addresses/validate", bytes.NewBuffer(body))
	if err != nil {
		return nil, nil, err
	}

	resp, err := apiClient.Do(req)
	if err != nil {
		return nil, nil, err
	}

	bodyBytes, err := io.ReadAll(resp.Body) // Read the entire body into a byte slice
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	log := models.ShipEngineLog{
		EventType:  "api_request",
		Headers:    resp.Header,
		Endpoint:   "/addresses/validate",
		StatusCode: resp.StatusCode,
		RequestURL: apiClient.GetApiHost() + "/addresses/validate",
		Method:     "POST",
		CreatedAt:  time.Now(),
		Body:       string(bodyBytes),
	}
	log.Create()

	if resp.StatusCode != 200 {
		var errorResponse ShipengineModels.ShipengineError
		err = json.Unmarshal(bodyBytes, &errorResponse) // Use json.Unmarshal with bodyBytes
		if err != nil {
			return nil, nil, err
		}
		return nil, &errorResponse, nil
	}

	var result []*ShipengineModels.AddressResponse
	err = json.Unmarshal(bodyBytes, &result) // Use json.Unmarshal with bodyBytes
	if err != nil {
		return nil, nil, err
	}

	return result, nil, nil
}

func ParseAddress(toParse ShipengineModels.AddressParseRequest) (*ShipengineModels.AddressParseResponse, *ShipengineModels.ShipengineError, error) {
	body, err := json.Marshal(toParse)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequest("PUT", apiClient.GetApiHost()+"/addresses/recognize", bytes.NewBuffer(body))
	if err != nil {
		return nil, nil, err
	}

	resp, err := apiClient.Do(req)
	if err != nil {
		return nil, nil, err
	}

	log := models.ShipEngineLog{
		EventType:  "api_request",
		Headers:    resp.Header,
		Body:       string(body),
		Endpoint:   "/addresses/recognize",
		StatusCode: resp.StatusCode,
		RequestURL: apiClient.GetApiHost() + "/addresses/recognize",
		Method:     "PUT",
		CreatedAt:  time.Now(),
	}
	log.Create()

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		var errorResponse ShipengineModels.ShipengineError
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		if err != nil {
			return nil, nil, err
		}
		return nil, &errorResponse, nil
	}

	var result *ShipengineModels.AddressParseResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, nil, err
	}

	return result, nil, nil
}
