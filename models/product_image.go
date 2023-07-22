package models

import (
	"context"
	"fmt"

	"github.com/shipply-io/shipply-io-backend/util"
)

/* ------------------------- Main ProductImage Model ------------------------ */
type ProductImage struct {
	ID        int    `json:"id"`
	ProductID int    `json:"product_id"`
	FileName  string `json:"file_name"`
	Order     int    `json:"order"`
}

func (productImage *ProductImage) GetImageURL(ctx context.Context) string {
	if productImage.FileName == "" {
		return ""
	}
	return fmt.Sprintf("%s/%s", util.CDNFromContext(ctx), productImage.FileName)
}

func GetProductImageByID(ctx context.Context, ID int) (*ProductImage, error) {
	productImage := ProductImage{}
	err := util.DBFromContext(ctx).Model(&ProductImage{}).Where("id = ?", ID).First(&productImage).Error
	if err != nil {
		return nil, err
	}
	return &productImage, nil
}

/* --------------------------- Create ProductImage --------------------------- */
type CreateProductImageInput struct {
	ProductID int    `json:"product_id"`
	FileName  string `json:"file_name"`
	Order     int    `json:"order"`
}

func CreateProductImage(ctx context.Context, input CreateProductImageInput) (*ProductImage, error) {

	productImage := ProductImage{
		ProductID: input.ProductID,
		FileName:  input.FileName,
		Order:     input.Order,
	}

	err := util.DBFromContext(ctx).Model(&ProductImage{}).Create(&productImage).Error
	if err != nil {
		return nil, err
	}

	return &productImage, nil
}

/* --------------------------- Update ProductImage --------------------------- */
type UpdateProductImageInput struct {
	ID    int `json:"id"`
	Order int `json:"order"`
}

func UpdateProductImage(ctx context.Context, input UpdateProductImageInput) (*ProductImage, error) {
	productImage := ProductImage{
		Order: input.Order,
	}

	err := util.DBFromContext(ctx).Model(&ProductImage{}).Where("id = ?", input.ID).Updates(&productImage).Error
	if err != nil {
		return nil, err
	}

	return &productImage, nil
}

/* --------------------------- Delete ProductImage --------------------------- */
func DeleteProductImage(ctx context.Context, ID int) error {
	err := util.DBFromContext(ctx).Model(&ProductImage{}).Where("id = ?", ID).Delete(&ProductImage{}).Error
	if err != nil {
		return err
	}
	return nil
}
