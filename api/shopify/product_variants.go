package shopify

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/shipply-io/shipply-io-backend/api/shopify/gen"
	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/util"
)

func RetrieveProductVariants(shopName string, accessToken string) ([]gen.GetProductVariants_ProductVariants_Edges_Node, error) {

	client := NewClient(shopName, accessToken)

	var cursor *string
	productVariants := make([]gen.GetProductVariants_ProductVariants_Edges_Node, 0)

	for {
		getProductVariantsResponse, err := client.GetProductVariants(context.Background(), 100, cursor)
		if err != nil {
			return nil, err
		}

		for _, edge := range getProductVariantsResponse.ProductVariants.Edges {
			productVariants = append(productVariants, edge.Node)
		}

		if !getProductVariantsResponse.ProductVariants.PageInfo.HasNextPage {
			break
		}

		cursor = getProductVariantsResponse.ProductVariants.PageInfo.EndCursor
	}

	return productVariants, nil
}

func MapProductVariantToDatabaseProduct(productVariant gen.GetProductVariants_ProductVariants_Edges_Node) models.Product {
	product := models.Product{}

	if productVariant.Title != "Default Title" && productVariant.Title != "" {
		product.Name = productVariant.Title
	} else {
		product.Name = productVariant.Product.Title
	}

	if productVariant.Image != nil {
		product.ImageURL = productVariant.Image.URL
	}

	product.Sku = *productVariant.Sku
	product.Weight = *productVariant.Weight
	product.WeightUnit = productVariant.WeightUnit.String()
	product.Grams = util.ConvertWeightToGrams(product.Weight, product.WeightUnit)

	if productVariant.Barcode != nil {
		product.Barcode = *productVariant.Barcode
	}

	if productVariant.Product.PriceRange.MaxVariantPrice.CurrencyCode.String() != "" {
		product.PriceCurrency = productVariant.Product.PriceRange.MaxVariantPrice.CurrencyCode.String()
	}

	if productVariant.InventoryItem.UnitCost != nil {
		value, err := strconv.ParseFloat(productVariant.InventoryItem.UnitCost.Amount, 64)
		if err == nil {
			product.Value = value
		}
		product.ValueCurrency = productVariant.InventoryItem.UnitCost.CurrencyCode.String()
	}

	if productVariant.InventoryItem.CountryCodeOfOrigin != nil {
		product.CountryOfManufacture = (productVariant.InventoryItem.CountryCodeOfOrigin.String())
	}

	if productVariant.InventoryItem.HarmonizedSystemCode != nil {
		product.TariffCode = *productVariant.InventoryItem.HarmonizedSystemCode
	}

	product.CustomsValue = product.Price

	price, err := strconv.Atoi(productVariant.Price)
	if err == nil {
		product.Price = float64(price)
	}

	return product
}

func MapProductVariantToDatabaseShopifyProduct(productVariant gen.GetProductVariants_ProductVariants_Edges_Node) models.ShopifyProduct {
	shopifyProduct := models.ShopifyProduct{}
	shopifyProduct.VariantShopifyGraphqlID = productVariant.ID
	shopifyProduct.ProductShopifyGraphqlID = productVariant.Product.ID
	shopifyProduct.InventoryPolicy = productVariant.InventoryPolicy.String()
	shopifyProduct.ProductStatus = productVariant.Product.Status.String()
	shopifyProduct.Sku = *productVariant.Sku
	shopifyProduct.ProductHandle = productVariant.Product.Handle
	shopifyProduct.ProductTitle = productVariant.Product.Title
	shopifyProduct.ProductType = productVariant.Product.ProductType
	shopifyProduct.ProductVendor = productVariant.Product.Vendor
	shopifyProduct.VariantTitle = productVariant.Title
	shopifyProduct.InventoryItemShopifyGraphqlID = productVariant.InventoryItem.ID

	return shopifyProduct
}

