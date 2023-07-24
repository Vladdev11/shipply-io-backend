package models

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"
)

var (
	//ErrStoreisNotAShopifyStore is returned when a store is not a shopify store
	ErrStoreisNotAShopifyStore = errors.New("store is not a shopify store")
	//ErrInvalidShopifyCredentialsJSON is returned when the shopify credentials json is invalid
	ErrInvalidShopifyCredentialsJSON = errors.New("invalid shopify credentials json")
)

type Store struct {
	ID             int
	ClientID       int
	MarketplaceID  int
	Name           string
	APICredentials json.RawMessage `gorm:"type:jsonb"`
	Settings       json.RawMessage `gorm:"type:jsonb"`
	Active         bool

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	Marketplace Marketplace `gorm:"-"`
	Client      Client
}

type StoreReturnJSON struct {
	ID          int                   `json:"id"`
	ClientID    int                   `json:"client_id"`
	Name        string                `json:"name"`
	Marketplace MarketplaceReturnJSON `json:"marketplace"`
	Settings    json.RawMessage       `json:"settings"`
	Active      bool                  `json:"active"`
}

type StoreListRequest struct {
	ClientID       *int `json:"client_id"`
	OrganizationID int  `json:"organization_id"`
}

type StoreUpdateRequest struct {
	Name string `json:"name"`
	// TODO add settings
}

func (s *Store) GetMarketplace() error {

	if s.MarketplaceID == 0 {
		return errors.New("no marketplace id")
	}

	marketplace, ok := Marketplaces[s.MarketplaceID]
	if !ok {
		return errors.New("invalid marketplace id")
	}

	s.Marketplace = marketplace

	return nil
}

func (s *Store) ConvertToReturnJSON(ctx context.Context) *StoreReturnJSON {

	if s == nil {
		return nil
	}

	err := s.GetMarketplace()
	if err != nil {
		return nil
	}

	return &StoreReturnJSON{
		ID:          s.ID,
		ClientID:    s.ClientID,
		Name:        s.Name,
		Settings:    s.Settings,
		Active:      s.Active,
		Marketplace: *s.Marketplace.ConvertToReturnJSON(),
	}
}

