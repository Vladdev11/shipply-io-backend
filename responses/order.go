package responses

import (
	"fmt"
	"time"

	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/util"
)

/* -------------------------------- ListOrders -------------------------------- */
// OrderResponseForListOrders represents the expected response body for the ListOrders endpoint
type ListOrdersResponse struct {
	TotalCount    int                          `json:"total_count"`
	FilteredCount int                          `json:"filtered_count"`
	Data          []OrderResponseForListOrders `json:"data"`
}

// OrderResponseForListOrders represents the expected response body for the ListOrders endpoint
type OrderResponseForListOrders struct {
	ID               int       `json:"id"`
	Priority         int       `json:"priority"`
	OrderNumber      string    `json:"order_number"`
	ReadyToShip      bool      `json:"ready_to_ship"`
	FraudHold        bool      `json:"fraud_hold"`
	AddressHold      bool      `json:"address_hold"`
	PaymentHold      bool      `json:"payment_hold"`
	WarehouseHold    bool      `json:"warehouse_hold"`
	Total            float64   `json:"total"`
	OrderDate        time.Time `json:"order_date"`
	RequiredShipDate time.Time `json:"required_ship_date"`
	HoldUntilDate    time.Time `json:"hold_until_date"`
	Pickable         bool      `json:"pickable"`
	Status           string    `json:"status"`
	ShipToName       string    `json:"ship_to_name"`
	ShipToAddress1   string    `json:"ship_to_address1"`
	ShipToAddress2   string    `json:"ship_to_address2"`
	ShipToCity       string    `json:"ship_to_city"`
	ShipToState      string    `json:"ship_to_state"`
	ShipToPostalCode string    `json:"ship_to_postal_code"`
	ShipToCountry    string    `json:"ship_to_country"`
	Warehouse        string    `json:"warehouse"`
	StoreName        string    `json:"store_name"`
}

// GenerateListOrdersResponse generates the response body for the ListOrders endpoint
func GenerateListOrdersResponse(orders []models.OrdersListOrder, count int, total int) *ListOrdersResponse {

	orderResponses := make([]OrderResponseForListOrders, len(orders))
	for i, order := range orders {
		orderResponses[i] = OrderResponseForListOrders{
			ID:               order.ID,
			Priority:         order.Priority,
			OrderNumber:      order.OrderNumber,
			ReadyToShip:      order.ReadyToShip,
			FraudHold:        util.HasFraudHold(order.Holds),
			AddressHold:      util.HasAddressHold(order.Holds),
			PaymentHold:      util.HasPaymentHold(order.Holds),
			WarehouseHold:    util.HasWarehouseHold(order.Holds),
			Total:            order.Total,
			OrderDate:        order.OrderDate,
			RequiredShipDate: order.RequiredShipDate,
			HoldUntilDate:    order.HoldUntilDate,
			//TODO - handle pickleable
			Pickable:         false,
			Status:           order.Status,
			ShipToName:       order.ShipToName,
			ShipToAddress1:   order.ShipToAddress1,
			ShipToAddress2:   order.ShipToAddress2,
			ShipToCity:       order.ShipToCity,
			ShipToState:      order.ShipToState,
			ShipToPostalCode: order.ShipToPostalCode,
			ShipToCountry:    order.ShipToCountry,
			Warehouse:        order.Warehouse,
			StoreName:        order.StoreName,
		}

	}
	return &ListOrdersResponse{
		TotalCount:    total,
		FilteredCount: count,
		Data:          orderResponses,
	}
}

/* -------------------------------- GetOrder -------------------------------- */

