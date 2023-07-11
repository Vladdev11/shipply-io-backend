package OrganizationHandlers

import (
	"net/http"

	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/responses"
	"github.com/shipply-io/shipply-io-backend/util"
	"gorm.io/gorm"
)

func SearchProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusBadRequest)
		return
	}

	err = user.GetOrganization(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusBadRequest)
		return
	}

	productSearch := models.ProductSearchRequest{}
	err = productSearch.ParseAndValidateRequest(r)
	if err != nil {
		util.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if productSearch.ClientID != 0 {
		if !user.Organization.IsClientOwner(ctx, productSearch.ClientID) {
			util.ErrorResponse(w, "client does not belong to organization", http.StatusForbidden)
			return
		}
	}

	productSearch.OrganizationID = user.OwnerID
	products, err := user.Organization.SearchProducts(ctx, productSearch)
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

	err = user.GetOrganization(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusUnauthorized)
		return
	}

	// TODO - move this to validation package
	request := models.ProductListRequest{}
	request.OrganizationID = user.Organization.ID
	errors := request.ParseAndValidateRequest(r)
	if errors != nil {
		util.ErrorsResponse(w, errors, http.StatusBadRequest)
		return
	}

	products, count, total, err := user.Organization.GetProducts(ctx, request)
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

	err = user.GetOrganization(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusUnauthorized)
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

	if !user.Organization.IsClientOwner(ctx, product.ClientID) {
		util.ErrorResponse(w, "user does not have access to this product", http.StatusForbidden)
		return
	}

	err = product.GetProductBundles(ctx)
	if err != nil && err != gorm.ErrRecordNotFound {
		util.ErrorResponse(w, "failed to get product bundles", http.StatusBadRequest)
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

	err = user.GetOrganization(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusUnauthorized)
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

	if !user.Organization.IsClientOwner(ctx, product.ClientID) {
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

	err = user.GetOrganization(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusUnauthorized)
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

	if !user.Organization.IsClientOwner(ctx, product.ClientID) {
		util.ErrorResponse(w, "user does not have access to this product", http.StatusForbidden)
		return
	}

	err = product.GetInventoryLocations(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get inventory locations", http.StatusBadRequest)
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

	response := responses.GenerateGetProductInventoryResponse(product)
	util.JSONResponse(w, response, http.StatusOK)

}

func GetProductBundles(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	user, err := models.GetRequestingUser(r)
	if err != nil {
		util.ErrorResponse(w, "failed to get user", http.StatusUnauthorized)
		return
	}

	err = user.GetOrganization(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusUnauthorized)
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

	if !user.Organization.IsClientOwner(ctx, product.ClientID) {
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
		util.ErrorResponse(w, "failed to get bundles", http.StatusBadRequest)
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

	err = user.GetOrganization(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusUnauthorized)
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

	if !user.Organization.IsClientOwner(ctx, product.ClientID) {
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

	err = user.GetOrganization(ctx)
	if err != nil {
		util.ErrorResponse(w, "failed to get organization", http.StatusUnauthorized)
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

	if !user.Organization.IsClientOwner(ctx, product.ClientID) {
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
