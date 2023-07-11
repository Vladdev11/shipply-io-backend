package models

import (
	"context"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
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

func (s *ShopifyProduct) Create(ctx context.Context) error {
	return util.DBFromContext(ctx).Create(s).Error
}

func (s *ShopifyProduct) Update(ctx context.Context) error {
	return util.DBFromContext(ctx).Save(s).Error
}

func GetShopifyProductByGraphqlID(ctx context.Context, graphqlID string) (*ShopifyProduct, error) {
	shopifyProduct := &ShopifyProduct{}
	err := util.DBFromContext(ctx).Where("variant_shopify_graphql_id = ?", graphqlID).First(shopifyProduct).Error
	if err != nil {
		return nil, err
	}

	return shopifyProduct, nil
}

func GetShopifyProductAndProductByGraphqlID(ctx context.Context, graphqlID string) (*ShopifyProduct, error) {
	shopifyProduct := &ShopifyProduct{}
	err := util.DBFromContext(ctx).Where("variant_shopify_graphql_id = ?", graphqlID).First(shopifyProduct).Error
	if err != nil {
		return nil, err
	}

	err = util.DBFromContext(ctx).Model(shopifyProduct).Association("Product").Find(&shopifyProduct.Product)
	if err != nil {
		return nil, err
	}

	return shopifyProduct, nil
}

func GetShopifyProductsByProductID(ctx context.Context, productID int) ([]*ShopifyProduct, error) {
	shopifyProducts := []*ShopifyProduct{}
	err := util.DBFromContext(ctx).Where("product_id = ?", productID).Find(&shopifyProducts).Error
	if err != nil {
		return nil, err
	}

	return shopifyProducts, nil
}
