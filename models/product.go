package models

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"
)

/* --------------------------- Main Product Model --------------------------- */
type Product struct {
	ID                   int
	Name                 string
	Sku                  string `gorm:"uniqueIndex:idx_client_sku"`
	Barcode              string
	Kit                  bool
	Active               bool
	Weight               float64
	WeightUnit           string
	Grams                int
	Height               float64
	Width                float64
	Length               float64
	Value                float64
	FinalSale            bool
	NoAir                bool
	ValueCurrency        string
	Price                float64
	PriceCurrency        string
	CustomsValue         float64
	CustomsDescription   string
	ReorderLevel         int
	LastCounted          time.Time
	CountryOfManufacture string
	TariffCode           string
	NeedsSerialNumber    bool
	NeedsLotNumber       bool
	IgnoreOnInvoice      bool
	IgnoreOnCustoms      bool
	Virtual              bool
	ClientID             int `gorm:"uniqueIndex:idx_client_sku"`
	ProductNote          string
	PackNote             string
	ReturnNote           string
	Description          string

	OnHand      int
	NonSellable int
	Allocated   int
	Available   int
	Backordered int
	Reserve     int
	SellAhead   int
	Inbound     int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	Client             Client
	ProductLots        []ProductLot               `gorm:"-"`
	InventoryLocations []ProductInventoryLocation `gorm:"-"`
	Aliases            []ProductAlias             `gorm:"-"`
	InventoryHistory   []InventoryAuditLog        `gorm:"-"`
	ProductBundles     []ProductBundle            `gorm:"many2many:bundle_products;"`
	Components         []Product                  `gorm:"-"`
	Stores             []Store                    `gorm:"-"`
	ProductImages      []ProductImage             `gorm:"-"`
}

func (p *Product) GetProductImages(ctx context.Context) error {
	var productImages []ProductImage
	err := util.DBFromContext(ctx).Where("product_id = ?", p.ID).Order(`"order" ASC`).Find(&productImages).Error
	if err != nil {
		return err
	}
	p.ProductImages = productImages
	return nil
}

/* ----------------------------- Create Product ----------------------------- */
type CreateProductInput struct {
	ClientID    int
	Sku         string
	Barcode     *string
	Name        string
	Value       float64
	Description *string

	Weight     *Weight
	Dimensions *Dimensions
}

