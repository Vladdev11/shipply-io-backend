package models

import (
	"context"
	"encoding/json"
	"errors"
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
	if err := util.DBFromContext(ctx).Create(pa).Error; err != nil {
		return ErrCreateFailed{Err: err, Object: "product alias"}
	}
	return nil
}

func (pa *ProductAlias) Delete(ctx context.Context) error {
	if err := util.DBFromContext(ctx).Delete(pa).Error; err != nil {
		return ErrDeleteFailed{Err: err, Object: "product alias"}
	}
	return nil
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

func (pa *ProductAliasCreateRequest) ParseAndValidateRequest(r *http.Request) error {
	var errs []error

	//REVIEW- brady don't we need to confirm these fields are of correct type?
	err := json.NewDecoder(r.Body).Decode(&pa)
	if err != nil {
		return err
	}

	if pa.ProductID == 0 {
		errs = append(errs, ErrRequiredField{"product_id"})
	}

	if pa.Barcode == "" {
		errs = append(errs, ErrRequiredField{"barcode"})
	}

	if pa.Quantity == 0 {
		errs = append(errs, ErrRequiredField{"quantity"})
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
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
		return ErrUpdateFailed{Err: err, Object: "product alias"}
	}

	return nil
}

func EnsureProductAliasBarcodeUnique(ctx context.Context, clientID int, barcode string) error {
	var count int64
	if err := util.DBFromContext(ctx).Model(&ProductAlias{}).Where("barcode = ?", barcode).Count(&count).Error; err != nil {
		return ErrQueryFailed{Err: err, Object: "product alias by barcode"}
	}

	// Check if products with the same barcode and client exist in main product table
	if count == 0 {
		if err := util.DBFromContext(ctx).Model(&Product{}).Where("client_id = ? AND barcode = ?", clientID, barcode).Count(&count).Error; err != nil {
			return ErrQueryFailed{Err: err, Object: "product by client and barcode"}
		}
	}

	if count > 0 {
		return ErrExists{Object: "product or alias with this barcode"}
	}

	return nil
}

func GetProductAliasByBarcode(ctx context.Context, barcode string) (*ProductAlias, error) {
	productAlias := &ProductAlias{}
	if err := util.DBFromContext(ctx).Where("barcode = ?", barcode).First(productAlias).Error; err != nil {
		return nil, ErrQueryFailed{Err: err, Object: "product alias by barcode"}
	}

	return productAlias, nil
}

type ProductAliasReturnJSON struct {
	ID        int    `json:"id"`
	ProductID int    `json:"product_id"`
	Barcode   string `json:"barcode"`
	Quantity  int    `json:"quantity"`
}

func (pa *ProductAlias) ConvertToReturnJSON() *ProductAliasReturnJSON {
	return &ProductAliasReturnJSON{
		ID:        pa.ID,
		ProductID: pa.ProductID,
		Barcode:   pa.Barcode,
		Quantity:  pa.Quantity,
	}
}
