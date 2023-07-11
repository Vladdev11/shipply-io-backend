package ShipengineModels

import (
	"github.com/shipply-io/shipply-io-backend/models"

	"github.com/shipply-io/shipply-io-backend/util"
)

type Address struct {
	Name                        string `json:"name"`
	Email                       string `json:"email"`
	Phone                       string `json:"phone"`
	CompanyName                 string `json:"company_name"`
	AddressLine1                string `json:"address_line1"`
	AddressLine2                string `json:"address_line2"`
	AddressLine3                string `json:"address_line3"`
	CityLocality                string `json:"city_locality"`
	StateProvince               string `json:"state_province"`
	PostalCode                  string `json:"postal_code"`
	CountryCode                 string `json:"country_code"`
	AddressResidentialIndicator string `json:"address_residential_indicator,omitempty"`
}

type AddressResponse struct {
	Status          string                   `json:"status"` //this will be "verified" or "unverified"
	OriginalAddress Address                  `json:"original_address"`
	MatchedAddress  Address                  `json:"matched_address"`
	Messages        []AddressResponseMessage `json:"messages"`
}

type AddressResponseMessage struct {
	Code       string `json:"code"`
	DetailCode string `json:"detail_code"`
	Message    string `json:"message"`
	Type       string `json:"type"`
}

type AddressParseRequest struct {
	Text string `json:"text"`
}

type AddressParseResponse struct {
	Score    float64                      `json:"score"`
	Address  Address                      `json:"address"`
	Entities []AddressParseResponseEntity `json:"entities"`
}

type AddressParseResponseEntity struct {
	Type       string  `json:"type"`
	Score      float64 `json:"score"`
	Text       string  `json:"text"`
	StartIndex int     `json:"start_index"`
	EndIndex   int     `json:"end_index"`
	Result     struct {
		Value interface{} `json:"value"`
		Line  int         `json:"line"`
		Type  string      `json:"type"`
	} `json:"result"`
}

func ConvertAddressToShipengineAddress(address *models.Address) Address {
	return Address{
		Name:                        address.FirstName + " " + address.LastName,
		Phone:                       address.Phone,
		CompanyName:                 address.Company,
		AddressLine1:                address.Street1,
		AddressLine2:                address.Street2,
		AddressLine3:                address.Street3,
		CityLocality:                address.City,
		StateProvince:               address.State,
		PostalCode:                  address.PostalCode,
		CountryCode:                 address.Country,
		AddressResidentialIndicator: util.BooleanToYesNo(address.Residential),
	}
}

func ConvertShipengineAddressToModelAddress(address Address) models.Address {
	return models.Address{
		Company:     address.CompanyName,
		Street1:     address.AddressLine1,
		Street2:     address.AddressLine2,
		Street3:     address.AddressLine3,
		City:        address.CityLocality,
		State:       address.StateProvince,
		PostalCode:  address.PostalCode,
		Country:     address.CountryCode,
		Phone:       address.Phone,
		Residential: util.YesNoToBoolean(address.AddressResidentialIndicator),
	}
}