func CreateProduct(ctx context.Context, input CreateProductInput) (*Product, error) {

	product := Product{
		ClientID: input.ClientID,
		Sku:      input.Sku,
		Name:     input.Name,
		Value:    input.Value,
	}

	if input.Barcode != nil {
		product.Barcode = *input.Barcode
	}

	if input.Description != nil {
		product.Description = *input.Description
	}

	if input.Weight != nil {
		product.Weight = input.Weight.Value
		product.WeightUnit = input.Weight.Unit
	}

	if input.Dimensions != nil {
		product.Length = input.Dimensions.Length
		product.Width = input.Dimensions.Width
		product.Height = input.Dimensions.Height
	}

	if err := util.DBFromContext(ctx).Create(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil

}

/* ------------------------------ Update Product ----------------------------- */
type UpdateProductInput struct {
	ID          int
	Sku         *string
	Barcode     *string
	Name        *string
	Value       *float64
	Description *string
	Weight      *Weight
	Dimensions  *Dimensions
}

func UpdateProduct(ctx context.Context, input UpdateProductInput) (*Product, error) {

	var product Product
	if err := util.DBFromContext(ctx).First(&product, input.ID).Error; err != nil {
		return nil, err
	}

	if input.Sku != nil {
		product.Sku = *input.Sku
	}

	if input.Barcode != nil {
		product.Barcode = *input.Barcode
	}

	if input.Name != nil {
		product.Name = *input.Name
	}

	if input.Value != nil {
		product.Value = *input.Value
	}

	if input.Description != nil {
		product.Description = *input.Description
	}

	if input.Weight != nil {
		product.Weight = input.Weight.Value
		product.WeightUnit = input.Weight.Unit
	}

	if input.Dimensions != nil {
		product.Length = input.Dimensions.Length
		product.Width = input.Dimensions.Width
		product.Height = input.Dimensions.Height
	}

	if err := util.DBFromContext(ctx).Save(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil

}

/* ----------------------------- End of new code ---------------------------- */

type ProductReturnJSON struct {
	ID                   int       `json:"id,omitempty"`
	Name                 string    `json:"name,omitempty"`
	Sku                  string    `json:"sku,omitempty"`
	Barcode              string    `json:"barcode,omitempty"`
	Kit                  bool      `json:"kit,omitempty"`
	Active               bool      `json:"active,omitempty"`
	Weight               float64   `json:"weight,omitempty"`
	WeightUnit           string    `json:"weight_unit,omitempty"`
	Grams                int       `json:"grams,omitempty"`
	Height               float64   `json:"height,omitempty"`
	Width                float64   `json:"width,omitempty"`
	Length               float64   `json:"length,omitempty"`
	Value                float64   `json:"value,omitempty"`
	FinalSale            bool      `json:"final_sale,omitempty"`
	NoAir                bool      `json:"no_air,omitempty"`
	ValueCurrency        string    `json:"value_currency,omitempty"`
	Price                float64   `json:"price,omitempty"`
	PriceCurrency        string    `json:"price_currency,omitempty"`
	CustomsValue         float64   `json:"customs_value,omitempty"`
	CustomsDescription   string    `json:"customs_description,omitempty"`
	ReorderLevel         int       `json:"reorder_level,omitempty"`
	LastCounted          time.Time `json:"last_counted,omitempty"`
	CountryOfManufacture string    `json:"country_of_manufacture,omitempty"`
	TariffCode           string    `json:"tariff_code,omitempty"`
	NeedsSerialNumber    bool      `json:"needs_serial_number,omitempty"`
	NeedsLotNumber       bool      `json:"needs_lot_number,omitempty"`
	IgnoreOnInvoice      bool      `json:"ignore_on_invoice,omitempty"`
	IgnoreOnCustoms      bool      `json:"ignore_on_customs,omitempty"`
	Virtual              bool      `json:"virtual,omitempty"`
	ClientID             int       `json:"client_id,omitempty"`
	ProductNote          string    `json:"product_note,omitempty"`
	PackNote             string    `json:"pack_note,omitempty"`
	ReturnNote           string    `json:"return_note,omitempty"`
	ImageURL             string    `json:"image_url,omitempty"`
	OnHand               int       `json:"on_hand,omitempty"`
	NonSellable          int       `json:"non_sellable,omitempty"`
	Allocated            int       `json:"allocated,omitempty"`
	Available            int       `json:"available,omitempty"`
	Backordered          int       `json:"backordered,omitempty"`
	Reserve              int       `json:"reserve,omitempty"`
	SellAhead            int       `json:"sell_ahead,omitempty"`
	Inbound              int       `json:"inbound,omitempty"`
}

type ProductSearchRequest struct {
	SearchValue    string `json:"search_value"`
	ClientID       int    `json:"client_id"`
	OrganizationID int    `json:"organization_id"`
}

func (p *Product) Create(ctx context.Context) error {
	return util.DBFromContext(ctx).Create(p).Error
}

func (p *Product) Update(ctx context.Context) error {
	return util.DBFromContext(ctx).Save(p).Error
}

func (p *Product) ConvertToReturnJSON() *ProductReturnJSON {
	if p == nil {
		return &ProductReturnJSON{}
	}

	productReturnJSON := ProductReturnJSON{
		ID:                   p.ID,
		Name:                 p.Name,
		Sku:                  p.Sku,
		Barcode:              p.Barcode,
		Kit:                  p.Kit,
		Active:               p.Active,
		Weight:               p.Weight,
		WeightUnit:           p.WeightUnit,
		Grams:                p.Grams,
		Height:               p.Height,
		Width:                p.Width,
		Length:               p.Length,
		Value:                p.Value,
		FinalSale:            p.FinalSale,
		NoAir:                p.NoAir,
		ValueCurrency:        p.ValueCurrency,
		Price:                p.Price,
		PriceCurrency:        p.PriceCurrency,
		CustomsValue:         p.CustomsValue,
		CustomsDescription:   p.CustomsDescription,
		ReorderLevel:         p.ReorderLevel,
		LastCounted:          p.LastCounted,
		CountryOfManufacture: p.CountryOfManufacture,
		TariffCode:           p.TariffCode,
		NeedsSerialNumber:    p.NeedsSerialNumber,
		NeedsLotNumber:       p.NeedsLotNumber,
		IgnoreOnInvoice:      p.IgnoreOnInvoice,
		IgnoreOnCustoms:      p.IgnoreOnCustoms,
		Virtual:              p.Virtual,
		ClientID:             p.ClientID,
		ProductNote:          p.ProductNote,
		PackNote:             p.PackNote,
		ReturnNote:           p.ReturnNote,
		OnHand:               p.OnHand,
		NonSellable:          p.NonSellable,
		Allocated:            p.Allocated,
		Available:            p.Available,
		Backordered:          p.Backordered,
		Reserve:              p.Reserve,
		SellAhead:            p.SellAhead,
		Inbound:              p.Inbound,
	}
	return &productReturnJSON
}

func (p *ProductSearchRequest) ParseAndValidateRequest(r *http.Request) error {
	searchValue, err := util.GetStringQueryParam(r, "search_value")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			return errors.New("search_value must be of type string")
		} else {
			p.SearchValue = searchValue
		}
	}

	clientID, err := util.GetIntQueryParam(r, "client_id")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			return errors.New("customer_id must be of type int")
		} else {
			p.ClientID = clientID
		}
	}
	return nil
}

