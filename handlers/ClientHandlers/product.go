package ClientHandlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/shipply-io/shipply-io-backend/api"
	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/responses"
	"github.com/shipply-io/shipply-io-backend/util"
	"github.com/shipply-io/shipply-io-backend/validation"
	"gorm.io/gorm"
)

var (
	//ErrProductIsBundle is returned when a product is a bundle
	ErrProductIsBundle = errors.New("product is a bundle")
	//ErrProductIsComponent is returned when a product is a component
	ErrProductIsComponent = errors.New("product is a component")
	//ErrProductBarcodeAlreadyExists is returned when a product with the provided barcode already exists
	ErrProductBarcodeAlreadyExists = errors.New("product with barcode already exists")
	//ErrProductSkuAlreadyExists is returned when a product with the provided sku already exists
	ErrProductSkuAlreadyExists = errors.New("product with sku already exists")
	//ErrDuplicateProductImageIDs is returned when duplicate product image ids are provided
	ErrDuplicateProductImageIDs = errors.New("duplicate product image ids provided")
	//ErrInvalidNumberOfProductImages is returned when the number of product images provided is invalid
	ErrInvalidNumberOfProductImages = errors.New("invalid number of product images provided")
	//ErrProductImageDoesNotBelongToProduct is returned when a product image does not belong to the product
	ErrProductImageDoesNotBelongToProduct = errors.New("product image does not belong to product")
)

func SearchProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := models.UserFromContext(ctx)

	// TODO - move this to validation package
	productSearch := models.ProductSearchRequest{}
	err := productSearch.ParseAndValidateRequest(r)
	if err != nil {
		util.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	productSearch.ClientID = user.OwnerID

	products, err := user.Client.SearchProducts(ctx, productSearch)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	response := responses.GenerateSearchProductsResponse(products)
	util.JSONResponse(w, response, http.StatusOK)

}

