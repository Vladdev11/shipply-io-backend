package models

import (
	"errors"
	"net/http"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"
)

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
	ImageURL             string

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

	Client      Client
	ProductLots []ProductLot `gorm:"-"`
}

type ProductReturnJSON struct {
	ID                   int       `json:"id"`
	Name                 string    `json:"name"`
	Sku                  string    `json:"sku"`
	Barcode              string    `json:"barcode"`
	Kit                  bool      `json:"kit"`
	Active               bool      `json:"active"`
	Weight               float64   `json:"weight"`
	WeightUnit           string    `json:"weight_unit"`
	Grams                int       `json:"grams"`
	Height               float64   `json:"height"`
	Width                float64   `json:"width"`
	Length               float64   `json:"length"`
	Value                float64   `json:"value"`
	FinalSale            bool      `json:"final_sale"`
	NoAir                bool      `json:"no_air"`
	ValueCurrency        string    `json:"value_currency"`
	Price                float64   `json:"price"`
	PriceCurrency        string    `json:"price_currency"`
	CustomsValue         float64   `json:"customs_value"`
	CustomsDescription   string    `json:"customs_description"`
	ReorderLevel         int       `json:"reorder_level"`
	LastCounted          time.Time `json:"last_counted"`
	CountryOfManufacture string    `json:"country_of_manufacture"`
	TariffCode           string    `json:"tariff_code"`
	NeedsSerialNumber    bool      `json:"needs_serial_number"`
	NeedsLotNumber       bool      `json:"needs_lot_number"`
	IgnoreOnInvoice      bool      `json:"ignore_on_invoice"`
	IgnoreOnCustoms      bool      `json:"ignore_on_customs"`
	Virtual              bool      `json:"virtual"`
	ClientID             int       `json:"client_id"`
	ProductNote          string    `json:"product_note"`
	PackNote             string    `json:"pack_note"`
	ImageURL             string    `json:"image_url"`
}

type ProductSearchRequest struct {
	SearchValue    string `json:"search_value"`
	ClientID       int    `json:"client_id"`
	OrganizationID int    `json:"organization_id"`
}

func (p *Product) Create() error {
	return PGDB.Create(p).Error
}

func (p *Product) Update() error {
	return PGDB.Save(p).Error
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
		ImageURL:             p.ImageURL,
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

func GetProductByID(id int) (*Product, error) {
	var product Product
	err := PGDB.First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func GetProductByClientIDAndSku(clientID int, sku string) (*Product, error) {
	var product Product
	err := PGDB.Where("client_id = ? AND sku = ?", clientID, sku).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func UpdateProductInventoryLevelsByProductID(productID int) error {
	product, err := GetProductByID(productID)
	if err != nil {
		return err
	}
	return product.UpdateProductInventoryLevels()
}

func (p *Product) UpdateProductInventoryLevels() error {

	var onHand int64
	var nonSellable int64
	var allocated int64
	var backordered int64
	var inbound int64

	//on hand = total non shipped inventory
	err := PGDB.Model(&Inventory{}).Where("product_id = ?", p.ID).Count(&onHand).Error
	if err != nil {
		return err
	}

	//non sellable = total non sellable inventory
	err = PGDB.Model(&Inventory{}).Where("product_id = ? AND damaged = ?", p.ID, true).Count(&nonSellable).Error
	if err != nil {
		return err
	}

	//allocated = has order item id and not shipped
	var allocatedResult NullInt64
	err = PGDB.Model(&OrderItem{}).Where("product_id = ?", p.ID).Select("SUM(allocated)").Scan(&allocatedResult).Error
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
	err = PGDB.Model(&OrderItem{}).Where("product_id = ? AND backordered > 0", p.ID).Select("SUM(backordered)").Scan(&backorderedResult).Error
	if err != nil {
		return err
	}
	backordered = backorderedResult.Int64

	//inbound = total inbound inventory from purchase orders
	//get purchase orders that are not closed
	var purchaseOrders []PurchaseOrder
	err = PGDB.Model(&PurchaseOrder{}).Where("client_id = ? AND closed IS NOT TRUE", p.ClientID).Find(&purchaseOrders).Error
	if err != nil {
		return err
	}

	for _, po := range purchaseOrders {
		var purchaseOrderItems []PurchaseOrderItem
		err = PGDB.Model(&PurchaseOrderItem{}).Where("purchase_order_id = ? AND product_id = ?", po.ID, p.ID).Find(&purchaseOrderItems).Error
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

	err = PGDB.Exec("UPDATE products SET on_hand = ?, non_sellable = ?, allocated = ?, available = ?, backordered = ?, inbound = ? WHERE id = ?", p.OnHand, p.NonSellable, p.Allocated, p.Available, p.Backordered, p.Inbound, p.ID).Error
	if err != nil {
		return err
	}

	return nil
}

func (p *Product) GetClient() error {
	var client Client
	err := PGDB.First(&client, p.ClientID).Error
	if err != nil {
		return err
	}

	p.Client = client
	return nil

}

func (p *Product) GetProductLots() error {
	var productLots []ProductLot
	err := PGDB.Where("product_id = ?", p.ID).Find(&productLots).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	p.ProductLots = productLots
	return nil
}

func (p *Product) UpdateIPA() error {
	err := PGDB.Exec("UPDATE products SET length = ?, width = ?, height = ?, weight = ?, weight_unit = ? WHERE id = ?", p.Length, p.Width, p.Height, p.Weight, p.WeightUnit, p.ID).Error
	if err != nil {
		return err
	}

	return nil
}

func GetProductByBarcodeAndClientID(barcode string, clientID int) (*Product, error) {
	var product Product
	err := PGDB.Where("client_id = ? AND barcode = ?", clientID, barcode).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}