func GetProductByID(ctx context.Context, id int) (*Product, error) {
	var product Product
	err := util.DBFromContext(ctx).First(&product, id).Error
	if err != nil {
		return nil, ErrQueryFailed{Err: err, Object: "product"}
	}
	return &product, nil
}

func GetProductByClientIDAndSku(ctx context.Context, clientID int, sku string) (*Product, error) {
	var product Product
	err := util.DBFromContext(ctx).Where("client_id = ? AND sku = ?", clientID, sku).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func UpdateProductInventoryLevelsByProductID(ctx context.Context, productID int) error {
	product, err := GetProductByID(ctx, productID)
	if err != nil {
		return err
	}
	return product.UpdateProductInventoryLevels(ctx)
}

func (p *Product) UpdateProductInventoryLevels(ctx context.Context) error {
	db := util.DBFromContext(ctx)
	var onHand int64
	var nonSellable int64
	var allocated int64
	var backordered int64
	var inbound int64

	//on hand = total non shipped inventory
	err := db.Model(&Inventory{}).Where("product_id = ?", p.ID).Count(&onHand).Error
	if err != nil {
		return err
	}

	//non sellable = total non sellable inventory
	err = db.Model(&Inventory{}).Where("product_id = ? AND damaged = ?", p.ID, true).Count(&nonSellable).Error
	if err != nil {
		return err
	}

	//allocated = has order item id and not shipped
	var allocatedResult NullInt64
	err = db.Model(&OrderItem{}).Where("product_id = ?", p.ID).Select("SUM(allocated)").Scan(&allocatedResult).Error
	if err != nil {
		return err
	}
	allocated = allocatedResult.Int64

	//available = on hand - allocated - non sellable - reserved
	available := onHand - allocated - nonSellable

	//reserved = total reserved inventory
	if p.Reserve > 0 {
		allocatedResult.Int64 -= int64(p.Reserve)
		available -= int64(p.Reserve)
	}

	//backordered = total backordered inventory
	var backorderedResult NullInt64
	err = db.Model(&OrderItem{}).Where("product_id = ? AND backordered > 0", p.ID).Select("SUM(backordered)").Scan(&backorderedResult).Error
	if err != nil {
		return err
	}
	backordered = backorderedResult.Int64

	//inbound = total inbound inventory from purchase orders
	//get purchase orders that are not closed
	var purchaseOrders []PurchaseOrder
	err = db.Model(&PurchaseOrder{}).Where("client_id = ? AND closed IS NOT TRUE", p.ClientID).Find(&purchaseOrders).Error
	if err != nil {
		return err
	}

	for _, po := range purchaseOrders {
		var purchaseOrderItems []PurchaseOrderItem
		err = db.Model(&PurchaseOrderItem{}).Where("purchase_order_id = ? AND product_id = ?", po.ID, p.ID).Find(&purchaseOrderItems).Error
		if err != nil {
			return err
		}

		for _, poi := range purchaseOrderItems {
			if !(poi.Received+poi.Damaged > poi.Ordered) {
				inbound += int64(poi.Ordered) - int64(poi.Received) - int64(poi.Damaged)
			}
		}
	}

	//update product
	p.OnHand = int(onHand)
	p.NonSellable = int(nonSellable)
	p.Allocated = int(allocated)
	p.Available = int(available)
	p.Backordered = int(backordered)
	p.Inbound = int(inbound)

	err = db.Exec("UPDATE products SET on_hand = ?, non_sellable = ?, allocated = ?, available = ?, backordered = ?, inbound = ? WHERE id = ?", p.OnHand, p.NonSellable, p.Allocated, p.Available, p.Backordered, p.Inbound, p.ID).Error
	if err != nil {
		return err
	}

	return nil
}

func (p *Product) GetClient(ctx context.Context) error {
	var client Client
	err := util.DBFromContext(ctx).First(&client, p.ClientID).Error
	if err != nil {
		return ErrQueryFailed{Err: err, Object: "product's client association"}
	}

	p.Client = client
	return nil

}

func (p *Product) GetProductLots(ctx context.Context) error {
	var productLots []ProductLot
	err := util.DBFromContext(ctx).Where("product_id = ?", p.ID).Find(&productLots).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	p.ProductLots = productLots
	return nil
}

func (p *Product) UpdateIPA(ctx context.Context) error {
	err := util.DBFromContext(ctx).Exec("UPDATE products SET length = ?, width = ?, height = ?, weight = ?, weight_unit = ? WHERE id = ?", p.Length, p.Width, p.Height, p.Weight, p.WeightUnit, p.ID).Error
	if err != nil {
		return err
	}

	return nil
}

func (p *Product) IsBundle(ctx context.Context) (bool, error) {
	var productBundle ProductBundle
	err := util.DBFromContext(ctx).Where("product_id = ?", p.ID).First(&productBundle).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func GetProductByBarcodeAndClientID(ctx context.Context, barcode string, clientID int) (*Product, error) {
	var product Product
	err := util.DBFromContext(ctx).Where("client_id = ? AND barcode = ?", clientID, barcode).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

type ProductListRequest struct {
	ClientID       int    `json:"client_id"`
	OrganizationID int    `json:"organization_id"`
	Limit          int    `json:"limit"`
	Offset         int    `json:"offset"`
	OrderBy        string `json:"order_by"`
	OrderByColumn  string `json:"order_by_column"`
	SearchValue    string `json:"search_value"`

	*ProductListFilters
}

type ProductListFilters struct {
	Backorder *bool `json:"backorder"`
	Allocated *bool `json:"allocated"`
	Bundle    *bool `json:"bundle"`
	Active    *bool `json:"active"`
	Lot       *bool `json:"lot"`
}

func (plf *ProductListFilters) ParseAndValidateRequest(r *http.Request) []error {
	errs := []error{}

	backorder, err := util.GetOptionalBoolQueryParam(r, "backorder")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errs = append(errs, err)
		} else {
			plf.Backorder = backorder
		}
	}

	allocated, err := util.GetOptionalBoolQueryParam(r, "allocated")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errs = append(errs, err)
		} else {
			plf.Allocated = allocated
		}
	}

	bundle, err := util.GetOptionalBoolQueryParam(r, "bundle")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errs = append(errs, err)
		} else {
			plf.Bundle = bundle
		}
	}

	active, err := util.GetOptionalBoolQueryParam(r, "active")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errs = append(errs, err)
		} else {
			plf.Active = active
		}
	}

	lot, err := util.GetOptionalBoolQueryParam(r, "lot")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errs = append(errs, err)
		} else {
			plf.Lot = lot
		}
	}

	return errs
}

