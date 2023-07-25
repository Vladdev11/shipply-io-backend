package ShipengineAPI

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	ShipengineModels "github.com/shipply-io/shipply-io-backend/api/shipengine/models"
	"github.com/shipply-io/shipply-io-backend/models"
)

func ConnectCarrier(ctx context.Context, carrier ShipengineModels.CarrierConnect, carrierName string) (*ShipengineModels.CarrierConnectResponse, *ShipengineModels.ShipengineError, error) {

	body, err := json.Marshal(carrier.Carrier)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequest("POST", "/connections/carriers/"+carrierName, bytes.NewBuffer(body))
	if err != nil {
		return nil, nil, err
	}

	resp, err := ShipengineClientFromContext(ctx).Do(req)
	if err != nil {
		return nil, nil, err
	}

	log := models.ShipEngineLog{
		EventType:  "api_request",
		Headers:    resp.Header,
		Body:       string(body),
		Endpoint:   "/connections/carriers/" + carrierName,
		StatusCode: resp.StatusCode,
		RequestURL: resp.Request.URL.String(),
		Method:     "POST",
		CreatedAt:  time.Now(),
	}
	log.Create(ctx)

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

func DeleteCarrier(ctx context.Context, carrierName string, carrierId string) (bool, *ShipengineModels.ShipengineError, error) {

	req, err := http.NewRequest("DELETE", "/connections/carriers/"+carrierName+"/"+carrierId, nil)
	if err != nil {
		return false, nil, err
	}

	resp, err := ShipengineClientFromContext(ctx).Do(req)
	if err != nil {
		return false, nil, err
	}

	log := models.ShipEngineLog{
		EventType:  "api_request",
		Headers:    resp.Header,
		Endpoint:   "/connections/carriers/" + carrierName + "/" + carrierId,
		StatusCode: resp.StatusCode,
		RequestURL: resp.Request.URL.String(),
		Method:     "DELETE",
		CreatedAt:  time.Now(),
	}
	log.Create(ctx)

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

func GetCarrierConnectionOptions(ctx context.Context, carrier_id string) (*ShipengineModels.CarrierOptionsResponse, *ShipengineModels.ShipengineError, error) {

	req, err := http.NewRequest("GET", "/carriers/"+carrier_id+"/options", nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := ShipengineClientFromContext(ctx).Do(req)
	if err != nil {
		return nil, nil, err
	}

	log := models.ShipEngineLog{
		EventType:  "api_request",
		Headers:    resp.Header,
		Endpoint:   "/carriers/" + carrier_id + "/options",
		StatusCode: resp.StatusCode,
		RequestURL: resp.Request.URL.String(),
		Method:     "GET",
		CreatedAt:  time.Now(),
	}
	log.Create(ctx)

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

func GetCarrierConnectionPackageTypes(ctx context.Context, carrier_id string) (*ShipengineModels.CarrierPackageTypesResponse, *ShipengineModels.ShipengineError, error) {

	req, err := http.NewRequest("GET", "/carriers/"+carrier_id+"/packages", nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := ShipengineClientFromContext(ctx).Do(req)
	if err != nil {
		return nil, nil, err
	}

	log := models.ShipEngineLog{
		EventType:  "api_request",
		Headers:    resp.Header,
		Endpoint:   "/carriers/" + carrier_id + "/packages",
		StatusCode: resp.StatusCode,
		RequestURL: resp.Request.URL.String(),
		Method:     "GET",
		CreatedAt:  time.Now(),
	}
	log.Create(ctx)

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

func GetCarrierConnectionServices(ctx context.Context, carrier_id string) (*ShipengineModels.CarrierServicesResponse, *ShipengineModels.ShipengineError, error) {

	req, err := http.NewRequest("GET", "/carriers/"+carrier_id+"/services", nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := ShipengineClientFromContext(ctx).Do(req)
	if err != nil {
		return nil, nil, err
	}

	log := models.ShipEngineLog{
		EventType:  "api_request",
		Headers:    resp.Header,
		Endpoint:   "/carriers/" + carrier_id + "/services",
		StatusCode: resp.StatusCode,
		RequestURL: resp.Request.URL.String(),
		Method:     "GET",
		CreatedAt:  time.Now(),
	}
	log.Create(ctx)

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
