package models

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type ShippingScanToteRequest struct {
	Barcode string `json:"barcode"`
}

type ShippingScanToteResponse struct {
	PickSessionOrderID int `json:"pick_session_order_id"`
}

type ShippingGetPickSessionOrderResponse struct {
	PickSessionOrderID      int  `json:"pick_session_order_id"`
	ToteID                  int  `json:"tote_id"`
	HasError                bool `json:"has_error"`
	PickSessionOrderErrorID *int `json:"pick_session_order_error_id"`
	OrderID                 int  `json:"order_id"`

	Rate                  ShippingRate                     `json:"rate,omitempty"`
	Order                 OrderReturnJSON                  `json:"order"`
	PickSessionOrderItems []PickSessionOrderItemReturnJSON `json:"pick_session_order_items"`
	Boxes                 []BoxReturnJSON                  `json:"boxes"`
}

type ShippingGetRatesRequest struct {
	BoxID  int     `json:"box_id"`
	Weight float64 `json:"weight"`
}

type ShippingGetRatesResponse struct {
	CarrierConnections []CarrierConnectionRatesResponse `json:"carrier_connections"`
}

type CarrierConnectionRatesResponse struct {
	CarrierConnectionID       int            `json:"carrier_connection_id"`
	CarrierConnectionNickName string         `json:"carrier_name"`
	Rates                     []ShippingRate `json:"rates"`
}

type ShippingSelectRateRequest struct {
	ShippingRateID int `json:"shipping_rate_id"`
}

type ShippingPurchaseLabelRequest struct {
	Weight float64 `json:"weight"`
	BoxID  int     `json:"box_id"`
}

func (sctr *ShippingScanToteRequest) ParseAndValidateRequest(r *http.Request) []string {

	errors := []string{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		Barcode json.RawMessage `json:"barcode"`
	}{}

	if err := json.Unmarshal(body, aux); err != nil {
		return []string{"invalid JSON"}
	}

	if aux.Barcode == nil {
		errors = append(errors, "barcode is required")
	} else if err := json.Unmarshal(aux.Barcode, &sctr.Barcode); err != nil {
		errors = append(errors, "barcode must be a string")
	} else if len(sctr.Barcode) < 1 {
		errors = append(errors, "barcode must be at least 1 character")
	}

	if len(errors) > 0 {
		return errors
	}

	return nil

}

func (sgrr *ShippingGetRatesRequest) ParseAndValidateRequest(r *http.Request) []string {

	errors := []string{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		BoxID  json.RawMessage `json:"box_id"`
		Weight json.RawMessage `json:"weight"`
	}{}

	if err := json.Unmarshal(body, aux); err != nil {
		return []string{"invalid JSON"}
	}

	user, err := GetRequestingUser(r)

	if aux.Weight == nil {
		errors = append(errors, "weight is required")
	} else if err := json.Unmarshal(aux.Weight, &sgrr.Weight); err != nil {
		errors = append(errors, "weight must be a number")
	} else if sgrr.Weight <= 0 {
		errors = append(errors, "weight must be greater than 0")
	}

	if aux.BoxID == nil {
		errors = append(errors, "box_id is required")
	} else if err := json.Unmarshal(aux.BoxID, &sgrr.BoxID); err != nil {
		errors = append(errors, "box_id must be an integer")
	} else if sgrr.BoxID < 1 {
		errors = append(errors, "box_id must be greater than 0")
	} else {
		box, err := GetBoxByID(sgrr.BoxID)
		if err != nil {
			errors = append(errors, "invalid box_id")
		}

		if box.OrganizationID != user.OwnerID {
			errors = append(errors, "invalid box_id")
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil

}

func (ssrr *ShippingSelectRateRequest) ParseAndValidateRequest(r *http.Request) []string {

	errors := []string{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		ShippingRateID json.RawMessage `json:"shipping_rate_id"`
	}{}

	if err := json.Unmarshal(body, aux); err != nil {
		return []string{"invalid JSON"}
	}

	if aux.ShippingRateID == nil {
		errors = append(errors, "shipping_rate_id is required")
	} else if err := json.Unmarshal(aux.ShippingRateID, &ssrr.ShippingRateID); err != nil {
		errors = append(errors, "shipping_rate_id must be an integer")
	} else if ssrr.ShippingRateID < 1 {
		errors = append(errors, "shipping_rate_id must be greater than 0")
	}

	if len(errors) > 0 {
		return errors
	}

	return nil

}

func (splr *ShippingPurchaseLabelRequest) ParseAndValidateRequest(r *http.Request) []string {

	errors := []string{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		RateID json.RawMessage `json:"rate_id"`
		Weight json.RawMessage `json:"weight"`
		BoxID  json.RawMessage `json:"box_id"`
	}{}

	if err := json.Unmarshal(body, aux); err != nil {
		return []string{"invalid JSON"}
	}

	if aux.Weight == nil {
		errors = append(errors, "weight is required")
	} else if err := json.Unmarshal(aux.Weight, &splr.Weight); err != nil {
		errors = append(errors, "weight must be a number")
	} else if splr.Weight < 1 {
		errors = append(errors, "weight must be greater than 0")
	}

	if aux.BoxID == nil {
		errors = append(errors, "box_id is required")
	} else if err := json.Unmarshal(aux.BoxID, &splr.BoxID); err != nil {
		errors = append(errors, "box_id must be an integer")
	} else if splr.BoxID < 1 {
		errors = append(errors, "box_id must be greater than 0")
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}
