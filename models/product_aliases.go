package models

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/shipply-io/shipply-io-backend/util"
)

// override table name for product aliases
func (ProductAlias) TableName() string {
	return "products_aliases"
}

type ProductAlias struct {
	ID        int
	ProductID int
	Barcode   string
	Quantity  int
}

func (pa *ProductAlias) Create(ctx context.Context) error {
	return util.DBFromContext(ctx).Create(pa).Error
}

func (pa *ProductAlias) Delete(ctx context.Context) error {
	return util.DBFromContext(ctx).Delete(pa).Error
}

type ProductAliasCreateRequest struct {
	Barcode   string `json:"barcode"`
	ProductID int    `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type ProductAliasUpdateRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

func (pa *ProductAliasCreateRequest) ParseAndValidateRequest(r *http.Request) []string {
	var errs []string

	err := json.NewDecoder(r.Body).Decode(&pa)
	if err != nil {
		return []string{"invalid JSON"}
	}

	if pa.ProductID == 0 {
		errs = append(errs, "product_id is required")
	}

	if pa.Barcode == "" {
		errs = append(errs, "barcode is required")
	}

	if pa.Quantity == 0 {
		errs = append(errs, "quantity is required")
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func (pa *ProductAlias) UpdateWithRequest(ctx context.Context, paur *ProductAliasUpdateRequest) error {
	if paur.ProductID != 0 {
		pa.ProductID = paur.ProductID
	}
	if paur.Quantity != 0 {
		pa.Quantity = paur.Quantity
	}

	if err := util.DBFromContext(ctx).Save(pa).Error; err != nil {
		return err
	}

	return nil
}

func IsProductAliasBarcodeUnique(ctx context.Context, clientID int, barcode string) bool {
	var count int64
	util.DBFromContext(ctx).Model(&ProductAlias{}).Where("barcode = ?", barcode).Count(&count)

	// Check if products with the same barcode and client exist in main product table
	if count == 0 {
		util.DBFromContext(ctx).Model(&Product{}).Where("client_id = ? AND barcode = ?", clientID, barcode).Count(&count)
	}

	return count == 0
}

func GetProductAliasByBarcode(ctx context.Context, barcode string) (*ProductAlias, error) {
	productAlias := &ProductAlias{}
	err := util.DBFromContext(ctx).Where("barcode = ?", barcode).First(productAlias).Error
	return productAlias, err
}
