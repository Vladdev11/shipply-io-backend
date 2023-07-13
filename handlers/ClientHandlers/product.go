package ClientHandlers

import (
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

func SearchProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusBadRequest)
		return
	}

	// TODO - move this to validation package
	productSearch := models.ProductSearchRequest{}
	err = productSearch.ParseAndValidateRequest(r)
	if err != nil {
		util.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	productSearch.ClientID = user.OwnerID

	products, err := user.Client.SearchProducts(ctx, productSearch)
	if err != nil {
		util.ErrorResponse(w, "failed to search products", http.StatusBadRequest)
		return
	}

	response := responses.GenerateSearchProductsResponse(products)
	util.JSONResponse(w, response, http.StatusOK)

}

func ListProducts(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetClient(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get client", http.StatusUnauthorized)
		return
	}

	request := models.ProductListRequest{}
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	if user.Client.ID != request.ClientID {
		util.ErrorResponse(w, "user does not have access to this client", http.StatusForbidden)
		return
	}

	products, count, total, err := user.Client.GetProducts(ctx, request)
	if err != nil {
		util.ErrorResponse(w, "failed to get products", http.StatusBadRequest)
		return
	}

	response := responses.GenerateListProductsResponse(products, count, total)
	util.JSONResponse(w, response, http.StatusOK)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetClient(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get client", http.StatusUnauthorized)
		return
	}

	productID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "invalid product id", http.StatusBadRequest)
		return
	}

	product, err := models.GetProductByID(ctx, productID)
	if err != nil {
		util.ErrorResponse(w, "failed to find product", http.StatusBadRequest)
		return
	}

	if user.Client.ID != product.ClientID {
		util.ErrorResponse(w, "user does not have access to this product", http.StatusForbidden)
		return
	}

	response := responses.GenerateGetProductResponse(product)
	util.JSONResponse(w, response, http.StatusOK)

}

func GetProductOrders(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetClient(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get client", http.StatusUnauthorized)
		return
	}

	productID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "invalid product id", http.StatusBadRequest)
		return
	}

	product, err := models.GetProductByID(ctx, productID)
	if err != nil {
		util.ErrorResponse(w, "failed to find product", http.StatusBadRequest)
		return
	}

	if user.Client.ID != product.ClientID {
		util.ErrorResponse(w, "user does not have access to this product", http.StatusForbidden)
		return
	}

	orders, err := models.GetOrdersByProductID(ctx, productID)
	if err != nil {
		util.ErrorResponse(w, "failed to get orders", http.StatusBadRequest)
		return
	}

	response := responses.GenerateGetProductOrdersResponse(orders)
	util.JSONResponse(w, response, http.StatusOK)

}

func GetProductInventory(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetClient(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get client", http.StatusUnauthorized)
		return
	}

	productID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "invalid product id", http.StatusBadRequest)
		return
	}

	product, err := models.GetProductByID(ctx, productID)
	if err != nil {
		util.ErrorResponse(w, "failed to find product", http.StatusBadRequest)
		return
	}

	if user.Client.ID != product.ClientID {
		util.ErrorResponse(w, "user does not have access to this product", http.StatusForbidden)
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

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetClient(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get client", http.StatusUnauthorized)
		return
	}

	productID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "invalid product id", http.StatusBadRequest)
		return
	}

	product, err := models.GetProductByID(ctx, productID)
	if err != nil {
		util.ErrorResponse(w, "failed to find product", http.StatusBadRequest)
		return
	}

	if user.Client.ID != product.ClientID {
		util.ErrorResponse(w, "user does not have access to this product", http.StatusForbidden)
		return
	}

	productIsBundle, err := product.IsBundle(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to check if product is bundle", http.StatusBadRequest)
		return
	}

	if productIsBundle {
		util.ErrorResponse(w, "product is a bundle", http.StatusBadRequest)
		return
	}

	err = product.GetProductBundles(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get product bundles", http.StatusBadRequest)
		return
	}

	for i := range product.ProductBundles {
		err = product.ProductBundles[i].GetProduct(ctx)
		if err != nil {
			util.ErrorResponse(w, "failed to get bundle product", http.StatusBadRequest)
			return
		}
	}

	response := responses.GenerateGetProductBundlesResponse(product)
	util.JSONResponse(w, response, http.StatusOK)

}

