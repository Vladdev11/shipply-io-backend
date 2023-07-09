package models

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"
)

type CarrierConnection struct {
	ID                  int
	CarrierID           int
	OwnerID             int
	OwnerType           int // 1 = Organization, 2 = Client
	Active              bool
	ShipengineCarrierID string
	ShipengineNickname  string
	CarrierOptions      json.RawMessage `gorm:"type:jsonb"` //this is []CarrierOption
	CarrierServices     json.RawMessage `gorm:"type:jsonb"` //this is []CarrierService
	CarrierPackageTypes json.RawMessage `gorm:"type:jsonb"` //this is []CarrierPackageType
	Credentials         json.RawMessage `gorm:"type:jsonb"`
	Settings            json.RawMessage `gorm:"type:jsonb"` //we use the settings field to store the settings as json in the database

	ParsedSettings CarrierConnectionSettings `gorm:"-"` //this is the parsed settings field as a struct

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	Carrier Carrier
}

type CarrierConnectionSettings struct {
	GeneralSettings         struct{}                          //Use this to store general settings for the carrier
	CarrierSpecificSettings struct{}                          //Use this to store carrier specific settings for the carrier
	EnabledServices         []CarrierConnectionEnabledService //Use this to store the enabled services for the carrier
}

type CarrierConnectionEnabledService struct {
	ServiceCode string `json:"service_code"`
	ServiceName string `json:"service_name"`
	Enabled     bool   `json:"enabled"`
}

func (ccs *CarrierConnectionSettings) ConvertToReturnJSON() interface{} {
	return map[string]interface{}{
		"general_settings":          ccs.GeneralSettings,
		"carrier_specific_settings": ccs.CarrierSpecificSettings,
		"enabled_services":          ccs.EnabledServices,
	}
}

type CarrierOption struct {
	Name         string      `json:"name"`
	DefaultValue interface{} `json:"default_value"`
	Description  string      `json:"description"`
}

