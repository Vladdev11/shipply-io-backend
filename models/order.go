package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"
)

type Order struct {
	ID                            int
	APIID                         string
	Priority                      int
	OrderNumber                   string
	StoreID                       int
	Status                        string
	WarehouseID                   *int `gorm:"foreignKey:WarehouseID;type:integer;null"`
	GiftNote                      string
	PackingNote                   string
	ReadyToShip                   bool
	BillToAddressID               *int `gorm:"foreignKey:BillToAddressID;type:integer;null"`
	ShipToAddressID               *int `gorm:"foreignKey:ShipToAddressID;type:integer;null"`
	ShippingMethodID              *int `gorm:"foreignKey:ShippingMethodID;type:integer;null"`
	Holds                         uint64
	Subtotal                      float64
	Tax                           float64
	Shipping                      float64
	Discount                      float64
	DiscountCodes                 json.RawMessage `gorm:"type:jsonb"`
	Tip                           float64
	Total                         float64
	OrderDate                     time.Time
	RequiredShipDate              time.Time
	HoldUntilDate                 time.Time
	FulfillmentStatus             string
	MarketplaceFinacialStatus     string
	MarketplaceFulfillmentStatus  string
	BoxID                         *int `gorm:"foreignKey:BoxID;type:integer;null"`
	AutoPrintReturnLabel          bool
	CustomerEmail                 string
	CustomerPhone                 string
	SaturdayDelivery              bool
	IgnoreAddressValidationErrors bool
	SkipAddressValidation         bool
	AllocationPriority            int
	AllowPartial                  bool
	GiftInvoice                   bool
	RequireSignature              bool
	AdultSignatureRequired        bool
	Alcohol                       bool
	Insurance                     bool
	InsuranceValue                float64
	Currency                      string
	HasDryIce                     bool
	DryIceWeightInLbs             float64
	AllowSplit                    bool
	FTRExemption                  bool

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	Box            Box
	OrderItems     []OrderItem
	Warehouse      Warehouse
	BillToAddress  Address
	ShipToAddress  Address
	ShippingMethod ShippingMethod
	Store          Store `gorm:"foreignKey:StoreID;type:integer;null"`
}

type InventoryCount struct {
	ProductID      uint
	LocationID     uint
	AvailableCount int
}