// GetOrderResponse represents the expected response body for the GetOrder endpoint
type GetOrderResponse struct {
	ID                            int       `json:"id"`
	Priority                      int       `json:"priority"`
	OrderNumber                   string    `json:"order_number"`
	GiftNote                      string    `json:"gift_note"`
	PackingNote                   string    `json:"packing_note"`
	ReadyToShip                   bool      `json:"ready_to_ship"`
	FraudHold                     bool      `json:"fraud_hold"`
	AddressHold                   bool      `json:"address_hold"`
	PaymentHold                   bool      `json:"payment_hold"`
	WarehouseHold                 bool      `json:"warehouse_hold"`
	Pickable                      bool      `json:"pickable"`
	Subtotal                      float64   `json:"subtotal"`
	Tax                           float64   `json:"tax"`
	Shipping                      float64   `json:"shipping"`
	Total                         float64   `json:"total"`
	OrderDate                     time.Time `json:"order_date"`
	RequiredShipDate              time.Time `json:"required_ship_date"`
	HoldUntilDate                 time.Time `json:"hold_until_date"`
	FulfillmentStatus             string    `json:"fulfillment_status"`
	AutoPrintReturnLabel          bool      `json:"auto_print_return_label"`
	CustomerEmail                 string    `json:"customer_email"`
	CustomerPhone                 string    `json:"customer_phone"`
	SaturdayDelivery              bool      `json:"saturday_delivery"`
	IgnoreAddressValidationErrors bool      `json:"ignore_address_validation_errors"`
	AllocationPriority            int       `json:"allocation_priority"`
	AllowPartial                  bool      `json:"allow_partial"`
	GiftInvoice                   bool      `json:"gift_invoice"`
	RequireSignature              bool      `json:"require_signature"`
	AdultSignatureRequired        bool      `json:"adult_signature_required"`
	Alcohol                       bool      `json:"alchohol"`
	Insurance                     bool      `json:"insurance"`
	InsuranceValue                float64   `json:"insurance_value"`
	HasDryIce                     bool      `json:"has_dry_ice"`
	Tags                          []string  `json:"tags"`
	CreatedAt                     time.Time `json:"created_at"`
	UpdatedAt                     time.Time `json:"updated_at"`

	Store          StoreResponseForGetOrder          `json:"store"`
	Status         StatusResponseForGetOrder         `json:"status"`
	Warehouse      WarehouseResponseForGetOrder      `json:"warehouse"`
	BillToAddress  Address                           `json:"bill_to_address"`
	ShipToAddress  Address                           `json:"ship_to_address"`
	OrderItems     []OrderItemResponseForGetOrder    `json:"order_items"`
	ShippingMethod ShippingMethodResponseForGetOrder `json:"shipping_method"`
	OrderHistory   []OrderHistoryResponseForGetOrder `json:"order_history"`
}

// OrderHistoryResponseForGetOrder represents the expected response body for an individual order history in the GetOrder endpoint
type OrderHistoryResponseForGetOrder struct {
	CreatedAt     time.Time     `json:"created_at"`
	Note          string        `json:"note"`
	ChangedByUser ChangedByUser `json:"changed_by_user"`
}

// StoreResponseForGetOrder represents the expected response body for an individual store in the GetOrder endpoint
type StoreResponseForGetOrder struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// StatusResponseForGetOrder represents the expected response body for an individual status in the GetOrder endpoint
type StatusResponseForGetOrder struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// WarehouseResponseForGetOrder represents the expected response body for an individual warehouse in the GetOrder endpoint
type WarehouseResponseForGetOrder struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// OrderItemResponseForGetOrder represents the expected response body for an individual order item in the GetOrder endpoint
type OrderItemResponseForGetOrder struct {
	ID                 int     `json:"id"`
	Sku                string  `json:"sku"`
	Name               string  `json:"name"`
	PendingFulfillment int     `json:"pending_fulfillment"`
	Ordered            int     `json:"ordered"`
	Shipped            int     `json:"shipped"`
	Allocated          int     `json:"allocated"`
	Backordered        int     `json:"backordered"`
	UnitPrice          float64 `json:"unit_price"`
	TotalPrice         float64 `json:"total_price"`
	Status             string  `json:"status"`
}