type ProductCreateWebhookRequest struct {
	ID                int    `json:"id"`
	AdminGraphqlAPIID string `json:"admin_graphql_api_id"`
	Title             string `json:"title"`
	Handle            string `json:"handle"`
	ProductType       string `json:"product_type"`
	Vendor            string `json:"vendor"`
	Status            string `json:"status"`
	Tags              string `json:"tags"`
	Variants          []struct {
		ID                            int     `json:"id"`
		AdminGraphqlAPIID             string  `json:"admin_graphql_api_id"`
		Title                         string  `json:"title"`
		Sku                           string  `json:"sku"`
		Price                         string  `json:"price"`
		Weight                        float64 `json:"weight"`
		WeightUnit                    string  `json:"weight_unit"`
		InventoryItemID               int     `json:"inventory_item_id"`
		InventoryItemShopifyGraphqlID string  `json:"inventory_item_shopify_graphql_id"`
		Barcode                       string  `json:"barcode"`
		Grams                         int     `json:"grams"`
		InventoryPolicy               string  `json:"inventory_policy"`
	} `json:"variants"`
	Image struct {
		Src string `json:"src"`
	} `json:"image"`
	Images []struct {
		Src        string `json:"src"`
		VariantIDs []int  `json:"variant_ids"`
	} `json:"images"`
}

func ParseAndValidateProductCreateRequest(r *http.Request) (*string, []string) {

	var productID string

	var errs []string

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errs = append(errs, "failed to read body")
		return nil, errs
	}

	aux := &struct {
		AdminGraphqlAPIID json.RawMessage `json:"admin_graphql_api_id"`
	}{}

	err = json.Unmarshal(body, aux)
	if err != nil {
		return nil, []string{"invalid json"}
	}

	if aux.AdminGraphqlAPIID != nil {
		if err := json.Unmarshal(aux.AdminGraphqlAPIID, &productID); err != nil {
			errs = append(errs, "invalid product id")
		}
	}

	if len(errs) > 0 {
		return nil, errs
	}

	return &productID, nil
}

func ParseAndValidateProductUpdateRequest(r *http.Request) (*string, []string) {
	var productID string

	var errs []string

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errs = append(errs, "failed to read body")
		return nil, errs
	}

	aux := &struct {
		AdminGraphqlAPIID json.RawMessage `json:"admin_graphql_api_id"`
	}{}

	err = json.Unmarshal(body, aux)
	if err != nil {
		return nil, []string{"invalid json"}
	}

	if aux.AdminGraphqlAPIID != nil {
		if err := json.Unmarshal(aux.AdminGraphqlAPIID, &productID); err != nil {
			errs = append(errs, "invalid product id")
		}
	}

	if len(errs) > 0 {
		return nil, errs
	}

	return &productID, nil
}

func RetrieveProductVariantsByProductID(shopName string, accessToken string, productID string) ([]gen.GetProductVariantsByProductID_Product_Variants_Edges_Node, error) {
	client := NewClient(shopName, accessToken)

	var cursor *string
	productVariants := make([]gen.GetProductVariantsByProductID_Product_Variants_Edges_Node, 0)

	for {
		retries := 0
		var getProductVariantsResponse *gen.GetProductVariantsByProductID
		var err error

		for {
			getProductVariantsResponse, err = client.GetProductVariantsByProductID(context.Background(), productID, 10, cursor)
			if err == nil {
				break
			}

			errorType, err := DetermineShopifyGraphqlError(err)
			if err != nil {
				return nil, err
			}

			if errorType == ShopifyRateLimitError {
				if retries >= MaxRetries {
					return nil, errors.New("maximum rate limit retries exceeded")
				}

				waitDuration := ExponentialBackoffWithJitter(retries, InitialWaitDuration, JitterFactor)
				time.Sleep(waitDuration)
				retries++
			} else if errorType == ShopifyQueryLimitExceededError {
				return nil, errors.New("query limit exceeded")
			} else {
				return nil, err
			}
		}

		for _, edge := range getProductVariantsResponse.Product.Variants.Edges {
			productVariants = append(productVariants, edge.Node)
		}

		if !getProductVariantsResponse.Product.Variants.PageInfo.HasNextPage {
			break
		}
		cursor = getProductVariantsResponse.Product.Variants.PageInfo.EndCursor

	}

	return productVariants, nil
}

type ShopifyGraphqlModelProductVariant struct {
	ID                string
	Image             string
	InventoryItem     ShopifyGraphqlModelInventoryItem
	InventoryPolicy   string
	InventoryQuantity *int
	Price             float64
	Product           ShopifyGraphqlModelProduct
	Barcode           *string
	Sku               *string
	TaxCode           *string
	Taxable           bool
	Title             string
	UpdatedAt         time.Time
	Weight            *float64
	WeightUnit        string
	RequiresShipping  bool
}

