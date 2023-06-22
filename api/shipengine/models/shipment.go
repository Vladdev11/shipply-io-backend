package ShipengineModels

type CreateShipmentRequest struct {
	Shipments []Shipment `json:"shipments"`
}

type CreateShipmentResponse struct {
	HasErrors bool       `json:"has_errors"`
	Shipments []Shipment `json:"shipments"`
}

type Shipment struct {
	ValidateAddress ValidateAddress `json:"validate_address,omitempty"`  // Level of address validation to be performed.
	CarrierID       string          `json:"carrier_id,omitempty"`        // ID of the carrier for shipping.
	ServiceCode     string          `json:"service_code,omitempty"`      // Carrier service code for shipping.
	ExternalOrderID string          `json:"external_order_id,omitempty"` // Optional identifier for associated order.
	// TaxIdentifiers     []CustomsTaxIdentifier `json:"tax_identifiers"`      // Tax identifiers for international shipments.
	ExternalShipmentID string   `json:"external_shipment_id,omitempty"` // Optional external identifier for shipment.
	ShipmentNumber     string   `json:"shipment_number,omitempty"`      // Optional shipment number to identify shipment.
	ShipDate           string   `json:"ship_date,omitempty"`            // Expected ship date.
	ShipTo             Address  `json:"ship_to,omitempty"`              // Destination address for package.
	ShipFrom           *Address `json:"ship_from,omitempty"`            // Origin address for package.
	WarehouseID        string   `json:"warehouse_id"`                   // Optional ID of shipping warehouse.
	ReturnTo           *Address `json:"return_to,omitempty"`            // Return address for non-delivery.
	Confirmation       string   `json:"confirmation,omitempty"`         // Type of delivery confirmation required.
	// Customs            Customs                `json:"customs"`              // Customs info for international shipments.
	AdvancedOptions   *ShipmentAdvancedOptions `json:"advanced_options,omitempty"`   // Additional options for shipment.
	OriginType        string                   `json:"origin_type,omitempty"`        // Type of origin (residential, commercial).
	InsuranceProvider string                   `json:"insurance_provider,omitempty"` // Provider for insuring shipment.
	OrderSourceCode   string                   `json:"order_source_code,omitempty"`  // Source of associated order (e-commerce platform).
	Packages          *[]Package               `json:"packages,omitempty"`           // Array of packages in the shipment.
}

type ShipmentAdvancedOptions struct {
	BillToAccount     string  `json:"bill_to_account,omitempty"`
	BillToCountryCode string  `json:"bill_to_country_code,omitempty"`
	BillToParty       string  `json:"bill_to_party,omitempty"`
	BillToPostalCode  string  `json:"bill_to_postal_code,omitempty"`
	ContainsAlcohol   bool    `json:"contains_alcohol,omitempty"`
	DeliveryDutyPaid  bool    `json:"delivery_duty_paid,omitempty"`
	DryIce            bool    `json:"dry_ice,omitempty"`
	DryIceWeight      *Weight `json:"dry_ice_weight,omitempty"`
	NonMachinable     bool    `json:"non_machinable,omitempty"`
	SaturdayDelivery  bool    `json:"saturday_delivery,omitempty"`
	FedexFreight      *struct {
		ShipperLoadAndCount string `json:"shipper_load_and_count,omitempty"`
		BookingConfirmation string `json:"booking_confirmation,omitempty"`
	} `json:"fedex_freight,omitempty"`
	UseUPSGroundFreightPricing bool    `json:"use_ups_ground_freight_pricing,omitempty"`
	FreightClass               float64 `json:"freight_class,omitempty"`
	OriginType                 string  `json:"origin_type,omitempty"`
	ShipperRelease             bool    `json:"shipper_release,omitempty"`
	CollectOnDelivery          *struct {
		PaymentType   string `json:"payment_type,omitempty"`
		PaymentAmount Money  `json:"payment_amount,omitempty"`
	} `json:"collect_on_delivery,omitempty"`
	ThirdPartyConsignee bool `json:"third_party_consignee,omitempty"`
}

type Customs struct {
	Contents     string        `json:"contents"`
	CustomsItems []CustomsItem `json:"customs_items"`
	NonDelivery  string        `json:"non_delivery"`
}

type CustomsItem struct {
	CustomsItemID        string  `json:"customs_item_id"` //Shipengine ID (we don't submit this in the request)
	Description          string  `json:"description"`
	Quantity             int     `json:"quantity"`
	Value                float64 `json:"value"`
	HarmonziedTariffCode string  `json:"harmonized_tariff_code"`
	CountryOfOrigin      string  `json:"country_of_origin"`
	Sku                  string  `json:"sku"` //must be between 1 and 20 characters
	SkuDescription       string  `json:"sku_description"`
}

type ValidateAddress string

const (
	ValidateAddressNone             ValidateAddress = "no_validation"
	ValidateAddressValidateOnly     ValidateAddress = "validate_only"
	ValidateAddressValidateAndClean ValidateAddress = "validate_and_clean"
)