type CarrierPackageType struct {
	PackageID   string `json:"package_id"`
	PackageCode string `json:"package_code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CarrierService struct {
	CarrierID               string `json:"carrier_id"`
	CarrierCode             string `json:"carrier_code"`
	ServiceCode             string `json:"service_code"`
	Name                    string `json:"name"`
	Domestic                bool   `json:"domestic"`
	International           bool   `json:"international"`
	IsMultiPackageSupported bool   `json:"is_multi_package_supported"`
}

type CarrierConnectionReturnJSON struct {
	ID        int               `json:"id"`
	CarrierID int               `json:"carrier_id"`
	OwnerID   int               `json:"owner_id"`
	OwnerType int               `json:"owner_type"`
	Active    bool              `json:"active"`
	Settings  interface{}       `json:"settings"`
	Carrier   CarrierReturnJSON `json:"carrier,omitempty"`
}

func (cc *CarrierConnection) ConvertToReturnJSON() CarrierConnectionReturnJSON {
	return CarrierConnectionReturnJSON{
		ID:        cc.ID,
		CarrierID: cc.CarrierID,
		OwnerID:   cc.OwnerID,
		OwnerType: cc.OwnerType,
		Active:    cc.Active,
		Settings:  cc.ParsedSettings.ConvertToReturnJSON(),
		Carrier:   cc.Carrier.ConvertToReturnJSON(),
	}
}

func (cc *CarrierConnection) SaveCarrierOptions(ctx context.Context, carrierOptions []CarrierOption) error {

	var err error
	cc.CarrierOptions, err = util.ConvertToJSONRawMessage(carrierOptions)
	if err != nil {
		return err
	}
	//save only the carrier options field in the database
	err = util.DBFromContext(ctx).Model(cc).Update("carrier_options", cc.CarrierOptions).Error
	if err != nil {
		return err
	}
	return nil
}

func (cc *CarrierConnection) GetCarrier(ctx context.Context) error {
	var carrier Carrier
	err := util.DBFromContext(ctx).First(&carrier, cc.CarrierID).Error
	if err != nil {
		return err
	}
	cc.Carrier = carrier
	return nil
}

func (cc *CarrierConnection) SaveCarrierPackageTypes(ctx context.Context, carrierPackageTypes []CarrierPackageType) error {
	var err error
	cc.CarrierPackageTypes, err = util.ConvertToJSONRawMessage(carrierPackageTypes)
	if err != nil {
		return err
	}

	//save only the carrier package types field in the database
	err = util.DBFromContext(ctx).Model(cc).Update("carrier_package_types", cc.CarrierPackageTypes).Error
	if err != nil {
		return err
	}
	return nil
}

func (cc *CarrierConnection) SaveCarrierServices(ctx context.Context, carrierServices []CarrierService) error {
	var err error
	cc.CarrierServices, err = util.ConvertToJSONRawMessage(carrierServices)
	if err != nil {
		return err
	}

	//save only the carrier services field in the database
	err = util.DBFromContext(ctx).Model(cc).Update("carrier_services", cc.CarrierServices).Error
	if err != nil {
		return err
	}

	return nil
}

func GetCarrierConnectionByID(ctx context.Context, id int) (*CarrierConnection, error) {
	var carrierConnection CarrierConnection

	err := util.DBFromContext(ctx).Preload("Carrier").First(&carrierConnection, id).Error
	if err != nil {
		return nil, err
	}

	return &carrierConnection, nil
}

func CreateCarrierConnection(carrier *Carrier, r *http.Request) (*CarrierConnect, []string) {

	var errs []string
	var carrierConnect CarrierConnect

	//TODO make sure exisitng account number does not already exist so we don't get duplicate accounts

	switch carrier.ShipEngineID {
	case "asendia":
		carrierConnect = &CarrierConnectAsendia{}
		errs = carrierConnect.ParseAndValidateRequest(r)
	case "endicia":
		carrierConnect = &CarrierConnectEndicia{}
		errs = carrierConnect.ParseAndValidateRequest(r)
	case "dhl_ecommerce":
		carrierConnect = &CarrierConnectDHLeCommerce{}
		errs = carrierConnect.ParseAndValidateRequest(r)
	case "dhl_express":
		carrierConnect = &CarrierConnectDHLExpress{}
		errs = carrierConnect.ParseAndValidateRequest(r)
	case "ontrac":
		carrierConnect = &CarrierConnectOnTrac{}
		errs = carrierConnect.ParseAndValidateRequest(r)
	case "stamps_com":
		carrierConnect = &CarrierConnectStampsCom{}
		errs = carrierConnect.ParseAndValidateRequest(r)
	case "ups":
		carrierConnect = &CarrierConnectUPS{}
		errs = carrierConnect.ParseAndValidateRequest(r)
	case "fedex":
		carrierConnect = &CarrierConnectFedexUSCA{}
		errs = carrierConnect.ParseAndValidateRequest(r)
	default:
		return nil, []string{"carrier not supported"}
	}

	if len(errs) > 0 {
		return nil, errs
	}

	return &carrierConnect, nil

}

func (cc *CarrierConnection) Create(ctx context.Context) error {
	err := util.DBFromContext(ctx).Create(cc).Error
	if err != nil {
		return err
	}

	return nil

}

func (cc *CarrierConnection) Update(ctx context.Context) error {
	err := util.DBFromContext(ctx).Save(cc).Error
	if err != nil {
		return err
	}

	return nil

}

func (cc *CarrierConnection) Delete(ctx context.Context) error {
	err := util.DBFromContext(ctx).Delete(cc).Error
	if err != nil {
		return err
	}

	return nil
}

type CarrierConnect interface {
	ParseAndValidateRequest(r *http.Request) []string
	GetNickname() (string, error)
}

type CarrierConnectAsendia struct {
	Nickname      string `json:"nickname"`
	FTPUsername   string `json:"ftp_username"`
	FTPPassword   string `json:"ftp_password"`
	AccountNumber int    `json:"account_number"`
}

func (ccar *CarrierConnectAsendia) ParseAndValidateRequest(r *http.Request) []string {

	errs := []string{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		Nickname      json.RawMessage `json:"nickname"`
		FTPUsername   json.RawMessage `json:"ftp_username"`
		FTPPassword   json.RawMessage `json:"ftp_password"`
		AccountNumber json.RawMessage `json:"account_number"`
	}{}

	if err := json.Unmarshal(body, &aux); err != nil {
		errs = append(errs, "invalid json")
	}

	if aux.Nickname != nil {
		if err := json.Unmarshal(aux.Nickname, &ccar.Nickname); err != nil {
			errs = append(errs, "nickname must be string")
		}
	} else {
		errs = append(errs, "nickname is required")
	}

	if aux.FTPUsername != nil {
		if err := json.Unmarshal(aux.FTPUsername, &ccar.FTPUsername); err != nil {
			errs = append(errs, "ftp_username must be string")
		}
	} else {
		errs = append(errs, "ftp_username is required")
	}

	if aux.FTPPassword != nil {
		if err := json.Unmarshal(aux.FTPPassword, &ccar.FTPPassword); err != nil {
			errs = append(errs, "ftp_password must be string")
		}
	} else {
		errs = append(errs, "ftp_password is required")
	}

	if aux.AccountNumber != nil {
		if err := json.Unmarshal(aux.AccountNumber, &ccar.AccountNumber); err != nil {
			errs = append(errs, "account_number must be int")
		}
	} else {
		errs = append(errs, "account_number is required")
	}

	if len(errs) > 0 {
		return errs
	}

	return nil

}

func (ccar *CarrierConnectAsendia) GetNickname() (string, error) {
	return ccar.Nickname, nil
}

type CarrierConnectEndicia struct {
	Nickname   string `json:"nickname"`
	Account    string `json:"account"`
	Passphrase string `json:"passphrase"`
}

func (cceu *CarrierConnectEndicia) ParseAndValidateRequest(r *http.Request) []string {

	errs := []string{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		Nickname   json.RawMessage `json:"nickname"`
		Account    json.RawMessage `json:"account"`
		Passphrase json.RawMessage `json:"passphrase"`
	}{}

	if err := json.Unmarshal(body, &aux); err != nil {
		errs = append(errs, "invalid json")
	}

	if aux.Nickname != nil {
		if err := json.Unmarshal(aux.Nickname, &cceu.Nickname); err != nil {
			errs = append(errs, "nickname must be string")
		}
	} else {
		errs = append(errs, "nickname is required")
	}

	if aux.Account != nil {
		if err := json.Unmarshal(aux.Account, &cceu.Account); err != nil {
			errs = append(errs, "account must be string")
		}
	} else {
		errs = append(errs, "account is required")
	}

	if aux.Passphrase != nil {
		if err := json.Unmarshal(aux.Passphrase, &cceu.Passphrase); err != nil {
			errs = append(errs, "passphrase must be string")
		}
	} else {
		errs = append(errs, "passphrase is required")
	}

	if len(errs) > 0 {
		return errs
	}

	return nil

}

func (cceu *CarrierConnectEndicia) GetNickname() (string, error) {
	return cceu.Nickname, nil
}

type CarrierConnectDHLeCommerce struct {
	Nickname           string `json:"nickname"`
	ClientID           string `json:"client_id"`
	APISecret          string `json:"api_secret"`
	PickupNumber       string `json:"pickup_number"`
	DistributionCenter string `json:"distribution_center"`
	SoldTo             string `json:"sold_to"`
}

func (ccdhl *CarrierConnectDHLeCommerce) ParseAndValidateRequest(r *http.Request) []string {

	errs := []string{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		Nickname           json.RawMessage `json:"nickname"`
		ClientID           json.RawMessage `json:"client_id"`
		APISecret          json.RawMessage `json:"api_secret"`
		PickupNumber       json.RawMessage `json:"pickup_number"`
		DistributionCenter json.RawMessage `json:"distribution_center"`
		SoldTo             json.RawMessage `json:"sold_to"`
	}{}

	if err := json.Unmarshal(body, &aux); err != nil {
		errs = append(errs, "invalid json")
	}

	if aux.Nickname != nil {
		if err := json.Unmarshal(aux.Nickname, &ccdhl.Nickname); err != nil {
			errs = append(errs, "nickname must be string")
		}
	} else {
		errs = append(errs, "nickname is required")
	}

	if aux.ClientID != nil {
		if err := json.Unmarshal(aux.ClientID, &ccdhl.ClientID); err != nil {
			errs = append(errs, "client_id must be string")
		}
	} else {
		errs = append(errs, "client_id is required")
	}

	if aux.APISecret != nil {
		if err := json.Unmarshal(aux.APISecret, &ccdhl.APISecret); err != nil {
			errs = append(errs, "api_secret must be string")
		}
	} else {
		errs = append(errs, "api_secret is required")
	}

	if aux.PickupNumber != nil {
		if err := json.Unmarshal(aux.PickupNumber, &ccdhl.PickupNumber); err != nil {
			errs = append(errs, "pickup_number must be string")
		}
	} else {
		errs = append(errs, "pickup_number is required")
	}

	if aux.DistributionCenter != nil {
		if err := json.Unmarshal(aux.DistributionCenter, &ccdhl.DistributionCenter); err != nil {
			errs = append(errs, "distribution_center must be string")
		}
	} else {
		errs = append(errs, "distribution_center is required")
	}

	if aux.SoldTo != nil {
		if err := json.Unmarshal(aux.SoldTo, &ccdhl.SoldTo); err != nil {
			errs = append(errs, "sold_to must be string")
		}
	} else {
		errs = append(errs, "sold_to is required")
	}

	if len(errs) > 0 {
		return errs
	}

	return nil

}

func (ccdhl *CarrierConnectDHLeCommerce) GetNickname() (string, error) {
	return ccdhl.Nickname, nil
}

type CarrierConnectDHLExpress struct {
	Nickname      string `json:"nickname"`
	AccountNumber string `json:"account_number"`
	SiteID        string `json:"site_id,omitempty"`
	Password      string `json:"password,omitempty"`
	CountryCode   string `json:"country_code,omitempty"`
}

func (ccde *CarrierConnectDHLExpress) ParseAndValidateRequest(r *http.Request) []string {

	errs := []string{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		Nickname      json.RawMessage `json:"nickname"`
		AccountNumber json.RawMessage `json:"account_number"`
		SiteID        json.RawMessage `json:"site_id,omitempty"`
		Password      json.RawMessage `json:"password,omitempty"`
		CountryCode   json.RawMessage `json:"country_code,omitempty"`
	}{}

	if err := json.Unmarshal(body, &aux); err != nil {
		errs = append(errs, "invalid json")
	}

	if aux.Nickname != nil {
		if err := json.Unmarshal(aux.Nickname, &ccde.Nickname); err != nil {
			errs = append(errs, "nickname must be string")
		}
	} else {
		errs = append(errs, "nickname is required")
	}

	if aux.AccountNumber != nil {
		if err := json.Unmarshal(aux.AccountNumber, &ccde.AccountNumber); err != nil {
			errs = append(errs, "account_number must be string")
		}
	} else {
		errs = append(errs, "account_number is required")
	}

	if aux.SiteID != nil {
		if err := json.Unmarshal(aux.SiteID, &ccde.SiteID); err != nil {
			errs = append(errs, "site_id must be string")
		}
	}

	if aux.Password != nil {
		if err := json.Unmarshal(aux.Password, &ccde.Password); err != nil {
			errs = append(errs, "password must be string")
		}
	}

	if aux.CountryCode != nil {
		if err := json.Unmarshal(aux.CountryCode, &ccde.CountryCode); err != nil {
			errs = append(errs, "country_code must be string")
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil

}

func (ccde *CarrierConnectDHLExpress) GetNickname() (string, error) {
	return ccde.Nickname, nil
}

type CarrierConnectOnTrac struct {
	Nickname      string `json:"nickname"`
	AccountNumber int    `json:"account_number"`
	Password      string `json:"password"`
}

func (ccot *CarrierConnectOnTrac) ParseAndValidateRequest(r *http.Request) []string {

	errs := []string{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		Nickname      json.RawMessage `json:"nickname"`
		AccountNumber json.RawMessage `json:"account_number"`
		Password      json.RawMessage `json:"password"`
	}{}

	if err := json.Unmarshal(body, &aux); err != nil {
		errs = append(errs, "invalid json")
	}

	if aux.Nickname != nil {
		if err := json.Unmarshal(aux.Nickname, &ccot.Nickname); err != nil {
			errs = append(errs, "nickname must be string")
		}
	} else {
		errs = append(errs, "nickname is required")
	}

	if aux.AccountNumber != nil {
		if err := json.Unmarshal(aux.AccountNumber, &ccot.AccountNumber); err != nil {
			errs = append(errs, "account_number must be int")
		}
	} else {
		errs = append(errs, "account_number is required")
	}

	if aux.Password != nil {
		if err := json.Unmarshal(aux.Password, &ccot.Password); err != nil {
			errs = append(errs, "password must be string")
		}
	} else {
		errs = append(errs, "password is required")
	}

	if len(errs) > 0 {
		return errs
	}

	return nil

}

func (ccot *CarrierConnectOnTrac) GetNickname() (string, error) {
	return ccot.Nickname, nil
}

type CarrierConnectStampsCom struct {
	Nickname string `json:"nickname"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (ccsc *CarrierConnectStampsCom) ParseAndValidateRequest(r *http.Request) []string {

	errs := []string{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		Nickname json.RawMessage `json:"nickname"`
		Username json.RawMessage `json:"username"`
		Password json.RawMessage `json:"password"`
	}{}

	if err := json.Unmarshal(body, &aux); err != nil {
		errs = append(errs, "invalid json")
	}

	if aux.Nickname != nil {
		if err := json.Unmarshal(aux.Nickname, &ccsc.Nickname); err != nil {
			errs = append(errs, "nickname must be string")
		}
	} else {
		errs = append(errs, "nickname is required")
	}

	if aux.Username != nil {
		if err := json.Unmarshal(aux.Username, &ccsc.Username); err != nil {
			errs = append(errs, "username must be string")
		}
	} else {
		errs = append(errs, "username is required")
	}

	if aux.Password != nil {
		if err := json.Unmarshal(aux.Password, &ccsc.Password); err != nil {
			errs = append(errs, "password must be string")
		}
	} else {
		errs = append(errs, "password is required")
	}

	if len(errs) > 0 {
		return errs
	}

	return nil

}

func (ccsc *CarrierConnectStampsCom) GetNickname() (string, error) {
	return ccsc.Nickname, nil
}

type CarrierConnectUPS struct {
	Nickname             string                    `json:"nickname"`
	AccountNumber        string                    `json:"account_number"`
	AccountCountryCode   string                    `json:"account_country_code"`
	AccountPostalCode    string                    `json:"account_postal_code"`
	Title                string                    `json:"title"`
	FirstName            string                    `json:"first_name"`
	LastName             string                    `json:"last_name"`
	Company              string                    `json:"company,omitempty"`
	Address1             string                    `json:"address1"`
	Address2             string                    `json:"address2,omitempty"`
	City                 string                    `json:"city"`
	State                string                    `json:"state"`
	PostalCode           string                    `json:"postal_code"`
	CountryCode          string                    `json:"country_code"`
	Email                string                    `json:"email"`
	Phone                string                    `json:"phone"`
	Invoice              *CarrierConnectUPSInvoice `json:"invoice,omitempty"`
	AgreeToTechAgreement bool                      `json:"agree_to_technology_agreement"`
}

// use this struct to give front-end the proper fields to fill out
type CarrierConnectUPSFrontend struct {
	Nickname             string  `json:"nickname"`
	AccountNumber        string  `json:"account_number"`
	AccountCountryCode   string  `json:"account_country_code"`
	AccountPostalCode    string  `json:"account_postal_code"`
	Title                string  `json:"title"`
	FirstName            string  `json:"first_name"`
	LastName             string  `json:"last_name"`
	Company              string  `json:"company,omitempty"`
	Address1             string  `json:"address1"`
	Address2             string  `json:"address2,omitempty"`
	City                 string  `json:"city"`
	State                string  `json:"state"`
	PostalCode           string  `json:"postal_code"`
	CountryCode          string  `json:"country_code"`
	Email                string  `json:"email"`
	Phone                string  `json:"phone"`
	InvoiceControlID     string  `json:"invoice_control_id,omitempty"`
	InvoiceNumber        string  `json:"invoice_number,omitempty"`
	InvoiceAmount        float64 `json:"invoice_amount,omitempty"`
	InvoiceDate          string  `json:"invoice_date,omitempty"`
	AgreeToTechAgreement bool    `json:"agree_to_technology_agreement"`
}

type CarrierConnectUPSInvoice struct {
	ControlID     string  `json:"control_id"`
	InvoiceNumber string  `json:"invoice_number"`
	InvoiceAmount float64 `json:"invoice_amount"`
	InvoiceDate   string  `json:"invoice_date"`
}

func (ccups *CarrierConnectUPS) ParseAndValidateRequest(r *http.Request) []string {

	errs := []string{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		Nickname             json.RawMessage `json:"nickname"`
		AccountNumber        json.RawMessage `json:"account_number"`
		AccountCountryCode   json.RawMessage `json:"account_country_code"`
		AccountPostalCode    json.RawMessage `json:"account_postal_code"`
		Title                json.RawMessage `json:"title"`
		FirstName            json.RawMessage `json:"first_name"`
		LastName             json.RawMessage `json:"last_name"`
		Company              json.RawMessage `json:"company,omitempty"`
		Address1             json.RawMessage `json:"address1"`
		Address2             json.RawMessage `json:"address2,omitempty"`
		City                 json.RawMessage `json:"city"`
		State                json.RawMessage `json:"state"`
		PostalCode           json.RawMessage `json:"postal_code"`
		CountryCode          json.RawMessage `json:"country_code"`
		Email                json.RawMessage `json:"email"`
		Phone                json.RawMessage `json:"phone"`
		InvoiceControlID     json.RawMessage `json:"invoice_control_id,omitempty"`
		InvoiceNumber        json.RawMessage `json:"invoice_number,omitempty"`
		InvoiceAmount        json.RawMessage `json:"invoice_amount,omitempty"`
		InvoiceDate          json.RawMessage `json:"invoice_date,omitempty"`
		AgreeToTechAgreement json.RawMessage `json:"agree_to_technology_agreement"`
	}{}

	if err := json.Unmarshal(body, &aux); err != nil {
		errs = append(errs, "invalid json")
	}

	if aux.Nickname != nil {
		if err := json.Unmarshal(aux.Nickname, &ccups.Nickname); err != nil {
			errs = append(errs, "nickname must be string")
		}
	} else {
		errs = append(errs, "nickname is required")
	}

	if aux.AccountNumber != nil {
		if err := json.Unmarshal(aux.AccountNumber, &ccups.AccountNumber); err != nil {
			errs = append(errs, "account_number must be string")
		}
	} else {
		errs = append(errs, "account_number is required")
	}

	if aux.AccountCountryCode != nil {
		if err := json.Unmarshal(aux.AccountCountryCode, &ccups.AccountCountryCode); err != nil {
			errs = append(errs, "account_country_code must be string")
		}
	} else {
		errs = append(errs, "account_country_code is required")
	}

	if aux.AccountPostalCode != nil {
		if err := json.Unmarshal(aux.AccountPostalCode, &ccups.AccountPostalCode); err != nil {
			errs = append(errs, "account_postal_code must be string")
		}
	} else {
		errs = append(errs, "account_postal_code is required")
	}

	if aux.Title != nil {
		if err := json.Unmarshal(aux.Title, &ccups.Title); err != nil {
			errs = append(errs, "title must be string")
		}
	} else {
		errs = append(errs, "title is required")
	}

	if aux.FirstName != nil {
		if err := json.Unmarshal(aux.FirstName, &ccups.FirstName); err != nil {
			errs = append(errs, "first_name must be string")
		}
	} else {
		errs = append(errs, "first_name is required")
	}

	if aux.LastName != nil {
		if err := json.Unmarshal(aux.LastName, &ccups.LastName); err != nil {
			errs = append(errs, "last_name must be string")
		}
	} else {
		errs = append(errs, "last_name is required")
	}

	if aux.Address1 != nil {
		if err := json.Unmarshal(aux.Address1, &ccups.Address1); err != nil {
			errs = append(errs, "address1 must be string")
		}
	} else {
		errs = append(errs, "address1 is required")
	}

	if aux.City != nil {
		if err := json.Unmarshal(aux.City, &ccups.City); err != nil {
			errs = append(errs, "city must be string")
		}
	} else {
		errs = append(errs, "city is required")
	}

	if aux.State != nil {
		if err := json.Unmarshal(aux.State, &ccups.State); err != nil {
			errs = append(errs, "state must be string")
		}
	} else {
		errs = append(errs, "state is required")
	}
	if aux.PostalCode != nil {
		if err := json.Unmarshal(aux.PostalCode, &ccups.PostalCode); err != nil {
			errs = append(errs, "postal_code must be string")
		}
	} else {
		errs = append(errs, "postal_code is required")
	}

	if aux.CountryCode != nil {
		if err := json.Unmarshal(aux.CountryCode, &ccups.CountryCode); err != nil {
			errs = append(errs, "country_code must be string")
		}
	} else {
		errs = append(errs, "country_code is required")
	}

	if aux.Email != nil {
		if err := json.Unmarshal(aux.Email, &ccups.Email); err != nil {
			errs = append(errs, "email must be string")
		}
	} else {
		errs = append(errs, "email is required")
	}

	if aux.Phone != nil {
		if err := json.Unmarshal(aux.Phone, &ccups.Phone); err != nil {
			errs = append(errs, "phone must be string")
		}
	} else {
		errs = append(errs, "phone is required")
	}

	if aux.AgreeToTechAgreement != nil {
		if err := json.Unmarshal(aux.AgreeToTechAgreement, &ccups.AgreeToTechAgreement); err != nil {
			errs = append(errs, "agree_to_technology_agreement must be bool")
		}
	} else {
		errs = append(errs, "agree_to_technology_agreement is required")
	}

	// Parse and validate optional fields
	if aux.Company != nil {
		if err := json.Unmarshal(aux.Company, &ccups.Company); err != nil {
			errs = append(errs, "company must be string")
		}
	}

	if aux.Address2 != nil {
		if err := json.Unmarshal(aux.Address2, &ccups.Address2); err != nil {
			errs = append(errs, "address2 must be string")
		}
	}

	if aux.InvoiceControlID != nil {
		if err := json.Unmarshal(aux.InvoiceControlID, &ccups.Invoice.ControlID); err != nil {
			errs = append(errs, "invoice_control_id must be string")
		}
	}

	if aux.InvoiceNumber != nil {
		if err := json.Unmarshal(aux.InvoiceNumber, &ccups.Invoice.InvoiceNumber); err != nil {
			errs = append(errs, "invoice_control_number must be string")
		}
	}

	if aux.InvoiceAmount != nil {
		if err := json.Unmarshal(aux.InvoiceAmount, &ccups.Invoice.InvoiceAmount); err != nil {
			errs = append(errs, "invoice_control_amount must be string")
		}
	}

	if aux.InvoiceDate != nil {
		if err := json.Unmarshal(aux.InvoiceDate, &ccups.Invoice.InvoiceDate); err != nil {
			errs = append(errs, "invoice_control_type must be string")
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil

}

func (ccsc *CarrierConnectUPS) GetNickname() (string, error) {
	return ccsc.Nickname, nil
}

func GetAllCarrierConnections(ctx context.Context) ([]CarrierConnection, error) {
	var carrierConnections []CarrierConnection
	err := util.DBFromContext(ctx).Find(&carrierConnections).Error
	if err != nil {
		return nil, err
	}

	return carrierConnections, nil

}

func (cc *CarrierConnection) GetCarrierServices() ([]CarrierService, error) {
	var carrierServices []CarrierService

	//parse cc.CarrierServices into a slice of carrier services
	err := json.Unmarshal(cc.CarrierServices, &carrierServices)
	if err != nil {
		return nil, err
	}

	return carrierServices, nil
}

func (cc *CarrierConnection) GetCarrierOptions() ([]CarrierOption, error) {
	var carrierOptions []CarrierOption

	//parse cc.CarrierOptions into a slice of carrier options
	err := json.Unmarshal(cc.CarrierOptions, &carrierOptions)
	if err != nil {
		return nil, err
	}

	return carrierOptions, nil
}

func (cc *CarrierConnection) GetSettings() error {

	var settings CarrierConnectionSettings

	if cc.Settings == nil {
		settings = CarrierConnectionSettings{}
	}

	carrierServices, err := cc.GetCarrierServices()
	if err != nil {
		return err
	}

	if settings.EnabledServices == nil {
		settings.EnabledServices = []CarrierConnectionEnabledService{}
	}

	for _, carrierService := range carrierServices {
		found := false
		for i, enabledService := range settings.EnabledServices {
			if enabledService.ServiceCode == carrierService.ServiceCode {
				found = true
				settings.EnabledServices[i].Enabled = true
				break
			}
		}

		if !found {
			enabledService := CarrierConnectionEnabledService{
				ServiceCode: carrierService.CarrierCode,
				ServiceName: carrierService.Name,
				Enabled:     true,
			}
			settings.EnabledServices = append(settings.EnabledServices, enabledService)
		}
	}

	cc.ParsedSettings = settings
	return nil
}

func GetCarrierConnectionsByClientID(ctx context.Context, clientID int, orgID int) ([]CarrierConnection, error) {
	var carrierConnections []CarrierConnection

	err := util.DBFromContext(ctx).Where("owner_id = ? AND owner_type = ?", clientID, 2).Or("owner_id = ? AND owner_type = ?", orgID, 1).Find(&carrierConnections).Error
	if err != nil {
		return nil, err
	}

	return carrierConnections, nil

}

func GetCheapestCarrierConnection(ctx context.Context) (*CarrierConnection, error) {

	var carrierConnections CarrierConnection

	err := util.DBFromContext(ctx).Where("id = 1000").Find(&carrierConnections).Error

	if err != nil {
		return nil, err
	}

	return &carrierConnections, nil

}

type CarrierConnectFedexUSCA struct {
	Nickname      string `json:"nickname"`
	AccountNumber string `json:"account_number"`
	Company       string `json:"company"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Phone         string `json:"phone"`
	Address1      string `json:"address1"`
	Address2      string `json:"address2"`
	City          string `json:"city"`
	State         string `json:"state"`
	PostalCode    string `json:"postal_code"`
	CountryCode   string `json:"country_code"`
	Email         string `json:"email"`
	AgreeToEula   bool   `json:"agree_to_eula"`
}

func (ccfusca *CarrierConnectFedexUSCA) ParseAndValidateRequest(r *http.Request) []string {

	errs := []string{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		Nickname      json.RawMessage `json:"nickname"`
		AccountNumber json.RawMessage `json:"account_number"`
		Company       json.RawMessage `json:"company"`
		FirstName     json.RawMessage `json:"first_name"`
		LastName      json.RawMessage `json:"last_name"`
		Phone         json.RawMessage `json:"phone"`
		Address1      json.RawMessage `json:"address1"`
		Address2      json.RawMessage `json:"address2"`
		City          json.RawMessage `json:"city"`
		State         json.RawMessage `json:"state"`
		PostalCode    json.RawMessage `json:"postal_code"`
		CountryCode   json.RawMessage `json:"country_code"`
		Email         json.RawMessage `json:"email"`
		AgreeToEula   json.RawMessage `json:"agree_to_eula"`
	}{}

	if err := json.Unmarshal(body, &aux); err != nil {
		errs = append(errs, "invalid json")
	}

	if aux.Nickname != nil {
		if err := json.Unmarshal(aux.Nickname, &ccfusca.Nickname); err != nil {
			errs = append(errs, "nickname must be string")
		}
	} else {
		errs = append(errs, "nickname is required")
	}

	if aux.AccountNumber != nil {
		if err := json.Unmarshal(aux.AccountNumber, &ccfusca.AccountNumber); err != nil {
			errs = append(errs, "account_number must be string")
		}
	} else {
		errs = append(errs, "account_number is required")
	}

	if aux.Company != nil {
		if err := json.Unmarshal(aux.Company, &ccfusca.Company); err != nil {
			errs = append(errs, "company must be string")
		}
	} else {
		errs = append(errs, "company is required")
	}

	if aux.FirstName != nil {
		if err := json.Unmarshal(aux.FirstName, &ccfusca.FirstName); err != nil {
			errs = append(errs, "first_name must be string")
		}
	} else {
		errs = append(errs, "first_name is required")
	}

	if aux.LastName != nil {
		if err := json.Unmarshal(aux.LastName, &ccfusca.LastName); err != nil {
			errs = append(errs, "last_name must be string")
		}
	} else {
		errs = append(errs, "last_name is required")
	}

	if aux.Phone != nil {
		if err := json.Unmarshal(aux.Phone, &ccfusca.Phone); err != nil {
			errs = append(errs, "phone must be string")
		}
	} else {
		errs = append(errs, "phone is required")
	}

	if aux.Address1 != nil {
		if err := json.Unmarshal(aux.Address1, &ccfusca.Address1); err != nil {
			errs = append(errs, "address1 must be string")
		}
	} else {
		errs = append(errs, "address1 is required")
	}

	if aux.Address2 != nil {
		if err := json.Unmarshal(aux.Address2, &ccfusca.Address2); err != nil {
			errs = append(errs, "address2 must be string")
		}
	}

	if aux.City != nil {
		if err := json.Unmarshal(aux.City, &ccfusca.City); err != nil {
			errs = append(errs, "city must be string")
		}
	} else {
		errs = append(errs, "city is required")
	}

	if aux.State != nil {
		if err := json.Unmarshal(aux.State, &ccfusca.State); err != nil {
			errs = append(errs, "state must be string")
		}
	} else {
		errs = append(errs, "state is required")
	}

	if aux.PostalCode != nil {
		if err := json.Unmarshal(aux.PostalCode, &ccfusca.PostalCode); err != nil {
			errs = append(errs, "postal_code must be string")
		}
	} else {
		errs = append(errs, "postal_code is required")
	}

	if aux.CountryCode != nil {
		if err := json.Unmarshal(aux.CountryCode, &ccfusca.CountryCode); err != nil {
			errs = append(errs, "country_code must be string")
		}
	} else {
		errs = append(errs, "country_code is required")
	}

	if aux.Email != nil {
		if err := json.Unmarshal(aux.Email, &ccfusca.Email); err != nil {
			errs = append(errs, "email must be string")
		}
	} else {
		errs = append(errs, "email is required")
	}

	if aux.AgreeToEula != nil {
		if err := json.Unmarshal(aux.AgreeToEula, &ccfusca.AgreeToEula); err != nil {
			errs = append(errs, "agree_to_eula must be boolean")
		}
	} else {
		errs = append(errs, "agree_to_eula is required")
	}

	if len(errs) > 0 {
		return errs
	}

	return nil

}

func (ccfusca *CarrierConnectFedexUSCA) GetNickname() (string, error) {
	return ccfusca.Nickname, nil
}