type ShopifyGraphqlModelProduct struct {
	Status         string
	Description    string
	Handle         string
	ID             string
	CurrencyCode   string
	ProductType    string
	PublishedAt    *time.Time
	Tags           []string
	Title          string
	UpdatedAt      string
	Vendor         string
	OnlineStoreURL *string
}

type ShopifyGraphqlModelInventoryItem struct {
	ID                   string
	CountryCodeOfOrigin  *string
	DuplicateSkuCount    int
	HarmonizedSystemCode *string
	InventoryHistoryURL  *string
	LocationsCount       int
	ProvinceCodeOfOrigin *string
	RequiresShipping     bool
	Sku                  *string
	Tracked              bool
	UnitCost             float64
}

func ConvertGetProductById_ProductToShopifyModelProduct(shopifyProduct gen.GetProductVariantsByProductID_Product_Variants_Edges_Node) (*ShopifyGraphqlModelProductVariant, error) {

	var err error
	var timeValue time.Time

	mProduct := ShopifyGraphqlModelProductVariant{}
	mProduct.ID = shopifyProduct.ID

	if shopifyProduct.Image != nil {
		mProduct.Image = shopifyProduct.Image.URL
	} else {
		for _, productImage := range shopifyProduct.Product.Images.Edges {
			mProduct.Image = productImage.Node.Src
		}
	}

	mProduct.InventoryItem.ID = shopifyProduct.InventoryItem.ID
	if shopifyProduct.InventoryItem.CountryCodeOfOrigin != nil {
		mProduct.InventoryItem.CountryCodeOfOrigin = (*string)(shopifyProduct.InventoryItem.CountryCodeOfOrigin)
		mProduct.InventoryItem.HarmonizedSystemCode = (*string)(shopifyProduct.InventoryItem.HarmonizedSystemCode)
		mProduct.InventoryItem.UnitCost, err = strconv.ParseFloat(shopifyProduct.InventoryItem.GetUnitCost().Amount, 64)
		if err != nil {
			mProduct.InventoryItem.UnitCost = 0
		}
	}

	mProduct.InventoryPolicy = (string)(shopifyProduct.InventoryPolicy)
	mProduct.InventoryQuantity = shopifyProduct.InventoryQuantity
	mProduct.Price, err = strconv.ParseFloat(shopifyProduct.GetPrice(), 64)
	if err != nil {
		mProduct.Price = 0
	}

	mProduct.Product.Status = shopifyProduct.Product.Status.String()
	mProduct.Product.Description = shopifyProduct.Product.Description
	mProduct.Product.Handle = shopifyProduct.Product.Handle
	mProduct.Product.ID = shopifyProduct.Product.ID
	mProduct.Product.CurrencyCode = shopifyProduct.Product.PriceRange.MaxVariantPrice.CurrencyCode.String()
	mProduct.Product.ProductType = shopifyProduct.Product.ProductType
	timeValue, err = util.ParseRFC3339Date(*shopifyProduct.Product.PublishedAt)
	if err != nil {
		mProduct.Product.PublishedAt = nil
	}
	mProduct.Product.PublishedAt = &timeValue
	mProduct.Product.Tags = shopifyProduct.Product.Tags
	mProduct.Product.Title = shopifyProduct.Product.Title
	timeValue, err = util.ParseRFC3339Date(shopifyProduct.Product.UpdatedAt)
	if err != nil {
		return nil, err
	}
	mProduct.Product.UpdatedAt = shopifyProduct.Product.UpdatedAt
	mProduct.Product.Vendor = shopifyProduct.Product.Vendor
	if shopifyProduct.Product.OnlineStoreURL != nil {
		mProduct.Product.OnlineStoreURL = (*string)(shopifyProduct.Product.OnlineStoreURL)
	}

	if shopifyProduct.Barcode != nil {
		mProduct.Barcode = (*string)(shopifyProduct.Barcode)
	}

	if shopifyProduct.Sku != nil {
		mProduct.Sku = (*string)(shopifyProduct.Sku)
	}

	if shopifyProduct.TaxCode != nil {
		mProduct.TaxCode = (*string)(shopifyProduct.TaxCode)
	}

	mProduct.Taxable = shopifyProduct.Taxable
	mProduct.Title = shopifyProduct.Title
	timeValue, err = util.ParseRFC3339Date(shopifyProduct.UpdatedAt)
	if err != nil {
		return nil, err
	}
	mProduct.UpdatedAt = timeValue

	if shopifyProduct.Weight != nil {
		mProduct.Weight = (*float64)(shopifyProduct.Weight)
	}

	mProduct.WeightUnit = shopifyProduct.WeightUnit.String()
	mProduct.RequiresShipping = shopifyProduct.RequiresShipping

	return &mProduct, nil
}