func (plr *ProductListRequest) ParseAndValidateRequest(r *http.Request) []string {

	errors := []string{}

	clientID, err := util.GetIntQueryParam(r, "client_id")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "client_id must be an integer")
		} else {
			plr.ClientID = clientID
		}
	}

	plr.Limit = 100
	limit, err := util.GetIntQueryParam(r, "limit")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "limit must be an integer")
		} else if limit < 0 {
			errors = append(errors, "limit must be greater than or equal to 0")
		} else {
			plr.Limit = limit
		}
	}

	plr.Offset = 0
	offset, err := util.GetIntQueryParam(r, "offset")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "offset must be an integer")
		} else if offset < 0 {
			errors = append(errors, "offset must be greater than or equal to 0")
		} else {
			plr.Offset = offset
		}
	}

	plr.OrderBy = "asc"
	orderBy, err := util.GetStringQueryParam(r, "order_by")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "order must be a string")
		} else if orderBy != "asc" && orderBy != "desc" {
			errors = append(errors, "order must be either asc or desc")
		} else {
			plr.OrderBy = orderBy
		}
	}

	searchValue, err := util.GetStringQueryParam(r, "search_value")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "search_value must be a string")
		} else {
			plr.SearchValue = searchValue
		}
	}

	orderByColumn, err := util.GetStringQueryParam(r, "order_by_column")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "order_by_column must be a string")
		} else {
			plr.OrderByColumn = orderByColumn
		}
	}

	saved_filter, err := util.GetIntQueryParam(r, "saved_filter")
	if err != util.ErrMissingQueryParam {
		if err != nil {
			errors = append(errors, "saved_filter must be an integer id")
		}
	}
	if saved_filter != 0 {
		user, err := GetRequestingUser(r)
		if err != nil {
			errors = append(errors, err.Error())
		} else {
			var filter ProductListFilters
			if err := LoadUserSavedFilter(r.Context(), user.ID, "products", saved_filter, &filter); err != nil {
				errors = append(errors, err.Error())
			}
			plr.ProductListFilters = &filter
		}
	}
	if plr.ProductListFilters == nil {
		plr.ProductListFilters = &ProductListFilters{}
		if errs := plr.ProductListFilters.ParseAndValidateRequest(r); len(errs) > 0 {
			for _, err := range errs {
				errors = append(errors, err.Error())
			}
		}
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

func (por *ProductListRequest) ConvertToOrganizationQuery(ctx context.Context) *gorm.DB {

	query := util.DBFromContext(ctx).Model(&Product{}).
		Select("DISTINCT products.*").
		Joins("LEFT JOIN clients ON clients.id = products.client_id").
		Joins("LEFT JOIN organizations ON organizations.id = clients.organization_id")

	query = query.Where("organizations.id = ?", por.OrganizationID)

	if por.ClientID != 0 {
		query = query.Where("products.client_id = ?", por.ClientID)
	}

	if por.SearchValue != "" {
		query = query.
			Where(util.DBFromContext(ctx).
				Where("to_tsvector('english', products.name) @@ to_tsquery('english', ?)", fmt.Sprintf("*%s:*", por.SearchValue)).
				Or("to_tsvector('english', products.sku) @@ to_tsquery('english', ?)", fmt.Sprintf("*%s:*", por.SearchValue)))
	}

	if por.OrderByColumn != "" {
		query = query.Order(por.OrderByColumn + " " + por.OrderBy)
	} else {
		if por.OrderBy == "asc" {
			query = query.Order("id asc")
		} else {
			query = query.Order("id desc")
		}
	}

	if por.ProductListFilters != nil {
		if por.ProductListFilters.Backorder != nil {
			if *por.ProductListFilters.Backorder {
				query = query.Where("products.backordered >= 1")
			} else {
				query = query.Where("products.backordered = 0")
			}
		}

		if por.ProductListFilters.Allocated != nil {
			if *por.ProductListFilters.Allocated {
				query = query.Where("products.allocated >= 1")
			} else {
				query = query.Where("products.allocated = 0")
			}
		}

		// TODO: Need to test this one throughly, I am not sure if this is correct
		// And setting up products for testing is a pain...
		if por.ProductListFilters.Bundle != nil {
			if *por.ProductListFilters.Bundle {
				query = query.Where("EXISTS (SELECT 1 from bundle_products where bundle_products.product_id = products.id)")
			} else {
				query = query.Where("NOT EXISTS (SELECT 1 from bundle_products where bundle_products.product_id = products.id)")
			}
		}

		if por.ProductListFilters.Active != nil {
			query = query.Where("products.active = ?", *por.ProductListFilters.Active)
		}

		if por.ProductListFilters.Lot != nil {
			query = query.Where("products.needs_lot_number = ?", *por.ProductListFilters.Lot)
		}
	}

	return query
}

func (por *ProductListRequest) ConvertToClientQuery(ctx context.Context) *gorm.DB {

	query := util.DBFromContext(ctx).Model(&Product{}).
		Select("DISTINCT products.*").
		Joins("LEFT JOIN clients ON clients.id = products.client_id")

	query = query.Where("products.client_id = ?", por.ClientID)

	if por.SearchValue != "" {
		query = query.
			Where(util.DBFromContext(ctx).
				Where("to_tsvector('english', products.name) @@ to_tsquery('english', ?)", fmt.Sprintf("*%s:*", por.SearchValue)).
				Or("to_tsvector('english', products.sku) @@ to_tsquery('english', ?)", fmt.Sprintf("*%s:*", por.SearchValue)))
	}

	if por.OrderByColumn != "" {
		query = query.Order(por.OrderByColumn + " " + por.OrderBy)
	} else {
		if por.OrderBy == "asc" {
			query = query.Order("id asc")
		} else {
			query = query.Order("id desc")
		}
	}

	if por.ProductListFilters != nil {
		if por.ProductListFilters.Backorder != nil {
			if *por.ProductListFilters.Backorder {
				query = query.Where("products.backordered >= 1")
			} else {
				query = query.Where("products.backordered = 0")
			}
		}

		if por.ProductListFilters.Allocated != nil {
			if *por.ProductListFilters.Allocated {
				query = query.Where("products.allocated >= 1")
			} else {
				query = query.Where("products.allocated = 0")
			}
		}

		// TODO: Need to test this one throughly, same issue as the org one
		if por.ProductListFilters.Bundle != nil {
			if *por.ProductListFilters.Bundle {
				query = query.Where("EXISTS (SELECT 1 from bundle_products where bundle_products.product_id = products.id)")
			} else {
				query = query.Where("NOT EXISTS (SELECT 1 from bundle_products where bundle_products.product_id = products.id)")
			}
		}

		if por.ProductListFilters.Active != nil {
			query = query.Where("products.active = ?", *por.ProductListFilters.Active)
		}

		if por.ProductListFilters.Lot != nil {
			query = query.Where("products.needs_lot_number = ?", *por.ProductListFilters.Lot)
		}
	}

	return query
}

func ConvertProductsToSearchResults(products []Product, total int, count int) (*SearchResults, error) {

	var productsJSON []*ProductReturnJSON

	for _, product := range products {
		productsJSON = append(productsJSON, product.ConvertToReturnJSON())
	}

	results, err := util.ConvertStructsToInterfaces(productsJSON)
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

func (client *Client) GetProducts(ctx context.Context, polr ProductListRequest) ([]Product, int, int, error) {

	query := polr.ConvertToClientQuery(ctx)
	countQuery := polr.ConvertToClientQuery(ctx)
	totalQuery := util.DBFromContext(ctx).Model(&Product{}).Where("client_id = ?", client.ID)

	var products []Product
	if err := query.Offset(polr.Offset).Limit(polr.Limit).Find(&products).Error; err != nil {
		return nil, 0, 0, err
	}

	// TODO: why are we doing a second query here for count when you can just use len(products)?
	var count int64
	if err := countQuery.Count(&count).Error; err != nil {
		return nil, 0, 0, err
	}

	var total int64
	if err := totalQuery.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	return products, int(count), int(total), nil
}

type ProductInventoryLocation struct {
	Quantity     int    `json:"quantity"`
	LocationName string `json:"location_name"`
	LocationID   int    `json:"location_id"`
	Pickable     bool   `json:"pickable"`
	Sellable     bool   `json:"sellable"`
}

type ProductInventoryLocationReturnJSON struct {
	Quantity     int    `json:"quantity"`
	LocationName string `json:"location_name"`
	LocationID   int    `json:"location_id"`
	Pickable     bool   `json:"pickable"`
	Sellable     bool   `json:"sellable"`
}

func (ProductInventoryLocation *ProductInventoryLocation) ConvertToReturnJSON() *ProductInventoryLocationReturnJSON {

	return &ProductInventoryLocationReturnJSON{
		Quantity:     ProductInventoryLocation.Quantity,
		LocationName: ProductInventoryLocation.LocationName,
		LocationID:   ProductInventoryLocation.LocationID,
		Pickable:     ProductInventoryLocation.Pickable,
		Sellable:     ProductInventoryLocation.Sellable,
	}

}

func (product *Product) GetInventoryLocations(ctx context.Context) error {

	var productInventoryLevels []ProductInventoryLocation

	err := util.DBFromContext(ctx).Raw(fmt.Sprintf(`
		SELECT COUNT
			( inventory.ID ) AS quantity,
			locations.NAME AS location_name,
			locations.ID AS location_id,
			locations.pickable,
			locations.sellable 
		FROM
			inventory
			INNER JOIN locations ON locations.ID = inventory.location_id 
		WHERE
			product_id = %d 
		GROUP BY
			locations.NAME,
			locations.pickable,
			locations.sellable,
			locations.id
	`, product.ID)).Scan(&productInventoryLevels).Error

	if err != nil {
		return err
	}

	product.InventoryLocations = productInventoryLevels

	return nil

}

func (product *Product) GetProductAliases(ctx context.Context) error {

	var productAliases []ProductAlias

	if err := util.DBFromContext(ctx).Where("product_id = ?", product.ID).Find(&productAliases).Error; err != nil {
		return err
	}

	product.Aliases = productAliases

	return nil

}

func (product *Product) GetInventoryHistory(ctx context.Context) error {

	var inventoryHistory []InventoryAuditLog

	if err := util.DBFromContext(ctx).Where("product_id = ?", product.ID).Find(&inventoryHistory).Error; err != nil {
		return err
	}

	product.InventoryHistory = inventoryHistory

	return nil
}

func (product *Product) GetLots(ctx context.Context) error {

	var lots []ProductLot

	if err := util.DBFromContext(ctx).Where("product_id = ?", product.ID).Find(&lots).Error; err != nil {
		return err
	}

	product.ProductLots = lots

	return nil
}

func (product *Product) GetProductBundles(ctx context.Context) error {

	err := util.DBFromContext(ctx).Model(product).Association("ProductBundles").Find(&product.ProductBundles)
	if err != nil {
		return err
	}

	return nil
}

func (product *Product) GetStores(ctx context.Context) error {

	var stores []Store
	if err := util.DBFromContext(ctx).
		Joins("JOIN shopify_products ON shopify_products.store_id = stores.id").
		Where("shopify_products.product_id = ?", product.ID).
		Find(&stores).Error; err != nil {
		return err
	}

	product.Stores = stores

	return nil

}

func (product *Product) IsComponent(ctx context.Context) bool {

	// find bundle product with product id
	var productBundle ProductBundle
	err := util.DBFromContext(ctx).
		Joins("JOIN bundle_products ON bundle_products.product_bundle_id = product_bundles.id").
		Where("bundle_products.product_id = ?", product.ID).First(&productBundle).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		}
		return false
	}

	return true

}