func (slr *StoreListRequest) ParseAndValidateRequest(r *http.Request) []string {
	errors := []string{}

	clientID, err := util.GetIntQueryParam(r, "client_id")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "client_id must be an integer")
		} else {
			slr.ClientID = &clientID
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

func (s *StoreUpdateRequest) ParseAndValidateRequest(r *http.Request) []string {

	var errs []string

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return []string{"invalid JSON"}
	}

	aux := &struct {
		Name json.RawMessage `json:"name"`
	}{}

	if err := json.Unmarshal(body, aux); err != nil {
		return []string{"invalid JSON"}
	}

	if aux.Name != nil {
		if err := json.Unmarshal(aux.Name, &s.Name); err != nil {
			errs = append(errs, "name must be a string")
		} else if len(s.Name) > 255 {
			errs = append(errs, "name must be less than 255 characters")
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil

}

func (s *Store) UpdateWithRequest(ctx context.Context, request *StoreUpdateRequest) error {

	if request.Name != "" {
		s.Name = request.Name
	}

	if err := util.DBFromContext(ctx).Save(s).Error; err != nil {
		return err
	}

	return nil
}

type ShopifyStoreSettings struct {

	//BY DEFAULT WE ASSUME WE DOWNLOAD / UPDATE FROM SHOPIFY
	//ONLY IF SELECTED DO WE OVERRIDE THE SETTINGS TO SHOPIFY

	//GENERAL SETTINGS

	//PRODUCT SETTINGS
	OnlyImportActiveProducts       bool `json:"only_import_active_products"`       //When true, we will only import published products from shopify
	OverwriteProductBarcode        bool `json:"overwrite_product_barcode"`         //When true, we will overwrite the barcode in shopify and manage it going forward
	OverwriteProductWeight         bool `json:"overwrite_product_weight"`          //When true, we will overwrite the weight in shopify and manage it going forward
	OverwriteProductPrice          bool `json:"overwrite_product_price"`           //When true, we will overwrite the price in shopify and manage it going forward
	OverwriteProductHarmonizedCode bool `json:"overwrite_product_harmonized_code"` //When true, we will overwrite the harmonized code in shopify and manage it going forward
	OverwriteCountryOfOrigin       bool `json:"overwrite_country_of_origin"`       //When true, we will overwrite the country of origin in shopify and manage it going forward

	//INVENTORY SETTINGS
	OverwriteShopifyInventory bool `json:"overwrite_shopify_inventory"` //When true, we will overwrite inventory in shopify and manage it going forward

	//FRAUD SETTINGS
	ApplyFraudHoldLevel ShopifyFraudLevel `json:"apply_fraud_hold"` //When true, we will apply a fraud hold to orders that are flagged as fraud based on your fraud settings

	//ORDER SETTINGS
	ShopifyOrderFulfilledEmail       bool                             `json:"shopify_order_fulfilled_email"`        //When true, we will trigger the shopify order fulfilled email
	OverwriteOrderAddress            bool                             `json:"overwrite_order_address"`              //When true, we will overwrite the address in shopify and manage it going forward
	OnlyImportPaidOrders             bool                             `json:"only_import_paid_orders"`              //When true, we will only import paid orders from shopify
	CustomerNotesFieldMapping        ShopifyCustomerNotesFieldMapping `json:"customer_notes_field_mapping"`         //This is the field in shopify that we will map to the customer notes field in shipply
	IgnorePaymentStatus              bool                             `json:"ignore_payment_status"`                //When true, we will ignore the payment status of the order and import it anyway
	UpdateShippingAddressFromShopify bool                             `json:"update_shipping_address_from_shopify"` //When true, we will update the shipping address in shipply from shopify
	UpdateShopifyAddressFromShipply  bool                             `json:"update_shopify_address_from_shipply"`  //When true, we will update the address in shopify from shipply

}

var DefaultShopifyStoreSettings = ShopifyStoreSettings{
	OnlyImportActiveProducts:       true,
	OverwriteProductBarcode:        true,
	OverwriteProductWeight:         true,
	OverwriteProductPrice:          true,
	OverwriteProductHarmonizedCode: true,
	OverwriteCountryOfOrigin:       true,
	OverwriteShopifyInventory:      true,
	ApplyFraudHoldLevel:            ShopifyFraudLevelLow,
	ShopifyOrderFulfilledEmail:     true,
	OverwriteOrderAddress:          true,
	OnlyImportPaidOrders:           true,
	CustomerNotesFieldMapping:      CustomerNotesFieldMappingNonePacking,
	IgnorePaymentStatus:            false,
}

type ShopifyFraudLevel string

const (
	ShopifyFraudLevelLow    ShopifyFraudLevel = "low"
	ShopifyFraudLevelMedium ShopifyFraudLevel = "medium"
	ShopifyFraudLevelHigh   ShopifyFraudLevel = "high"
)

type ShopifyCustomerNotesFieldMapping string

const (
	CustomerNotesFieldMappingNonePacking ShopifyCustomerNotesFieldMapping = "packing_note"
	CustomerNotesFieldMappingNoteGift    ShopifyCustomerNotesFieldMapping = "gift_note"
)

type ShopifyCredentials struct {
	ShopName            string `json:"shop_name"`
	AccessToken         string `json:"access_token"`
	ShopifyInstallNonce string `json:"shopify_install_nonce"`
}

func (s *Store) Create(ctx context.Context) error {
	var err error
	//init default settings
	if s.Settings == nil {
		switch s.MarketplaceID {
		case util.ShopifyMarketplaceID:
			s.Settings, err = json.Marshal(DefaultShopifyStoreSettings)
			if err != nil {
				return err
			}
		}
	}
	return util.DBFromContext(ctx).Create(s).Error
}

func (s *Store) Delete(ctx context.Context) error {
	err := util.DBFromContext(ctx).Delete(s).Error
	if err != nil {
		return ErrDeleteFailed{Object: "store", Err: err}
	}
	return nil
}

func (s *Store) Activate(ctx context.Context) error {
	s.Active = true
	err := util.DBFromContext(ctx).Save(s).Error
	if err != nil {
		return ErrUpdateFailed{Object: "store", Err: err}
	}
	return nil
}

func (s *Store) Deactivate(ctx context.Context) error {
	s.Active = false
	err := util.DBFromContext(ctx).Save(s).Error
	if err != nil {
		return ErrUpdateFailed{Object: "store", Err: err}
	}
	return nil
}

func (s *Store) UpdateAPICredentials(ctx context.Context, newCredentials interface{}) error {
	credentialsJSON, err := json.Marshal(newCredentials)
	if err != nil {
		return err
	}
	s.APICredentials = json.RawMessage(credentialsJSON)
	return util.DBFromContext(ctx).Save(s).Error
}

func (s *Store) UpdateSettings(ctx context.Context, newSettings interface{}) error {
	settingsJSON, err := json.Marshal(newSettings)
	if err != nil {
		return err
	}
	s.Settings = json.RawMessage(settingsJSON)
	return util.DBFromContext(ctx).Save(s).Error
}

func (s *Store) GetShopifySettings(ctx context.Context) (*ShopifyStoreSettings, error) {

	if s.MarketplaceID != util.ShopifyMarketplaceID {
		return nil, errors.New("Store is not a Shopify store")
	}

	var settings ShopifyStoreSettings

	if err := json.Unmarshal(s.Settings, &settings); err != nil {
		return nil, errors.New("invalid shopify settings json")
	}

	return &settings, nil
}

func (s *Store) GetShopifyCredentials(ctx context.Context) (*ShopifyCredentials, error) {

	if s.MarketplaceID != util.ShopifyMarketplaceID {
		return nil, ErrStoreisNotAShopifyStore
	}

	var credentials ShopifyCredentials

	aux := &struct {
		ShopName            json.RawMessage `json:"shop_name"`
		AccessToken         json.RawMessage `json:"access_token"`
		ShopifyInstallNonce json.RawMessage `json:"shopify_install_nonce"`
	}{}

	if err := json.Unmarshal(s.APICredentials, aux); err != nil {
		return nil, ErrInvalidShopifyCredentialsJSON
	}

	if aux.ShopName != nil {
		if err := json.Unmarshal(aux.ShopName, &credentials.ShopName); err != nil {
			return nil, errors.New("ShopName must be a string")
		} else if !strings.Contains(credentials.ShopName, ".myshopify.com") {
			return nil, errors.New("ShopName in an invalid shopify domain")
		}
	}

	if aux.AccessToken != nil {
		if err := json.Unmarshal(aux.AccessToken, &credentials.AccessToken); err != nil {
			return nil, errors.New("AccessToken must be a string")
		}
	}

	if aux.ShopifyInstallNonce != nil {
		if err := json.Unmarshal(aux.ShopifyInstallNonce, &credentials.ShopifyInstallNonce); err != nil {
			return nil, errors.New("ShopifyInstallNonce must be a string")
		}
	}

	return &credentials, nil
}

func (s *Store) DeleteAccessToken(ctx context.Context) error {

	credentials, err := s.GetShopifyCredentials(ctx)
	if err != nil {
		return err
	}

	credentials.AccessToken = ""
	return s.UpdateAPICredentials(ctx, credentials)

}

func (s *Store) DeleteShopifyLocations(ctx context.Context) error {
	return util.DBFromContext(ctx).Where("store_id = ?", s.ID).Delete(&ShopifyLocation{}).Error
}

func GetShopifyStoreByShopName(ctx context.Context, shopName string) (*Store, error) {
	var store Store
	err := util.DBFromContext(ctx).Where("api_credentials ->> 'shop_name' = ?", shopName).First(&store).Error
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (s *Store) ShopifyShopName(ctx context.Context) string {
	var shopName string
	util.DBFromContext(ctx).Raw("SELECT api_credentials ->> 'shop_name' FROM stores WHERE id = ?", s.ID).Scan(&shopName)
	return shopName
}

func (s *Store) ShopifyAccessToken(ctx context.Context) string {
	var accessToken string
	util.DBFromContext(ctx).Raw("SELECT api_credentials ->> 'access_token' FROM stores WHERE id = ?", s.ID).Scan(&accessToken)
	return accessToken
}

func GetStoreByID(ctx context.Context, id int) (*Store, error) {
	var store Store
	err := util.DBFromContext(ctx).First(&store, id).Error
	if err != nil {
		return nil, ErrQueryFailed{Object: "store", Err: err}
	}
	return &store, nil
}