func ConvertGetProductVariants_ProductToShopifyModelProduct(shopifyProduct gen.GetProductVariants_ProductVariants_Edges_Node) (*ShopifyGraphqlModelProductVariant, error) {

	var err error
	var timeValue time.Time

	mProduct := ShopifyGraphqlModelProductVariant{}
	mProduct.ID = shopifyProduct.ID

	if shopifyProduct.Image != nil {
		mProduct.Image = shopifyProduct.Image.URL
	} else {
		for _, productImage := range shopifyProduct.Product.Images.Edges {
			mProduct.Image = productImage.Node.Src
		}
	}

	mProduct.InventoryItem.ID = shopifyProduct.InventoryItem.ID
	if shopifyProduct.InventoryItem.CountryCodeOfOrigin != nil {
		mProduct.InventoryItem.CountryCodeOfOrigin = (*string)(shopifyProduct.InventoryItem.CountryCodeOfOrigin)
		mProduct.InventoryItem.HarmonizedSystemCode = (*string)(shopifyProduct.InventoryItem.HarmonizedSystemCode)
		mProduct.InventoryItem.UnitCost, err = strconv.ParseFloat(shopifyProduct.InventoryItem.GetUnitCost().Amount, 64)
		if err != nil {
			mProduct.InventoryItem.UnitCost = 0
		}
	}

	mProduct.InventoryPolicy = (string)(shopifyProduct.InventoryPolicy)
	mProduct.InventoryQuantity = shopifyProduct.InventoryQuantity
	mProduct.Price, err = strconv.ParseFloat(shopifyProduct.GetPrice(), 64)
	if err != nil {
		mProduct.Price = 0
	}

	mProduct.Product.Status = shopifyProduct.Product.Status.String()
	mProduct.Product.Description = shopifyProduct.Product.Description
	mProduct.Product.Handle = shopifyProduct.Product.Handle
	mProduct.Product.ID = shopifyProduct.Product.ID
	mProduct.Product.CurrencyCode = shopifyProduct.Product.PriceRange.MaxVariantPrice.CurrencyCode.String()
	mProduct.Product.ProductType = shopifyProduct.Product.ProductType
	if shopifyProduct.Product.PublishedAt != nil {
		timeValue, err = util.ParseRFC3339Date(*shopifyProduct.Product.PublishedAt)
		if err != nil {
			mProduct.Product.PublishedAt = nil
		}
		mProduct.Product.PublishedAt = &timeValue
	} else {
		mProduct.Product.PublishedAt = nil
	}
	mProduct.Product.Tags = shopifyProduct.Product.Tags
	mProduct.Product.Title = shopifyProduct.Product.Title
	timeValue, err = util.ParseRFC3339Date(shopifyProduct.Product.UpdatedAt)
	if err != nil {
		return nil, err
	}
	mProduct.Product.UpdatedAt = shopifyProduct.Product.UpdatedAt
	mProduct.Product.Vendor = shopifyProduct.Product.Vendor
	if shopifyProduct.Product.OnlineStoreURL != nil {
		mProduct.Product.OnlineStoreURL = (*string)(shopifyProduct.Product.OnlineStoreURL)
	}

	if shopifyProduct.Barcode != nil {
		mProduct.Barcode = (*string)(shopifyProduct.Barcode)
	}

	if shopifyProduct.Sku != nil {
		mProduct.Sku = (*string)(shopifyProduct.Sku)
	}

	if shopifyProduct.TaxCode != nil {
		mProduct.TaxCode = (*string)(shopifyProduct.TaxCode)
	}

	mProduct.Taxable = shopifyProduct.Taxable
	mProduct.Title = shopifyProduct.Title
	timeValue, err = util.ParseRFC3339Date(shopifyProduct.UpdatedAt)
	if err != nil {
		return nil, err
	}
	mProduct.UpdatedAt = timeValue

	if shopifyProduct.Weight != nil {
		mProduct.Weight = (*float64)(shopifyProduct.Weight)
	}

	mProduct.WeightUnit = shopifyProduct.WeightUnit.String()
	mProduct.RequiresShipping = shopifyProduct.RequiresShipping

	return &mProduct, nil
}
