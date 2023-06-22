package ShipengineAPI

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	ShipengineModels "github.com/shipply-io/shipply-io-backend/api/shipengine/models"
	"github.com/shipply-io/shipply-io-backend/models"
)

func CreatePackage(p ShipengineModels.Package) (*ShipengineModels.PackageResponse, *ShipengineModels.ShipengineError, error) {
	body, err := json.Marshal(p)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequest("POST", apiClient.GetApiHost()+"/packages", bytes.NewBuffer(body))
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
		Endpoint:   "/packages",
		StatusCode: resp.StatusCode,
		RequestURL: apiClient.GetApiHost() + "/packages",
		Method:     "POST",
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

	var result *ShipengineModels.PackageResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, nil, err
	}

	return result, nil, nil
}
