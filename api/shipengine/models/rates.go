package ShipengineModels

import (
	"encoding/json"
	"errors"

	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/util"
)

type RateShopRequest struct {
	ShipmentID  string      `json:"shipment_id,omitempty"`
	RateOptions RateOptions `json:"rate_options,omitempty"`
	Shipment    Shipment    `json:"shipment,omitempty"`
}

type RateOptions struct {
	CarrierIds         []string `json:"carrier_ids,omitempty"`
	ServiceCodes       []string `json:"service_codes,omitempty"`
	PackageTypes       []string `json:"package_types,omitempty"`
	CalculateTaxAmount bool     `json:"calculate_tax_amount,omitempty"`
	PreferredCurrency  string   `json:"preferred_currency,omitempty"`
	IsReturn           bool     `json:"is_return,omitempty"`
}

type RateShopResponse struct {
	RateResponse       RateResponse `json:"rate_response,omitempty"`
	ShipmentID         string       `json:"shipment_id,omitempty"`
	CarrierID          string       `json:"carrier_id,omitempty"`
	ServiceCode        string       `json:"service_code,omitempty"`
	ExternalOrderID    string       `json:"external_order_id,omitempty"`
	ExternalShipmentID string       `json:"external_shipment_id,omitempty"`
	ShipmentNumber     string       `json:"shipment_number,omitempty"`
	ShipDate           string       `json:"ship_date,omitempty"`
	CreatedAt          string       `json:"created_at,omitempty"`
	ModifiedAt         string       `json:"modified_at,omitempty"`
	ShipmentStatus     string       `json:"shipment_status,omitempty"`
	ShipTo             Address      `json:"ship_to,omitempty"`
	ShipFrom           Address      `json:"ship_from,omitempty"`
	WarehouseID        string       `json:"warehouse_id,omitempty"`
	ReturnTo           Address      `json:"return_to,omitempty"`
	IsReturn           bool         `json:"is_return,omitempty"`
	Confirmation       string       `json:"confirmation,omitempty"`
	Customs            CustomsItem  `json:"customs,omitempty"`
	// TaxIdentifiers	 []TaxIdentifier     `json:"tax_identifiers"`
	AdvancedOptions   ShipmentAdvancedOptions `json:"advanced_options,omitempty"`
	OriginType        string                  `json:"origin_type,omitempty"`
	InsuranceProvider string                  `json:"insurance_provider,omitempty"`
	OrderSourceCode   string                  `json:"order_source_code,omitempty"`
	Packages          []Package               `json:"packages,omitempty"`
	TotalWeight       Weight                  `json:"total_weight,omitempty"`
}

type RateResponse struct {
	Rates []Rate `json:"rates"`
	// InvalidRates []InvalidRate `json:"invalid_rates"`
	RateRequestID string            `json:"rate_request_id,omitempty"`
	ShipmentID    string            `json:"shipment_id,omitempty"`
	CreatedAt     string            `json:"created_at,omitempty"`
	Status        string            `json:"status,omitempty"`
	Errors        []ShipengineError `json:"errors,omitempty"`
}

type Rate struct {
	RateID                string   `json:"rate_id,omitempty"`
	RateType              string   `json:"rate_type,omitempty"`
	CarrierID             string   `json:"carrier_id,omitempty"`
	ShippingAmount        Money    `json:"shipping_amount,omitempty"`
	InsuranceAmount       Money    `json:"insurance_amount,omitempty"`
	ConfirmationAmount    Money    `json:"confirmation_amount,omitempty"`
	OtherAmount           Money    `json:"other_amount,omitempty"`
	TaxAmount             Money    `json:"tax_amount,omitempty"`
	Zone                  int      `json:"zone,omitempty"`
	PackageType           string   `json:"package_type,omitempty"`
	DeliveryDays          int      `json:"delivery_days,omitempty"`
	GuaranteedService     bool     `json:"guaranteed_service,omitempty"`
	EstimatedDeliveryDate string   `json:"estimated_delivery_date,omitempty"`
	CarrierDeliveryDays   string   `json:"carrier_delivery_days,omitempty"`
	ShipDate              string   `json:"ship_date,omitempty"`
	NegotiatedRate        bool     `json:"negotiated_rate,omitempty"`
	ServiceType           string   `json:"service_type,omitempty"`
	ServiceCode           string   `json:"service_code,omitempty"`
	Trackable             bool     `json:"trackable,omitempty"`
	CarrierCode           string   `json:"carrier_code,omitempty"`
	CarrierNickname       string   `json:"carrier_nickname,omitempty"`
	CarrierFriendlyName   string   `json:"carrier_friendly_name,omitempty"`
	ValidationStatus      string   `json:"validation_status,omitempty"`
	WarningMessages       []string `json:"warning_messages,omitempty"`
	ErrorMessages         []string `json:"error_messages,omitempty"`
}