// ShippingMethodResponseForGetOrder represents the expected response body for an individual shipping method in the GetOrder endpoint
type ShippingMethodResponseForGetOrder struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GenerateGetOrderResponse generates a GetOrderResponse struct
func GenerateGetOrderResponse(order models.Order) *GetOrderResponse {

	orderHistories := make([]OrderHistoryResponseForGetOrder, len(order.History))
	for i, orderHistory := range order.History {
		orderHistories[i] = OrderHistoryResponseForGetOrder{
			CreatedAt: orderHistory.CreatedAt,
			Note:      orderHistory.Note,
			ChangedByUser: ChangedByUser{
				ID:        orderHistory.CreatedByUser.ID,
				FirstName: orderHistory.CreatedByUser.FirstName,
				LastName:  orderHistory.CreatedByUser.LastName,
				ImageURL:  fmt.Sprintf("%s/%s", util.ConfigCDNHost, orderHistory.CreatedByUser.AvatarFileName),
			},
		}
	}

	orderItems := make([]OrderItemResponseForGetOrder, len(order.OrderItems))
	for i, orderItem := range order.OrderItems {
		orderItems[i] = OrderItemResponseForGetOrder{
			ID:                 orderItem.ID,
			Sku:                orderItem.Sku,
			Name:               orderItem.Name,
			PendingFulfillment: orderItem.Quantity - orderItem.QuantityShipped,
			Ordered:            orderItem.Quantity,
			Shipped:            orderItem.QuantityShipped,
			Allocated:          orderItem.Allocated,
			Backordered:        orderItem.Backordered,
			UnitPrice:          orderItem.ItemPrice,
			TotalPrice:         orderItem.ItemPrice * float64(orderItem.Quantity),
			Status:             orderItem.Status,
		}
	}

	orderTags := make([]string, len(order.Tags))
	for i, orderTag := range order.Tags {
		orderTags[i] = orderTag.Tag
	}

	return &GetOrderResponse{
		ID:            order.ID,
		Priority:      order.Priority,
		OrderNumber:   order.OrderNumber,
		GiftNote:      order.GiftNote,
		PackingNote:   order.PackingNote,
		ReadyToShip:   order.ReadyToShip,
		FraudHold:     util.HasFraudHold(order.Holds),
		AddressHold:   util.HasAddressHold(order.Holds),
		PaymentHold:   util.HasPaymentHold(order.Holds),
		WarehouseHold: util.HasWarehouseHold(order.Holds),
		//TODO - determine if pickable or not
		Pickable:                      false,
		Subtotal:                      order.Subtotal,
		Tax:                           order.Tax,
		Shipping:                      order.Shipping,
		Total:                         order.Total,
		OrderDate:                     order.OrderDate,
		RequiredShipDate:              order.RequiredShipDate,
		HoldUntilDate:                 order.HoldUntilDate,
		FulfillmentStatus:             order.FulfillmentStatus,
		AutoPrintReturnLabel:          order.AutoPrintReturnLabel,
		CustomerEmail:                 order.CustomerEmail,
		CustomerPhone:                 order.CustomerPhone,
		SaturdayDelivery:              order.SaturdayDelivery,
		IgnoreAddressValidationErrors: order.IgnoreAddressValidationErrors,
		AllocationPriority:            order.AllocationPriority,
		AllowPartial:                  order.AllowPartial,
		GiftInvoice:                   order.GiftInvoice,
		RequireSignature:              order.RequireSignature,
		AdultSignatureRequired:        order.AdultSignatureRequired,
		Alcohol:                       order.Alcohol,
		Insurance:                     order.Insurance,
		InsuranceValue:                order.InsuranceValue,
		HasDryIce:                     order.HasDryIce,
		Tags:                          orderTags,
		CreatedAt:                     order.CreatedAt,
		UpdatedAt:                     order.UpdatedAt,
		Store: StoreResponseForGetOrder{
			ID:   order.Store.ID,
			Name: order.Store.Name,
		},
		Status: StatusResponseForGetOrder{
			ID:   order.Status.ID,
			Name: order.Status.Name,
		},
		Warehouse: WarehouseResponseForGetOrder{
			ID:   order.Warehouse.ID,
			Name: order.Warehouse.Name,
		},
		BillToAddress: Address{
			ID:         order.BillToAddress.ID,
			Street1:    order.BillToAddress.Street1,
			Street2:    order.BillToAddress.Street2,
			Street3:    order.BillToAddress.Street3,
			City:       order.BillToAddress.City,
			State:      order.BillToAddress.State,
			PostalCode: order.BillToAddress.PostalCode,
			Country:    order.BillToAddress.Country,
			Company:    order.BillToAddress.Company,
			Phone:      order.BillToAddress.Phone,
		},
		ShipToAddress: Address{
			ID:         order.ShipToAddress.ID,
			Street1:    order.ShipToAddress.Street1,
			Street2:    order.ShipToAddress.Street2,
			Street3:    order.ShipToAddress.Street3,
			City:       order.ShipToAddress.City,
			State:      order.ShipToAddress.State,
			PostalCode: order.ShipToAddress.PostalCode,
			Country:    order.ShipToAddress.Country,
			Company:    order.ShipToAddress.Company,
			Phone:      order.ShipToAddress.Phone,
		},
		OrderItems: orderItems,
		ShippingMethod: ShippingMethodResponseForGetOrder{
			ID:   order.ShippingMethod.ID,
			Name: order.ShippingMethod.Name,
		},
		OrderHistory: orderHistories,
	}
}
