package models

import (
	"context"

	"github.com/shipply-io/shipply-io-backend/util"
)

type ProductBundle struct {
	ID        int `json:"id"`
	ClientID  int `json:"client_id"`
	ProductID int `json:"product_id"`

	Components []Product `gorm:"many2many:bundle_products;"`
	Product    Product   `gorm:"foreignkey:ProductID;" json:"product"`
}

type ProductBundleReturnJSON struct {
	ID       int `json:"id"`
	ClientID int `json:"client_id"`

	Components []ProductReturnJSON `json:"components"`
	Product    ProductReturnJSON   `json:"product"`
}

func (pb *ProductBundle) ConvertToReturnJSON(ctx context.Context) *ProductBundleReturnJSON {

	err := pb.GetComponents(ctx)
	if err != nil {
		return nil
	}

	components := []ProductReturnJSON{}
	for i := range pb.Components {
		components = append(components, *pb.Components[i].ConvertToReturnJSON())
	}

	err = pb.GetProduct(ctx)
	if err != nil {
		return nil
	}

	return &ProductBundleReturnJSON{
		ID:         pb.ID,
		ClientID:   pb.ClientID,
		Components: components,
		Product:    *pb.Product.ConvertToReturnJSON(),
	}
}

func (pb *ProductBundle) GetComponents(ctx context.Context) error {

	err := util.DBFromContext(ctx).Model(pb).Association("Components").Find(&pb.Components)
	if err != nil {
		return ErrQueryFailed{Err: err, Object: "product bundle components"}
	}

	return nil

}

func (pb *ProductBundle) GetProduct(ctx context.Context) error {

	err := util.DBFromContext(ctx).Model(pb).Association("Product").Find(&pb.Product)
	if err != nil {
		return ErrQueryFailed{Err: err, Object: "product"}
	}

	return nil

}

func GetBundleByProductID(ctx context.Context, productID int) (*ProductBundle, error) {

	bundle := ProductBundle{}

	err := util.DBFromContext(ctx).Where("product_id = ?", productID).First(&bundle).Error
	if err != nil {
		return nil, ErrQueryFailed{Err: err, Object: "product bundle"}
	}

	return &bundle, nil

}