func ConstructRateShopRequest(pickSessionOrder models.PickSessionOrder, boxID int, weight float64, shopAllServiceCodes bool) (*RateShopRequest, error) {
	box, err := models.GetBoxByID(boxID)
	if err != nil {
		return nil, errors.New("failed to get box")
	}

	order, err := pickSessionOrder.GetOrder()
	if err != nil {
		return nil, errors.New("failed to get order")
	}

	err = order.GetShippingMethod()
	if err != nil {
		return nil, errors.New("failed to get shipping method")
	}

	if !order.ShippingMethod.Mapped {
		return nil, errors.New("shipping method is not mapped")
	}

	shippingMethodCarriers := []models.ShippingMethodCarrier{}
	err = json.Unmarshal(order.ShippingMethod.Carriers, &shippingMethodCarriers)
	if err != nil {
		return nil, errors.New("failed to unmarshal shipping method carriers")
	}

	// Append carrier codes and enabled carrier service codes to arrays
	enabledCarrierServiceCodes := []string{}
	carrierCodes := []string{}
	for _, carrier := range shippingMethodCarriers {
		carrierConnection, err := models.GetCarrierConnectionByID(carrier.CarrierConnectionID)
		if err != nil {
			// TODO error log
			continue
		}

		if !carrierConnection.Active {
			continue
		}

		// skip cheapest (cheapest is only for selecting a rate automatically after shopping rates)
		if carrierConnection.ID == util.CheapestCarrierConnectionID {
			continue
		}

		carrierCodes = append(carrierCodes, carrierConnection.ShipengineCarrierID)

		for _, service := range carrier.Services {
			if service.Enabled || shopAllServiceCodes {
				enabledCarrierServiceCodes = append(enabledCarrierServiceCodes, service.ServiceCode)
			}
		}

	}

	if len(carrierCodes) == 0 {
		return nil, errors.New("no active carrier connections")
	}

	// TODO if box is carrier specific, only shop rates for that carrier

	err = order.GetShipToAddress()
	if err != nil {
		return nil, errors.New("failed to get ship to address")
	}

	err = order.GetWarehouse()
	if err != nil {
		return nil, errors.New("failed to get warehouse")
	}

	// REVIEW (for brennan)  -- Do we need to add box weight to the total weight? I am assuming the weight from the request will include the box weight, regardless of if they're creating a new box or not

	// TODO customs and other fields that I missed
	shipengineRateShopRequest := RateShopRequest{
		Shipment: Shipment{
			ValidateAddress: ValidateAddressNone,
			ShipTo:          ConvertAddressToShipengineAddress(&order.ShipToAddress),
			WarehouseID:     order.Warehouse.ShipengineWarehouseID,
			Packages: &[]Package{
				{
					Weight: &Weight{
						Value: weight,
						// TODO get unit from request
						Unit: "pound",
					},
					PackageID:   box.ShipenginePackageID,
					PackageCode: box.ShipenginePackageCode,
				},
			},
		},
		RateOptions: RateOptions{
			CarrierIds:        carrierCodes,
			ServiceCodes:      enabledCarrierServiceCodes,
			PreferredCurrency: order.Currency,
		},
		// TODO Account for potentially multiple packages

	}

	if order.Alcohol || order.SaturdayDelivery || order.HasDryIce {
		shipengineRateShopRequest.Shipment.AdvancedOptions = &ShipmentAdvancedOptions{}

		if order.Alcohol {
			shipengineRateShopRequest.Shipment.AdvancedOptions.ContainsAlcohol = true
		}

		if order.SaturdayDelivery {
			shipengineRateShopRequest.Shipment.AdvancedOptions.SaturdayDelivery = true
		}

		if order.HasDryIce {
			shipengineRateShopRequest.Shipment.AdvancedOptions.DryIce = true
			shipengineRateShopRequest.Shipment.AdvancedOptions.DryIceWeight = &Weight{
				Value: order.DryIceWeightInLbs,
				Unit:  "pound",
			}
		}
	}

	return &shipengineRateShopRequest, nil
}