func GetProductBundleComponents(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetClient(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get client", http.StatusUnauthorized)
		return
	}

	productID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "invalid bundle id", http.StatusBadRequest)
		return
	}

	product, err := models.GetProductByID(ctx, productID)
	if err != nil {
		util.ErrorResponse(w, "failed to find product", http.StatusBadRequest)
		return
	}

	if user.Client.ID != product.ClientID {
		util.ErrorResponse(w, "user does not have access to this product", http.StatusForbidden)
		return
	}

	if product.IsComponent(ctx) {
		util.ErrorResponse(w, "product is a component", http.StatusBadRequest)
		return
	}

	bundle, err := models.GetBundleByProductID(ctx, productID)
	if err != nil {

		if err == gorm.ErrRecordNotFound {
			response := responses.GenerateGetProductBundleComponentsResponse([]models.Product{})
			util.JSONResponse(w, response, http.StatusOK)
			return
		}

		util.ErrorResponse(w, "failed to find bundle", http.StatusBadRequest)
		return
	}

	err = bundle.GetComponents(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get bundle components", http.StatusBadRequest)
		return
	}

	response := responses.GenerateGetProductBundleComponentsResponse(bundle.Components)
	util.JSONResponse(w, response, http.StatusOK)

}

func GetProductStores(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetClient(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get client", http.StatusUnauthorized)
		return
	}

	productID, err := util.GetIntFromPath(r, "id")
	if err != nil {
		util.ErrorResponse(w, "invalid product id", http.StatusBadRequest)
		return
	}

	product, err := models.GetProductByID(ctx, productID)
	if err != nil {
		util.ErrorResponse(w, "failed to find product", http.StatusBadRequest)
		return
	}

	if user.Client.ID != product.ClientID {
		util.ErrorResponse(w, "user does not have access to this product", http.StatusForbidden)
		return
	}

	err = product.GetStores(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get stores", http.StatusBadRequest)
		return
	}

	response := responses.GenerateGetProductStoresResponse(product)
	util.JSONResponse(w, response, http.StatusOK)

}

func CreateProduct(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetClient(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get client", http.StatusUnauthorized)
		return
	}

	err = user.Client.GetOrganization(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusUnauthorized)
		return
	}

	createProductInput, errors := validation.ParseRequestToCreateProductInput(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	if !user.Client.Organization.IsWarehouseOwner(ctx, createProductInput.WarehouseID) {
		util.ErrorResponse(w, "user does not have access to this warehouse", http.StatusForbidden)
		return
	}

	existingProductByClientIDAndSku, err := models.GetProductByClientIDAndSku(ctx, user.Client.ID, createProductInput.Sku)
	if err != nil && err != gorm.ErrRecordNotFound {
		util.ErrorResponse(w, "failed to check if product sku is unique", http.StatusBadRequest)
		return
	}

	if existingProductByClientIDAndSku != nil {
		util.ErrorResponse(w, "product sku already exists", http.StatusBadRequest)
		return
	}

	existingProductByBarcodeAndClientID, err := models.GetProductByBarcodeAndClientID(ctx, createProductInput.Barcode, user.Client.ID)
	if err != nil && err != gorm.ErrRecordNotFound {
		util.ErrorResponse(w, "failed to check if product barcode is unique", http.StatusBadRequest)
		return
	}

	if existingProductByBarcodeAndClientID != nil {
		util.ErrorResponse(w, "product barcode already exists", http.StatusBadRequest)
		return
	}

	// TODO add warehouse id?
	// TODO add description?
	product := models.Product{
		ClientID:            createProductInput.ClientID,
		Sku:                 createProductInput.Sku,
		Barcode:             createProductInput.Barcode,
		Name:                createProductInput.Name,
		Value:               createProductInput.Value,
		AdditionalImageUrls: models.StringSlice{},
	}

	if createProductInput.Weight != nil {
		product.Weight = createProductInput.Weight.Value
		product.WeightUnit = createProductInput.Weight.Unit
	}

	if createProductInput.Dimensions != nil {
		product.Length = createProductInput.Dimensions.Length
		product.Width = createProductInput.Dimensions.Width
		product.Height = createProductInput.Dimensions.Height
	}

	for i, image := range createProductInput.Images {

		fileUUID := uuid.New().String()
		fileExtension := strings.Split(image.FileType, "/")[1]

		err = api.S3FromContext(ctx).UploadFileToCDN(image.ImageData, fileUUID, fileExtension, image.FileType)
		if err != nil {
			util.ErrorResponse(w, "failed to upload image to s3", http.StatusInternalServerError)
			return
		}

		imageURL := fmt.Sprintf("%s/%s.%s", util.CDNFromContext(ctx), fileUUID, fileExtension)
		if i == 0 {
			product.ImageURL = imageURL
		} else {
			product.AdditionalImageUrls = append(product.AdditionalImageUrls, imageURL)
		}

	}

	err = product.Create(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to create product", http.StatusBadRequest)
		return
	}

	response := responses.GenerateCreateProductResponse(product)
	util.JSONResponse(w, response, http.StatusOK)

}
