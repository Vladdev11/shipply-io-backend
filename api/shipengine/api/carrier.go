package ShipengineAPI

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	ShipengineModels "github.com/shipply-io/shipply-io-backend/api/shipengine/models"
	"github.com/shipply-io/shipply-io-backend/models"
)

func ConnectCarrier(carrier ShipengineModels.CarrierConnect, carrierName string) (*ShipengineModels.CarrierConnectResponse, *ShipengineModels.ShipengineError, error) {

	body, err := json.Marshal(carrier.Carrier)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequest("POST", apiClient.GetApiHost()+"/connections/carriers/"+carrierName, bytes.NewBuffer(body))
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
		Endpoint:   "/connections/carriers/" + carrierName,
		StatusCode: resp.StatusCode,
		RequestURL: apiClient.GetApiHost() + "/connections/carriers/" + carrierName,
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

	var result *ShipengineModels.CarrierConnectResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, nil, err
	}

	return result, nil, nil
}

func DeleteCarrier(carrierName string, carrierId string) (bool, *ShipengineModels.ShipengineError, error) {

	req, err := http.NewRequest("DELETE", apiClient.GetApiHost()+"/connections/carriers/"+carrierName+"/"+carrierId, nil)
	if err != nil {
		return false, nil, err
	}

	resp, err := apiClient.Do(req)
	if err != nil {
		return false, nil, err
	}

	log := models.ShipEngineLog{
		EventType:  "api_request",
		Headers:    resp.Header,
		Endpoint:   "/connections/carriers/" + carrierName + "/" + carrierId,
		StatusCode: resp.StatusCode,
		RequestURL: apiClient.GetApiHost() + "/connections/carriers/" + carrierName + "/" + carrierId,
		Method:     "DELETE",
		CreatedAt:  time.Now(),
	}
	log.Create()

	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		var errorResponse ShipengineModels.ShipengineError
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		if err != nil {
			return false, nil, err
		}
		return false, &errorResponse, nil
	}

	return true, nil, nil
}

func GetCarrierConnectionOptions(carrier_id string) (*ShipengineModels.CarrierOptionsResponse, *ShipengineModels.ShipengineError, error) {

	req, err := http.NewRequest("GET", apiClient.GetApiHost()+"/carriers/"+carrier_id+"/options", nil)
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
		Endpoint:   "/carriers/" + carrier_id + "/options",
		StatusCode: resp.StatusCode,
		RequestURL: apiClient.GetApiHost() + "/carriers/" + carrier_id + "/options",
		Method:     "GET",
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

	var result *ShipengineModels.CarrierOptionsResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, nil, err
	}

	return result, nil, nil
}

func GetCarrierConnectionPackageTypes(carrier_id string) (*ShipengineModels.CarrierPackageTypesResponse, *ShipengineModels.ShipengineError, error) {

	req, err := http.NewRequest("GET", apiClient.GetApiHost()+"/carriers/"+carrier_id+"/packages", nil)
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
		Endpoint:   "/carriers/" + carrier_id + "/packages",
		StatusCode: resp.StatusCode,
		RequestURL: apiClient.GetApiHost() + "/carriers/" + carrier_id + "/packages",
		Method:     "GET",
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

	var result *ShipengineModels.CarrierPackageTypesResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, nil, err
	}

	return result, nil, nil
}

func GetCarrierConnectionServices(carrier_id string) (*ShipengineModels.CarrierServicesResponse, *ShipengineModels.ShipengineError, error) {

	req, err := http.NewRequest("GET", apiClient.GetApiHost()+"/carriers/"+carrier_id+"/services", nil)
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
		Endpoint:   "/carriers/" + carrier_id + "/services",
		StatusCode: resp.StatusCode,
		RequestURL: apiClient.GetApiHost() + "/carriers/" + carrier_id + "/services",
		Method:     "GET",
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

	var result *ShipengineModels.CarrierServicesResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, nil, err
	}

	return result, nil, nil
}
