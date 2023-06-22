package ShipengineAPI

import (
	"bytes"
	"encoding/json"
	"net/http"

	ShipengineModels "github.com/shipply-io/shipply-io-backend/api/shipengine/models"
)

func CreateShipment(req ShipengineModels.CreateShipmentRequest) (*ShipengineModels.CreateShipmentResponse, *ShipengineModels.ShipengineError, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, nil, err
	}

	request, err := http.NewRequest("POST", apiClient.GetApiHost()+"/shipments", bytes.NewBuffer(body))
	if err != nil {
		return nil, nil, err
	}

	response, err := apiClient.Do(request)
	if err != nil {
		return nil, nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		var errorResponse ShipengineModels.ShipengineError
		err = json.NewDecoder(response.Body).Decode(&errorResponse)
		if err != nil {
			return nil, nil, err
		}
		return nil, &errorResponse, nil
	}

	var result *ShipengineModels.CreateShipmentResponse
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, nil, err
	}

	return result, nil, nil
}