type OrderReturnJSON struct {
	ID                            int             `json:"id"`
	APIID                         string          `json:"api_id,omitempty"`
	Priority                      int             `json:"priority"`
	OrderNumber                   string          `json:"order_number"`
	StoreID                       int             `json:"store_id,omitempty"`
	WarehouseID                   *int            `json:"warehouse_id,omitempty"`
	Status                        string          `json:"status"`
	GiftNote                      string          `json:"gift_note"`
	PackingNote                   string          `json:"packing_note"`
	ReadyToShip                   bool            `json:"ready_to_ship"`
	BillToAddressID               *int            `json:"bill_to_address_id,omitempty"`
	ShipToAddressID               *int            `json:"ship_to_address_id,omitempty"`
	ShippingMethodID              *int            `json:"shipping_method_id,omitempty"`
	Holds                         uint64          `json:"holds"`
	Subtotal                      float64         `json:"subtotal"`
	Tax                           float64         `json:"tax"`
	Shipping                      float64         `json:"shipping"`
	Discount                      float64         `json:"discount"`
	DiscountCodes                 json.RawMessage `json:"discount_codes"`
	Tip                           float64         `json:"tip"`
	Total                         float64         `json:"total"`
	OrderDate                     time.Time       `json:"order_date"`
	RequiredShipDate              time.Time       `json:"required_ship_date"`
	HoldUntilDate                 time.Time       `json:"hold_until_date"`
	FulfillmentStatus             string          `json:"fulfillment_status"`
	MarketplaceFinacialStatus     string          `json:"marketplace_finacial_status"`
	MarketplaceFulfillmentStatus  string          `json:"marketplace_fulfillment_status"`
	BoxID                         *int            `json:"box_id"`
	AutoPrintReturnLabel          bool            `json:"auto_print_return_label"`
	CustomerEmail                 string          `json:"customer_email"`
	CustomerPhone                 string          `json:"customer_phone"`
	SaturdayDelivery              bool            `json:"saturday_delivery"`
	IgnoreAddressValidationErrors bool            `json:"ignore_address_validation_errors"`
	SkipAddressValidation         bool            `json:"skip_address_validation"`
	AllocationPriority            int             `json:"allocation_priority"`
	AllowPartial                  bool            `json:"allow_partial"`
	GiftInvoice                   bool            `json:"gift_invoice"`
	RequireSignature              bool            `json:"require_signature"`
	AdultSignatureRequired        bool            `json:"adult_signature_required"`
	Alcohol                       bool            `json:"alcohol"`
	Insurance                     bool            `json:"insurance"`
	InsuranceValue                float64         `json:"insurance_value"`
	Currency                      string          `json:"currency"`
	HasDryIce                     bool            `json:"has_dry_ice"`
	DryIceWeightInLbs             float64         `json:"dry_ice_weight_in_lbs"`
	AllowSplit                    bool            `json:"allow_split"`
	FTRExemption                  bool            `json:"ftr_exemption"`

	Box            *BoxReturnJSON           `json:"box,omitempty"`
	Store          *StoreReturnJSON         `json:"store,omitempty"`
	BillToAddress  AddressReturnJSON        `json:"bill_to_address"`
	ShipToAddress  AddressReturnJSON        `json:"ship_to_address"`
	ShippingMethod ShippingMethodReturnJSON `json:"shipping_method,omitempty"`
	Warehouse      WarehouseReturnJSON      `json:"warehouse"`

	OrderItems []OrderItemReturnJSON `json:"order_items,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrdersListRequest struct {
	StoreID        int    `json:"store_id"`
	ClientID       int    `json:"client_id"`
	OrganizationID int    `json:"organization_id"`
	Limit          int    `json:"limit"`
	Offset         int    `json:"offset"`
	OrderBy        string `json:"order_by"`
	OrderByColumn  string `json:"order_by_column"`
	SearchValue    string `json:"search_value"`
}

func (o *Order) Create(ctx context.Context) error {
	return util.DBFromContext(ctx).Create(o).Error
}

func (o *Order) GetOrderItems(ctx context.Context) error {
	return util.DBFromContext(ctx).Model(o).Association("OrderItems").Find(&o.OrderItems)
}

func (o *Order) GetStore(ctx context.Context) error {
	return util.DBFromContext(ctx).Model(o).Association("Store").Find(&o.Store)
}

func (o *Order) GetWarehouse(ctx context.Context) error {
	return util.DBFromContext(ctx).Model(o).Association("Warehouse").Find(&o.Warehouse)
}

func (o *Order) GetBillToAddress(ctx context.Context) error {
	return util.DBFromContext(ctx).Model(o).Association("BillToAddress").Find(&o.BillToAddress)
}

func (o *Order) GetShipToAddress(ctx context.Context) error {
	return util.DBFromContext(ctx).Model(o).Association("ShipToAddress").Find(&o.ShipToAddress)
}

func (o *Order) GetShippingMethod(ctx context.Context) error {
	return util.DBFromContext(ctx).Model(o).Association("ShippingMethod").Find(&o.ShippingMethod)
}

func (o *Order) GetBox(ctx context.Context) error {
	return util.DBFromContext(ctx).Model(o).Association("Box").Find(&o.Box)
}

func GetNextOrderReadyForPicking(ctx context.Context, warehouseID int) (*Order, error) {

	subquery := util.DBFromContext(ctx).Model(&Inventory{}).
		Select("product_id, location_id, COUNT(*) AS available_count").
		Where("order_item_id IS NULL").
		Group("product_id, location_id")

	var order Order
	err := util.DBFromContext(ctx).Model(&Order{}).
		Select("orders.*").
		Joins("JOIN order_items ON order_items.order_id = orders.id").
		Joins("JOIN locations ON locations.warehouse_id = ? AND locations.pickable = true AND locations.is_tote IS NOT TRUE", warehouseID).
		Joins("JOIN inventory ON inventory.product_id = order_items.product_id").
		Joins("JOIN (?) AS inventory_counts ON inventory_counts.product_id = order_items.product_id AND inventory_counts.location_id = locations.id AND inventory_counts.available_count >= order_items.allocated", subquery).
		Where("order_items.quantity = order_items.allocated").
		Where("NOT EXISTS (SELECT 1 FROM pick_session_orders WHERE pick_session_orders.order_id = orders.id)").
		Limit(1).
		Find(&order).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.New("failed to get next order ready for picking")
	}

	if err == gorm.ErrRecordNotFound || order.ID == 0 {
		return nil, errors.New("no orders ready for picking")
	}

	return &order, nil
}

func GetOrderByStoreAndAPIID(ctx context.Context, storeID int, apiID string) (*Order, error) {
	order := Order{}
	err := util.DBFromContext(ctx).Where("store_id = ? AND api_id = ?", storeID, apiID).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func GetOrderByID(ctx context.Context, id int) (*Order, error) {
	order := Order{}
	err := util.DBFromContext(ctx).Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (o *Order) AddHold(holdType uint64) {
	o.Holds |= holdType
}

func (o *Order) RemoveHold(holdType uint64) {
	o.Holds &= ^holdType
}

func (o *Order) HasHold(holdType uint64) bool {
	return o.Holds&holdType != 0
}

func (o *Order) GetShippingAddress(ctx context.Context) error {
	//get address by id
	address := Address{}
	err := util.DBFromContext(ctx).Where("id = ?", o.ShipToAddressID).First(&address).Error
	if err != nil {
		return err
	}
	o.ShipToAddress = address
	return nil
}

func (o *Order) Update(ctx context.Context) error {
	return util.DBFromContext(ctx).Save(o).Error
}

// TODO delete this
func (o *Order) ConvertToProductReturnJSON() *OrderReturnJSON {

	// TODO add status
	//get status
	// status, _ := GetPurchaseOrderStatusByID(p.Status)

	return &OrderReturnJSON{
		ID:                            o.ID,
		Priority:                      o.Priority,
		OrderNumber:                   o.OrderNumber,
		GiftNote:                      o.GiftNote,
		PackingNote:                   o.PackingNote,
		ReadyToShip:                   o.ReadyToShip,
		Holds:                         o.Holds,
		Subtotal:                      o.Subtotal,
		Tax:                           o.Tax,
		Shipping:                      o.Shipping,
		Discount:                      o.Discount,
		DiscountCodes:                 o.DiscountCodes,
		Tip:                           o.Tip,
		Total:                         o.Total,
		OrderDate:                     o.OrderDate,
		RequiredShipDate:              o.RequiredShipDate,
		HoldUntilDate:                 o.HoldUntilDate,
		FulfillmentStatus:             o.FulfillmentStatus,
		MarketplaceFinacialStatus:     o.MarketplaceFinacialStatus,
		AutoPrintReturnLabel:          o.AutoPrintReturnLabel,
		CustomerEmail:                 o.CustomerEmail,
		CustomerPhone:                 o.CustomerPhone,
		SaturdayDelivery:              o.SaturdayDelivery,
		IgnoreAddressValidationErrors: o.IgnoreAddressValidationErrors,
		SkipAddressValidation:         o.SkipAddressValidation,
		AllocationPriority:            o.AllocationPriority,
		AllowPartial:                  o.AllowPartial,
		GiftInvoice:                   o.GiftInvoice,
		RequireSignature:              o.RequireSignature,
		AdultSignatureRequired:        o.AdultSignatureRequired,
		Alcohol:                       o.Alcohol,
		Insurance:                     o.Insurance,
		InsuranceValue:                o.InsuranceValue,
		Currency:                      o.Currency,
		HasDryIce:                     o.HasDryIce,
		DryIceWeightInLbs:             o.DryIceWeightInLbs,
		AllowSplit:                    o.AllowSplit,
		FTRExemption:                  o.FTRExemption,
	}
}

func (o *Order) ConvertToReturnJSON(ctx context.Context) *OrderReturnJSON {

	//make sure order items is an empty array if it is nil
	if o.OrderItems == nil {
		o.OrderItems = []OrderItem{}
	}

	orderItems := []OrderItemReturnJSON{}
	for _, item := range o.OrderItems {
		// set to 0 and nil to omit from json
		item.OrderID = 0
		item.ProductID = nil
		orderItems = append(orderItems, item.ConvertToReturnJSON(ctx))
	}

	// TODO add histories
	// //make sure notes is an empty array if it is nil
	// if p.History == nil {
	// 	p.History = []PurchaseOrderHistory{}
	// }

	// PurchaseOrderHistorys := []PurchaseOrderHistoryReturnJSON{}
	// for _, note := range p.History {
	// 	note.PurchaseOrderID = 0
	// 	PurchaseOrderHistorys = append(PurchaseOrderHistorys, note.ConvertToReturnJSON())
	// }

	// TODO add tags
	//convert purchase order tags to slice of strings
	// tags := []string{}
	// if p.Tags != nil {
	// 	for _, tag := range p.Tags {
	// 		tags = append(tags, tag.Tag)
	// 	}
	// }

	// TODO add status
	//get status
	// status, _ := GetPurchaseOrderStatusByID(p.Status)

	// TODO validation
	//get store
	o.GetStore(ctx)

	// Get Billing Address
	o.GetBillToAddress(ctx)

	// Get Shipping Address
	o.GetShippingAddress(ctx)

	// Get Warehouse
	if o.WarehouseID != nil {
		o.GetWarehouse(ctx)
		o.Warehouse.GetReturnAddress(ctx)
		o.Warehouse.GetShipFromAddress(ctx)
	}

	// Get Shipping Method
	o.GetShippingMethod(ctx)

	// Get Box
	o.GetBox(ctx)

	return &OrderReturnJSON{
		ID:                            o.ID,
		Priority:                      o.Priority,
		OrderNumber:                   o.OrderNumber,
		GiftNote:                      o.GiftNote,
		PackingNote:                   o.PackingNote,
		ReadyToShip:                   o.ReadyToShip,
		Holds:                         o.Holds,
		Subtotal:                      o.Subtotal,
		Tax:                           o.Tax,
		Shipping:                      o.Shipping,
		Discount:                      o.Discount,
		DiscountCodes:                 o.DiscountCodes,
		Tip:                           o.Tip,
		Total:                         o.Total,
		OrderDate:                     o.OrderDate,
		RequiredShipDate:              o.RequiredShipDate,
		HoldUntilDate:                 o.HoldUntilDate,
		FulfillmentStatus:             o.FulfillmentStatus,
		MarketplaceFinacialStatus:     o.MarketplaceFinacialStatus,
		AutoPrintReturnLabel:          o.AutoPrintReturnLabel,
		CustomerEmail:                 o.CustomerEmail,
		CustomerPhone:                 o.CustomerPhone,
		SaturdayDelivery:              o.SaturdayDelivery,
		IgnoreAddressValidationErrors: o.IgnoreAddressValidationErrors,
		SkipAddressValidation:         o.SkipAddressValidation,
		AllocationPriority:            o.AllocationPriority,
		AllowPartial:                  o.AllowPartial,
		GiftInvoice:                   o.GiftInvoice,
		RequireSignature:              o.RequireSignature,
		AdultSignatureRequired:        o.AdultSignatureRequired,
		Alcohol:                       o.Alcohol,
		Insurance:                     o.Insurance,
		InsuranceValue:                o.InsuranceValue,
		Currency:                      o.Currency,
		HasDryIce:                     o.HasDryIce,
		DryIceWeightInLbs:             o.DryIceWeightInLbs,
		AllowSplit:                    o.AllowSplit,
		FTRExemption:                  o.FTRExemption,
		Store:                         o.Store.ConvertToReturnJSON(),
		ShipToAddress:                 *o.ShipToAddress.ConvertToReturnJSON(),
		BillToAddress:                 *o.BillToAddress.ConvertToReturnJSON(),
		Warehouse:                     *o.Warehouse.ConvertToReturnJSON(),
		ShippingMethod:                o.ShippingMethod.ConvertToReturnJSON(),
		Box:                           o.Box.ConvertToReturnJSON(),
		OrderItems:                    orderItems,
	}
}

func (olr *OrdersListRequest) ParseAndValidateRequest(r *http.Request) []string {

	errors := []string{}

	storeID, err := util.GetIntQueryParam(r, "store_id")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "store_id must be an integer")
		} else {
			olr.StoreID = storeID
		}
	}

	clientID, err := util.GetIntQueryParam(r, "client_id")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "client_id must be an integer")
		} else {
			olr.ClientID = clientID
		}
	}

	olr.Limit = 100
	limit, err := util.GetIntQueryParam(r, "limit")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "limit must be an integer")
		} else if limit < 0 {
			errors = append(errors, "limit must be greater than or equal to 0")
		} else {
			olr.Limit = limit
		}
	}

	olr.Offset = 0
	offset, err := util.GetIntQueryParam(r, "offset")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "offset must be an integer")
		} else if offset < 0 {
			errors = append(errors, "offset must be greater than or equal to 0")
		} else {
			olr.Offset = offset
		}
	}

	olr.OrderBy = "asc"
	orderBy, err := util.GetStringQueryParam(r, "order_by")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "order_by must be a string")
		} else if orderBy != "asc" && orderBy != "desc" {
			errors = append(errors, "order_by must be 'asc' or 'desc'")
		} else {
			olr.OrderBy = orderBy
		}
	}

	orderByColumn, err := util.GetStringQueryParam(r, "order_by_column")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "order_by_column must be a string")
		} else {
			olr.OrderByColumn = orderByColumn
		}
	}

	searchValue, err := util.GetStringQueryParam(r, "search_value")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "search_value must be a string")
		} else {
			olr.SearchValue = searchValue
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil

}

func (olr *OrdersListRequest) ConvertToOrganizationQuery(ctx context.Context) *gorm.DB {

	query := util.DBFromContext(ctx).Model(&Order{}).
		Select("DISTINCT orders.*").
		Joins("LEFT JOIN order_items ON order_items.order_id = orders.id").
		Joins("LEFT JOIN products ON products.id = order_items.product_id").
		Joins("LEFT JOIN stores ON stores.id = orders.store_id").
		Joins("LEFT JOIN clients ON clients.id = stores.client_id").
		Joins("LEFT JOIN organizations ON organizations.id = clients.organization_id")

	query = query.Where("organizations.id = ?", olr.OrganizationID)

	if olr.ClientID != 0 {
		query = query.Where("stores.client_id = ?", olr.ClientID)
	}

	// TODO - add the rest of the filters
	if olr.SearchValue != "" {
		query = query.Where(util.DBFromContext(ctx).Where("to_tsvector('english', order_number) @@ to_tsquery('english', ?)", fmt.Sprintf("*%s:*", olr.SearchValue)).
			Or(util.DBFromContext(ctx).Where("to_tsvector('english', orders.customer_email) @@ to_tsquery('english', ?)", fmt.Sprintf("*%s:*", olr.SearchValue))).
			Or(util.DBFromContext(ctx).Where("to_tsvector('english', orders.customer_phone) @@ to_tsquery('english', ?)", fmt.Sprintf("*%s:*", olr.SearchValue))).
			Or(util.DBFromContext(ctx).Where("to_tsvector('english', order_statuses.name) @@ to_tsquery('english', ?)", fmt.Sprintf("*%s:*", olr.SearchValue))))
	}

	if olr.OrderByColumn != "" {
		query = query.Order(olr.OrderByColumn + " " + olr.OrderBy)
	} else {
		if olr.OrderBy == "asc" {
			query = query.Order("id asc")
		} else {
			query = query.Order("id desc")
		}
	}

	return query
}

func (olr *OrdersListRequest) ConvertToClientQuery(ctx context.Context) *gorm.DB {

	query := util.DBFromContext(ctx).Model(&PurchaseOrder{}).
		Select("DISTINCT orders.*").
		Joins("LEFT JOIN order_items ON order_items.order_id = orders.id").
		Joins("LEFT JOIN products ON products.id = order_items.product_id").
		Joins("LEFT JOIN stores ON stores.id = orders.store_id").
		Joins("LEFT JOIN clients ON clients.id = stores.client_id")

	query = query.Where("client.id = ?", olr.ClientID)

	// TODO - add the rest of the filters
	if olr.SearchValue != "" {
		query = query.Where(util.DBFromContext(ctx).Where("to_tsvector('english', order_number) @@ to_tsquery('english', ?)", fmt.Sprintf("*%s:*", olr.SearchValue)).
			Or(util.DBFromContext(ctx).Where("to_tsvector('english', orders.customer_email) @@ to_tsquery('english', ?)", fmt.Sprintf("*%s:*", olr.SearchValue))).
			Or(util.DBFromContext(ctx).Where("to_tsvector('english', orders.customer_phone) @@ to_tsquery('english', ?)", fmt.Sprintf("*%s:*", olr.SearchValue))).
			Or(util.DBFromContext(ctx).Where("to_tsvector('english', order_statuses.name) @@ to_tsquery('english', ?)", fmt.Sprintf("*%s:*", olr.SearchValue))))
	}

	if olr.OrderByColumn != "" {
		query = query.Order(olr.OrderByColumn + " " + olr.OrderBy)
	} else {
		if olr.OrderBy == "asc" {
			query = query.Order("id asc")
		} else {
			query = query.Order("id desc")
		}
	}

	return query
}

func ConvertOrdersToSearchResults(ctx context.Context, matchingOrders []Order, total int, count int) (*SearchResults, error) {

	var orders []*OrderReturnJSON

	for _, order := range matchingOrders {
		orders = append(orders, order.ConvertToReturnJSON(ctx))
	}

	results, err := util.ConvertStructsToInterfaces(orders)
	if err != nil {
		return nil, err
	}

	searchResults := &SearchResults{
		TotalCount:    total,
		FilteredCount: count,
		Data:          results,
	}

	return searchResults, nil
}

func GetOrdersByProductID(ctx context.Context, productID int) ([]Order, error) {

	var orders []Order

	err := util.DBFromContext(ctx).Model(&Order{}).
		Select("DISTINCT orders.*").
		Joins("LEFT JOIN order_items ON order_items.order_id = orders.id").
		Where("order_items.product_id = ?", productID).
		Find(&orders).Error

	if err != nil {
		return nil, err
	}

	return orders, nil
}