func ListProducts(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	user := models.UserFromContext(ctx)
	request := models.ProductListRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	if user.Client.ID != request.ClientID {
		util.ErrResponse(w, models.ErrUserDoesNotBelongToClient, http.StatusForbidden)
		return
	}

	products, count, total, err := user.Client.GetProducts(ctx, request)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	for i := range products {
		err = products[i].GetProductImages(ctx)
		if err != nil {
			util.ErrResponse(w, err, http.StatusBadRequest)
			return
		}
	}

	response := responses.GenerateListProductsResponse(ctx, products, count, total)
	util.JSONResponse(w, response, http.StatusOK)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	user := models.UserFromContext(ctx)

	getProductRequestData, err := validation.ParseRequestToGetProductRequestData(r)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	product, err := models.GetProductByID(ctx, getProductRequestData.ProductID)
	if err != nil {
		util.ErrResponse(w, ErrGetClient, http.StatusUnauthorized)
		return
	}

	productID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	product, err := models.GetProductByID(ctx, productID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if user.Client.ID != product.ClientID {
		util.ErrResponse(w, models.ErrUserDoesNotBelongToClient, http.StatusForbidden)
		return
	}

	err = product.GetProductImages(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	response := responses.GenerateGetProductResponse(ctx, product)
	util.JSONResponse(w, response, http.StatusOK)

}

func GetProductOrders(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	user := models.UserFromContext(ctx)

	getProductOrdersRequestData, err := validation.ParseRequestToGetProductOrdersRequestData(r)

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetClient(ctx)
	if err != nil {
		util.ErrResponse(w, ErrGetClient, http.StatusUnauthorized)
		return
	}

	productID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	product, err := models.GetProductByID(ctx, getProductOrdersRequestData.ProductID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if user.Client.ID != product.ClientID {
		util.ErrResponse(w, models.ErrUserDoesNotBelongToClient, http.StatusForbidden)
		return
	}

	orders, err := models.GetOrdersByProductID(ctx, product.ID)
	if err != nil {
		util.ErrorResponse(w, "failed to get orders", http.StatusBadRequest)
		return
	}

	response := responses.GenerateGetProductOrdersResponse(orders)
	util.JSONResponse(w, response, http.StatusOK)

}

func GetProductInventory(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	user := models.UserFromContext(ctx)

	getProductInventoryRequestData, err := validation.ParseRequestToGetProductInventoryRequestData(r)

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetClient(ctx)
	if err != nil {
		util.ErrResponse(w, ErrGetClient, http.StatusUnauthorized)
		return
	}

	productID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	product, err := models.GetProductByID(ctx, getProductInventoryRequestData.ProductID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if user.Client.ID != product.ClientID {
		util.ErrResponse(w, models.ErrUserDoesNotBelongToClient, http.StatusForbidden)
		return
	}

	//aliases
	err = product.GetProductAliases(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get product aliases", http.StatusBadRequest)
		return
	}

	//inventory history
	err = product.GetInventoryHistory(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get inventory history", http.StatusBadRequest)
		return
	}

	for i := range product.InventoryHistory {
		err = product.InventoryHistory[i].GetChangedByUser(ctx)
		if err != nil {
			util.ErrorResponse(w, "failed to get changed by user", http.StatusBadRequest)
			return
		}
	}

	//lots
	err = product.GetLots(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get lots", http.StatusBadRequest)
		return
	}

	response := responses.GenerateGetProductInventoryResponse(ctx, product)
	util.JSONResponse(w, response, http.StatusOK)

}

func GetProductBundles(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	user := models.UserFromContext(ctx)

	getProductBundlesRequestData, err := validation.ParseRequestToGetProductBundlesRequestData(r)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	product, err := models.GetProductByID(ctx, getProductBundlesRequestData.ProductID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	product, err := models.GetProductByID(ctx, productID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if user.Client.ID != product.ClientID {
		util.ErrResponse(w, models.ErrUserDoesNotBelongToClient, http.StatusForbidden)
		return
	}

	productIsBundle, err := product.IsBundle(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if productIsBundle {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	err = product.GetProductBundles(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	for i := range product.ProductBundles {
		err = product.ProductBundles[i].GetProduct(ctx)
		if err != nil {
			util.ErrResponse(w, err, http.StatusBadRequest)
			return
		}

		err = product.ProductBundles[i].Product.GetProductImages(ctx)
		if err != nil {
			util.ErrResponse(w, err, http.StatusBadRequest)
			return
		}
	}

	response := responses.GenerateGetProductBundlesResponse(ctx, product)
	util.JSONResponse(w, response, http.StatusOK)

}

func GetProductBundleComponents(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	user := models.UserFromContext(ctx)

	getProductBundleComponentsRequestData, err := validation.ParseRequestToGetProductBundleComponentsRequestData(r)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	product, err := models.GetProductByID(ctx, getProductBundleComponentsRequestData.ProductID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	product, err := models.GetProductByID(ctx, productID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if user.Client.ID != product.ClientID {
		util.ErrResponse(w, models.ErrUserDoesNotBelongToClient, http.StatusForbidden)
		return
	}

	if product.IsComponent(ctx) {
		util.ErrResponse(w, ErrProductIsComponent, http.StatusBadRequest)
		return
	}

	bundle, err := models.GetBundleByProductID(ctx, product.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response := responses.GenerateGetProductBundleComponentsResponse([]models.Product{})
			util.JSONResponse(w, response, http.StatusOK)
			return
		}

		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	err = bundle.GetComponents(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	response := responses.GenerateGetProductBundleComponentsResponse(bundle.Components)
	util.JSONResponse(w, response, http.StatusOK)

}

func GetProductStores(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	user := models.UserFromContext(ctx)

	getProductStoresRequestData, err := validation.ParseRequestToGetProductStoresRequestData(r)
	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetClient(ctx)
	if err != nil {
		util.ErrResponse(w, ErrGetClient, http.StatusUnauthorized)
		return
	}

	productID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	product, err := models.GetProductByID(ctx, getProductStoresRequestData.ProductID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if user.Client.ID != product.ClientID {
		util.ErrResponse(w, models.ErrUserDoesNotBelongToClient, http.StatusForbidden)
		return
	}

	err = product.GetStores(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	response := responses.GenerateGetProductStoresResponse(product)
	util.JSONResponse(w, response, http.StatusOK)

}

func CreateProduct(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	user := models.UserFromContext(ctx)

	createProductRequestData, err := validation.ParseRequestToCreateProductRequestData(r)
	if err != nil {
	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetClient(ctx)
	if err != nil {
		util.ErrResponse(w, ErrGetClient, http.StatusUnauthorized)
		return
	}

	err = user.Client.GetOrganization(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusUnauthorized)
		return
	}

	createProductRequestData, err := validation.ParseRequestToCreateProductRequestData(r)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	existingProductByClientIDAndSku, err := models.GetProductByClientIDAndSku(ctx, user.Client.ID, createProductRequestData.Sku)
	if err != nil && err != gorm.ErrRecordNotFound {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if existingProductByClientIDAndSku != nil {
		util.ErrResponse(w, models.ErrExists{Object: "product sku"}, http.StatusBadRequest)
		return
	}

	existingProductByBarcodeAndClientID, err := models.GetProductByBarcodeAndClientID(ctx, *createProductRequestData.Barcode, user.Client.ID)
	if err != nil && err != gorm.ErrRecordNotFound {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if existingProductByBarcodeAndClientID != nil {
		util.ErrResponse(w, ErrProductBarcodeAlreadyExists, http.StatusBadRequest)
		return
	}

	createProductInput := models.CreateProductInput{
		Name:        createProductRequestData.Name,
		ClientID:    user.Client.ID,
		Sku:         createProductRequestData.Sku,
		Barcode:     createProductRequestData.Barcode,
		Value:       createProductRequestData.Value,
		Description: createProductRequestData.Description,
	}

	if createProductRequestData.Weight != nil {
		createProductInput.Weight = &models.Weight{
			Value: createProductRequestData.Weight.Value,
			Unit:  createProductRequestData.Weight.Unit,
		}
	}

	if createProductRequestData.Dimensions != nil {
		createProductInput.Dimensions = &models.Dimensions{
			Length: createProductRequestData.Dimensions.Length,
			Width:  createProductRequestData.Dimensions.Width,
			Height: createProductRequestData.Dimensions.Height,
		}
	}

	product, err := models.CreateProduct(ctx, createProductInput)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if createProductRequestData.MainImage != nil {
		mainImageUUID := uuid.New().String()
		mainImageExtension := strings.Split(createProductRequestData.MainImage.FileType, "/")[1]

		err = api.S3FromContext(ctx).UploadFileToCDN(createProductRequestData.MainImage.ImageData, mainImageUUID, mainImageExtension, createProductRequestData.MainImage.FileType)
		if err != nil {
			util.ErrorResponse(w, "failed to upload image to s3", http.StatusInternalServerError)
			return
		}

		_, err := models.CreateProductImage(ctx, models.CreateProductImageInput{
			ProductID: product.ID,
			FileName:  fmt.Sprintf("%s.%s", mainImageUUID, mainImageExtension),
			Order:     0,
		})
		if err != nil {
			util.ErrResponse(w, err, http.StatusInternalServerError)
			return
		}
	}

	for i, image := range createProductRequestData.AdditionalImages {

		fileUUID := uuid.New().String()
		fileExtension := strings.Split(image.FileType, "/")[1]

		err = api.S3FromContext(ctx).UploadFileToCDN(image.ImageData, fileUUID, fileExtension, image.FileType)
		if err != nil {
			util.ErrorResponse(w, "failed to upload image to s3", http.StatusInternalServerError)
			return
		}

		_, err := models.CreateProductImage(ctx, models.CreateProductImageInput{
			ProductID: product.ID,
			FileName:  fmt.Sprintf("%s.%s", fileUUID, fileExtension),
			Order:     i + 1,
		})
		if err != nil {
			util.ErrResponse(w, err, http.StatusInternalServerError)
			return
		}

	}

	response := responses.GenerateCreateProductResponse(*product)
	util.JSONResponse(w, response, http.StatusOK)

}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	user := models.UserFromContext(ctx)

	updateProductRequestData, err := validation.ParseRequestToUpdateProductRequestData(r)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	product, err := models.GetProductByID(ctx, updateProductRequestData.ProductID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusUnauthorized)
		return
	}

	productID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	product, err := models.GetProductByID(ctx, productID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if product.ClientID != user.Client.ID {
		util.ErrResponse(w, models.ErrUserDoesNotBelongToClient, http.StatusForbidden)
		return
	}

	updateProductRequestData, err := validation.ParseRequestToUpdateProductRequestData(r)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if updateProductRequestData.Sku != nil {
		existingProductByClientIDAndSku, err := models.GetProductByClientIDAndSku(ctx, user.Client.ID, *updateProductRequestData.Sku)
		if err != nil && err != gorm.ErrRecordNotFound {
			util.ErrResponse(w, err, http.StatusBadRequest)
			return
		}

		if existingProductByClientIDAndSku != nil && existingProductByClientIDAndSku.ID != product.ID {
			util.ErrResponse(w, ErrProductSkuAlreadyExists, http.StatusBadRequest)
			return
		}
	}

	if updateProductRequestData.Barcode != nil {
		existingProductByBarcodeAndClientID, err := models.GetProductByBarcodeAndClientID(ctx, *updateProductRequestData.Barcode, user.Client.ID)
		if err != nil && err != gorm.ErrRecordNotFound {
			util.ErrResponse(w, err, http.StatusBadRequest)
			return
		}

		if existingProductByBarcodeAndClientID != nil && existingProductByBarcodeAndClientID.ID != product.ID {
			util.ErrResponse(w, ErrProductBarcodeAlreadyExists, http.StatusBadRequest)
			return
		}
	}

	updateProductInput := models.UpdateProductInput{
		ID:          product.ID,
		Name:        updateProductRequestData.Name,
		Sku:         updateProductRequestData.Sku,
		Barcode:     updateProductRequestData.Barcode,
		Value:       updateProductRequestData.Value,
		Description: updateProductRequestData.Description,
	}

	if updateProductRequestData.Weight != nil {
		updateProductInput.Weight = &models.Weight{
			Value: updateProductRequestData.Weight.Value,
			Unit:  updateProductRequestData.Weight.Unit,
		}
	}

	if updateProductRequestData.Dimensions != nil {
		updateProductInput.Dimensions = &models.Dimensions{
			Length: updateProductRequestData.Dimensions.Length,
			Width:  updateProductRequestData.Dimensions.Width,
			Height: updateProductRequestData.Dimensions.Height,
		}
	}

	product, err = models.UpdateProduct(ctx, updateProductInput)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	response := responses.GenerateUpdateProductResponse(*product)
	util.JSONResponse(w, response, http.StatusOK)
}

func AddProductImage(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	user := models.UserFromContext(ctx)

	addProductImageRequestData, err := validation.ParseRequestToAddProductImageRequestData(r)

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = user.GetClient(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusUnauthorized)
		return
	}

	productID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	product, err := models.GetProductByID(ctx, addProductImageRequestData.ProductID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if product.ClientID != user.Client.ID {
		util.ErrResponse(w, models.ErrUserDoesNotBelongToClient, http.StatusForbidden)
		return
	}

	err = product.GetProductImages(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	addProductImageRequestData, errors := validation.ParseRequestToAddProductImageRequestData(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	fileUUID := uuid.New().String()
	fileExtension := strings.Split(addProductImageRequestData.Image.FileType, "/")[1]

	err = api.S3FromContext(ctx).UploadFileToCDN(addProductImageRequestData.Image.ImageData, fileUUID, fileExtension, addProductImageRequestData.Image.FileType)
	if err != nil {
		util.ErrorResponse(w, "failed to upload image to s3", http.StatusInternalServerError)
		return
	}

	currentlyLargestOrder := 0
	for _, productImage := range product.ProductImages {
		if productImage.Order > currentlyLargestOrder {
			currentlyLargestOrder = productImage.Order
		}
	}

	productImage, err := models.CreateProductImage(ctx, models.CreateProductImageInput{
		ProductID: product.ID,
		FileName:  fmt.Sprintf("%s.%s", fileUUID, fileExtension),
		Order:     currentlyLargestOrder + 1,
	})
	if err != nil {
		util.ErrorResponse(w, "failed to create product image", http.StatusInternalServerError)
		return
	}

	response := responses.GenerateAddProductImageResponse(ctx, *productImage)
	util.JSONResponse(w, response, http.StatusOK)
}

func DeleteProductImage(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	user := models.UserFromContext(ctx)

	deleteProductImageRequestData, err := validation.ParseRequestToDeleteProductImageRequestData(r)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	product, err := models.GetProductByID(ctx, deleteProductImageRequestData.ProductID)
	if err != nil {
		util.ErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}

	productID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	product, err := models.GetProductByID(ctx, productID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if product.ClientID != user.Client.ID {
		util.ErrResponse(w, models.ErrUserDoesNotBelongToClient, http.StatusForbidden)
		return
	}

	productImage, err := models.GetProductImageByID(ctx, deleteProductImageRequestData.ProductImageID)
	if err != nil {
		util.ErrorResponse(w, "failed to find product image", http.StatusBadRequest)
		return
	}

	if productImage.ProductID != product.ID {
		util.ErrorResponse(w, "product image does not belong to this product", http.StatusForbidden)
		return
	}

	err = models.DeleteProductImage(ctx, productImage.ID)
	if err != nil {
		util.ErrorResponse(w, "failed to delete product image", http.StatusBadRequest)
		return
	}

	// Reorder product images
	err = product.GetProductImages(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	for i := range product.ProductImages {
		_, err = models.UpdateProductImage(ctx, models.UpdateProductImageInput{
			ID:    product.ProductImages[i].ID,
			Order: i,
		})

		if err != nil {
			util.ErrorResponse(w, "failed to update product image", http.StatusBadRequest)
			return
		}
	}

	util.SuccessResponse(w, http.StatusOK)
}

func UpdateProductImageOrder(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	user := models.UserFromContext(ctx)

	updateProductImageOrderRequestData, err := validation.ParseRequestToUpdateProductImageOrderRequestData(r)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	product, err := models.GetProductByID(ctx, updateProductImageOrderRequestData.ProductID)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	if product.ClientID != user.Client.ID {
		util.ErrResponse(w, models.ErrUserDoesNotBelongToClient, http.StatusForbidden)
		return
	}

	updateProductImageOrderRequestData, err := validation.ParseRequestToUpdateProductImageOrderRequestData(r)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	err = product.GetProductImages(ctx)
	if err != nil {
		util.ErrResponse(w, err, http.StatusBadRequest)
		return
	}

	// make sure all id's are unique
	uniqueOrder := []int{}
	for _, order := range updateProductImageOrderRequestData.Order {
		for _, uniqueOrderItem := range uniqueOrder {
			if order == uniqueOrderItem {
				util.ErrResponse(w, ErrDuplicateProductImageIDs, http.StatusBadRequest)
				return
			}
		}

		uniqueOrder = append(uniqueOrder, order)
	}

	if len(updateProductImageOrderRequestData.Order) != len(product.ProductImages) {
		util.ErrResponse(w, ErrInvalidNumberOfProductImages, http.StatusBadRequest)
		return
	}

	for i, productImageID := range updateProductImageOrderRequestData.Order {

		productImage, err := models.GetProductImageByID(ctx, productImageID)
		if err != nil {
			util.ErrResponse(w, err, http.StatusBadRequest)
			return
		}

		if productImage.ProductID != product.ID {
			util.ErrResponse(w, ErrProductImageDoesNotBelongToProduct, http.StatusBadRequest)
			return
		}

		_, err = models.UpdateProductImage(ctx, models.UpdateProductImageInput{
			ID:    productImageID,
			Order: i,
		})
		if err != nil {
			util.ErrResponse(w, err, http.StatusBadRequest)
			return
		}

	}

	util.SuccessResponse(w, http.StatusOK)

}
