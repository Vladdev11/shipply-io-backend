package models

import (
	"time"

	"gorm.io/gorm"
)

type ShopifyProduct struct {
	ID                            int    `json:"id"`
	VariantShopifyGraphqlID       string `json:"shopify_graphql_id"`
	ProductShopifyGraphqlID       string `json:"product_shopify_graphql_id"`
	ProductID                     int    `json:"product_id"`
	StoreID                       int    `json:"store_id"`
	Sku                           string `json:"sku"`
	InventoryPolicy               string `json:"inventory_policy"`
	InventoryItemShopifyGraphqlID string `json:"inventory_item_shopify_graphql_id"`
	ProductStatus                 string `json:"status"`
	ProductHandle                 string `json:"product_handle"`
	ProductType                   string `json:"product_type"`
	ProductTitle                  string `json:"product_title"`
	VariantTitle                  string `json:"variant_title"`
	ProductVendor                 string `json:"product_vendor"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	Store   Store    `json:"store"`
	Product *Product `json:"product"`
}

func (s *ShopifyProduct) Create() error {
	return PGDB.Create(s).Error
}

func (s *ShopifyProduct) Update() error {
	return PGDB.Save(s).Error
}

func GetShopifyProductByGraphqlID(graphqlID string) (*ShopifyProduct, error) {
	shopifyProduct := &ShopifyProduct{}
	err := PGDB.Where("variant_shopify_graphql_id = ?", graphqlID).First(shopifyProduct).Error
	if err != nil {
		return nil, err
	}

	return shopifyProduct, nil
}

func GetShopifyProductAndProductByGraphqlID(graphqlID string) (*ShopifyProduct, error) {
	shopifyProduct := &ShopifyProduct{}
	err := PGDB.Where("variant_shopify_graphql_id = ?", graphqlID).First(shopifyProduct).Error
	if err != nil {
		return nil, err
	}

	err = PGDB.Model(shopifyProduct).Association("Product").Find(&shopifyProduct.Product)
	if err != nil {
		return nil, err
	}

	return shopifyProduct, nil
}
