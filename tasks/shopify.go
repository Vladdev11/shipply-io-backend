package tasks

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/segmentio/encoding/json"
	"github.com/shipply-io/shipply-io-backend/api/shopify"
	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/util"

	shipengineHandlers "github.com/shipply-io/shipply-io-backend/api/shipengine/handlers"
	ShipengineModels "github.com/shipply-io/shipply-io-backend/api/shipengine/models"
)

func PullShopifyOrdersTask(ctx context.Context, storeID int, shopName string, accessToken string) error {

	//get orders from shopify
	orders, err := shopify.RetrieveOrders(shopName, accessToken)
	if err != nil {
		fmt.Println(err)
		return err
	}

	//loop through orders and process them
	for _, order := range orders {
		gqlOrder, err := shopify.ConvertGetOrders_Orders_Edges_NodeToShopifyModelOrder(order)
		if err != nil {
			models.CreateSystemError(ctx, fmt.Sprintf("Error converting shopify order to graphql order: %s", err.Error()))
			continue
		}
		err = SyncShopifyGraphqlOrder(ctx, storeID, shopName, accessToken, gqlOrder)
		if err != nil {
			models.CreateSystemError(ctx, fmt.Sprintf("Error syncing shopify order: %s", err.Error()))
			continue
		}
	}

	return nil
}

func PullShopifyProductsTask(ctx context.Context, storeID int, shopName string, accessToken string) error {

	store, err := models.GetStoreByID(ctx, storeID)
	if err != nil {
		return err
	}

	locations, err := models.GetShopifyLocationsByStoreID(ctx, storeID)
	if err != nil {
		return err
	}

	productVariants, err := shopify.RetrieveProductVariants(shopName, accessToken)
	if err != nil {
		return err
	}

	for _, productVariant := range productVariants {
		//Convert the graphql product to a internal shopify graphql model product
		gqlProduct, err := shopify.ConvertGetProductVariants_ProductToShopifyModelProduct(productVariant)
		if err != nil {
			models.CreateSystemError(ctx, "failed to convert shopify product to model product: "+err.Error())
			return err
		}

		// Activate Inventory Item for each location
		for _, location := range locations {
			err = shopify.ActivateInventoryItem(ctx, shopName, accessToken, gqlProduct.InventoryItem.ID, location.ShopifyLocationID)
			if err != nil {
				return err
			}
		}

		// Sync the product to the database
		err = SyncShopifyGraphqlProduct(ctx, store.ID, store.ShopifyShopName(ctx), store.ShopifyAccessToken(ctx), *gqlProduct)
		if err != nil {
			models.CreateSystemError(ctx, "failed to convert shopify product to model product: "+err.Error())
			continue
		}

	}

	return nil

}

func SyncShopifyGraphqlOrder(ctx context.Context, storeID int, shopName string, accessToken string, order shopify.ShopifyGraphqlModelOrder) error {

	//get client by store id
	client, err := models.GetClientByStoreID(ctx, storeID)
	if err != nil {
		return err
	}

	//get store by store id
	store, err := models.GetStoreByID(ctx, storeID)
	if err != nil {
		return err
	}

	//parse shopify settings
	settings, err := store.GetShopifySettings(ctx)
	if err != nil {
		return err
	}

	//get order items from shopify
	orderItems, err := shopify.RetrieveOrderItemsByOrderID(shopName, accessToken, order.ID)
	if err != nil {
		return err
	}

	//get order from database
	existingOrder, err := models.GetOrderByStoreAndAPIID(ctx, storeID, order.ID)

	//handle existing order
	if err == nil && existingOrder != nil {
		dbOrder := existingOrder

		//handle shipping address changes
		if settings.UpdateShippingAddressFromShopify {
			if order.ShippingAddress != nil {
				newOrderModelAddress := shopify.ConvertShopifyAddressToModelAddress(order.ShippingAddress)

				err = newOrderModelAddress.Create(ctx)
				if err != nil {
					return err
				}

				dbOrder.ShipToAddressID = &newOrderModelAddress.ID
			}
		}

		//handle financial status changes
		//financial status
		if order.DisplayFinacialStatus != nil {
			dbOrder.MarketplaceFinacialStatus = *order.DisplayFinacialStatus
			if dbOrder.HasHold(util.PaymentHold) && *order.DisplayFinacialStatus == "paid" {
				dbOrder.RemoveHold(util.PaymentHold)
			}
		}

		//handle fulfillment status changes
		dbOrder.MarketplaceFulfillmentStatus = order.DisplayFulfillmentStatus

		for _, orderItem := range orderItems {

			dbOrderItem, err := models.GetOrderItemByOrderIDAndAPIID(ctx, dbOrder.ID, orderItem.ID)
			if err != nil {
				continue
			}

			dbOrderItem.QuantityShipped = orderItem.Quantity - orderItem.UnfulfilledQuantity
			dbOrderItem.Quantity = orderItem.Quantity
			dbOrderItem.MarketplaceFulfillmentStatus = orderItem.FulfillmentStatus

			if dbOrderItem.QuantityShipped == dbOrderItem.Quantity {
				dbOrderItem.Fulfilled = true
			}

			err = dbOrderItem.Update(ctx)
			if err != nil {
				continue
			}
		}

		//update order
		err = dbOrder.Update(ctx)
		if err != nil {
			models.CreateSystemError(ctx, fmt.Sprintf("Error updating order: %s", err.Error()))
			return err
		}

		return nil
	}

	//if it doesn't, create it
	if err != nil {

		dbOrder := models.Order{}
		dbOrder.APIID = order.ID
		dbOrder.StoreID = storeID
		dbOrder.OrderNumber = order.Name

		if order.Note != nil {
			switch settings.CustomerNotesFieldMapping {
			case "gift_note":
				dbOrder.GiftNote = *order.Note
			case "packing_note":
				dbOrder.PackingNote = *order.Note
			default:
				dbOrder.PackingNote = *order.Note
			}
		}

		//address logic
		if order.ShippingAddress != nil {

			modelAddress := shopify.ConvertShopifyAddressToModelAddress(order.ShippingAddress)

			//validate address via shipengine
			shipengineAddress, _ := shipengineHandlers.ValidateAddress(ctx, ShipengineModels.ConvertAddressToShipengineAddress(modelAddress))

			if shipengineAddress.Status == "verified" {
				//update address with shipengine address
				updatedModelAddress := ShipengineModels.ConvertShipengineAddressToModelAddress(shipengineAddress.MatchedAddress)
				updatedModelAddress.FirstName = modelAddress.FirstName
				updatedModelAddress.LastName = modelAddress.LastName
				updatedModelAddress.AddressVerified = true
				modelAddress = &updatedModelAddress
			} else {
				dbOrder.AddHold(util.AddressHold)
			}
			err = modelAddress.Create(ctx)
			if err != nil {
				models.CreateSystemError(ctx, fmt.Sprintf("Error creating address for order %s: %s", dbOrder.OrderNumber, err.Error()))
			}
			dbOrder.ShipToAddressID = &modelAddress.ID
		} else {
			dbOrder.AddHold(util.AddressHold)
			dbOrder.ShipToAddressID = nil
		}

		if order.BillingAddress != nil {
			modelAddress := shopify.ConvertShopifyAddressToModelAddress(order.ShippingAddress)
			err = modelAddress.Create(ctx)
			if err != nil {
				models.CreateSystemError(ctx, fmt.Sprintf("Error creating address for order %s: %s", dbOrder.OrderNumber, err.Error()))
			}
			dbOrder.BillToAddressID = &modelAddress.ID
		}

		//handle shipping method
		if order.ShippingLine != nil {
			dbOrder.ShippingMethodID, err = shopify.GetOrCreateShippingMethod(ctx, storeID, order.ShippingLine.Title)
			if err != nil {
				models.CreateSystemError(ctx, fmt.Sprintf("Error creating shipping method for order %s: %s", dbOrder.OrderNumber, err.Error()))
			}
		}

		//handle order date
		dbOrder.OrderDate = order.CreatedAt

		//handle order pricing
		dbOrder.Subtotal = order.SubtotalPrice
		dbOrder.Tax = order.TotalTax
		dbOrder.Shipping = order.TotalShipping
		dbOrder.Discount = order.TotalDiscounts
		dbOrder.Tip = order.TotalTip
		dbOrder.Total = math.Round((order.TotalPrice+order.TotalDiscounts)*100) / 100

		//financial status
		if order.DisplayFinacialStatus != nil {
			dbOrder.MarketplaceFinacialStatus = *order.DisplayFinacialStatus
			if !settings.IgnorePaymentStatus {
				if settings.OnlyImportPaidOrders && strings.ToLower(*order.DisplayFinacialStatus) != "paid" {
					return nil
				}
				if strings.ToLower(*order.DisplayFinacialStatus) != "paid" {
					dbOrder.AddHold(util.PaymentHold)
				}
			}
		}

		//fulfillment status
		dbOrder.MarketplaceFulfillmentStatus = order.DisplayFulfillmentStatus

		//email and phone
		if order.Email != nil {
			dbOrder.CustomerEmail = *order.Email
		}
		if order.Phone != nil {
			dbOrder.CustomerPhone = *order.Phone
		}

		//order currency
		dbOrder.Currency = order.Currency

		//TODO SLA SHIP DATE
		dbOrder.RequiredShipDate = util.AddBusinessDays(time.Now(), 1)

		//discount codes
		dbOrder.DiscountCodes, err = json.Marshal(order.DiscountCodes)
		if err != nil {
			models.CreateSystemError(ctx, fmt.Sprintf("Error marshalling discount codes for order %s: %s", dbOrder.OrderNumber, err.Error()))
		}

		//fraud holds
		fraudLevel := strings.ToLower(order.RiskLevel)
		switch settings.ApplyFraudHoldLevel {
		case models.ShopifyFraudLevelLow:
			if fraudLevel == "low" {
				dbOrder.AddHold(util.FraudHold)
			}
		case models.ShopifyFraudLevelMedium:
			if fraudLevel == "medium" {
				dbOrder.AddHold(util.FraudHold)
			}
		case models.ShopifyFraudLevelHigh:
			if fraudLevel == "high" {
				dbOrder.AddHold(util.FraudHold)
			}
		}

		warehouse, err := models.GetWarehousesByOrganizationID(ctx, client.OrganizationID)
		if err != nil {
			models.CreateSystemError(ctx, fmt.Sprintf("Error getting warehouse for order %s: %s", dbOrder.OrderNumber, err.Error()))
		}

		if len(warehouse) != 1 {
			//TODO MULTI-WAREHOUSE ALLOCATION
		} else {
			dbOrder.WarehouseID = &warehouse[0].ID
		}

		err = dbOrder.Create(ctx)
		if err != nil {
			models.CreateSystemError(ctx, fmt.Sprintf("Error creating order %s: %s", dbOrder.OrderNumber, err.Error()))
		}

		//HANDLE ORDER ITEMS
		for _, orderItem := range orderItems {

			dbOrderItem := models.OrderItem{}
			dbOrderItem.OrderID = dbOrder.ID
			dbOrderItem.APIID = orderItem.ID
			dbOrderItem.Name = orderItem.Name
			dbOrderItem.Quantity = orderItem.Quantity
			dbOrderItem.QuantityShipped = orderItem.Quantity - orderItem.UnfulfilledQuantity
			dbOrderItem.MarketplaceFulfillmentStatus = orderItem.FulfillmentStatus

			dbOrderItem.Sku = ""
			if orderItem.Sku != nil && *orderItem.Sku != "" {
				dbOrderItem.Sku = *orderItem.Sku
				product, _ := models.GetProductByClientIDAndSku(ctx, client.ID, *orderItem.Sku)
				if product != nil {
					dbOrderItem.ProductID = &product.ID
				}
			}

			if dbOrderItem.QuantityShipped == dbOrderItem.Quantity {
				dbOrderItem.Fulfilled = true
			}

			//get price
			price, err := strconv.ParseFloat(orderItem.OriginalTotalSet.ShopMoney.Amount, 64)
			if err != nil {
				models.CreateSystemError(ctx, fmt.Sprintf("Error parsing order item price for order %s: %s", dbOrder.OrderNumber, err.Error()))
			}
			dbOrderItem.ItemPrice = price

			//create order item
			err = dbOrderItem.Create(ctx)
			if err != nil {
				models.CreateSystemError(ctx, fmt.Sprintf("Error creating order item for order %s: %s", dbOrder.OrderNumber, err.Error()))
			}

			//TODO run order inventory allocation

		}

	}
	return nil
}

func SyncShopifyGraphqlProduct(ctx context.Context, storeID int, shopName string, accessToken string, productVariant shopify.ShopifyGraphqlModelProductVariant) error {

	product := productVariant

	if product.Sku == nil {
		return errors.New("no product sku")
	}

	//get client by store id
	client, err := models.GetClientByStoreID(ctx, storeID)
	if err != nil {
		return err
	}

	//get store by store id
	store, err := models.GetStoreByID(ctx, storeID)
	if err != nil {
		return err
	}

	//check if store is active
	if !store.Active {
		return nil
	}

	// parse shopify settings
	settings, err := store.GetShopifySettings(ctx)
	if err != nil {
		return err
	}

	//check if product is published
	if settings.OnlyImportActiveProducts && strings.ToLower(productVariant.Product.Status) != "active" {
		return nil
	}

	//get product by API ID
	dbShopifyProduct, err := models.GetShopifyProductAndProductByGraphqlID(ctx, productVariant.ID)

	//handle existing product
	if dbShopifyProduct != nil && err == nil {
		if product.Sku != nil {
			dbShopifyProduct.Sku = *product.Sku
		}

		dbShopifyProduct.InventoryPolicy = product.InventoryPolicy
		dbShopifyProduct.ProductStatus = product.Product.Status
		dbShopifyProduct.ProductHandle = product.Product.Handle
		dbShopifyProduct.ProductType = product.Product.ProductType
		dbShopifyProduct.ProductTitle = product.Product.Title
		dbShopifyProduct.ProductVendor = product.Product.Vendor
		dbShopifyProduct.VariantTitle = product.Title

		//check if product exists by sku
		dbProduct, _ := models.GetProductByClientIDAndSku(ctx, client.ID, dbShopifyProduct.Sku)
		if dbProduct != nil && dbProduct.ID != dbShopifyProduct.ProductID {
			dbShopifyProduct.ProductID = dbProduct.ID
		}

		//if internal product doesn't exist, create it
		if dbProduct == nil {
			//create new product
			dbProduct = &models.Product{}
			dbProduct.Sku = dbShopifyProduct.Sku
			dbProduct.ClientID = client.ID

			if product.Title != "Default Title" && product.Title != "" {
				dbProduct.Name = product.Title
			} else {
				dbProduct.Name = product.Product.Title
			}

			// TODO map image
			if product.Image != "" {
				// dbProduct.ImageURL = product.Image
			}
		}

		if product.Weight != nil && product.WeightUnit != "" {
			dbProduct.Weight = *product.Weight
			dbProduct.WeightUnit = product.WeightUnit
			dbProduct.Grams = util.ConvertWeightToGrams(dbProduct.Weight, dbProduct.WeightUnit)
		}

		if product.Barcode != nil {
			dbProduct.Barcode = *product.Barcode
		} else {
			dbProduct.Barcode = util.GenerateRandomBarcode(12)
		}

		dbProduct.PriceCurrency = product.Product.CurrencyCode
		dbProduct.ValueCurrency = dbProduct.PriceCurrency

		dbProduct.Value = product.InventoryItem.UnitCost

		if product.InventoryItem.CountryCodeOfOrigin != nil {
			dbProduct.CountryOfManufacture = *product.InventoryItem.CountryCodeOfOrigin
		}

		if product.InventoryItem.HarmonizedSystemCode != nil {
			dbProduct.TariffCode = *product.InventoryItem.HarmonizedSystemCode
		}

		dbProduct.Price = product.Price
		dbProduct.CustomsValue = dbProduct.Price

		//update shopify product
		err = dbShopifyProduct.Update(ctx)
		if err != nil {
			return err
		}

		//update product
		err = dbProduct.Update(ctx)
		if err != nil {
			return err
		}

		return nil
	}

	//handle new product
	if err != nil {

		//handle shopify product mappings
		newShopifyProduct := models.ShopifyProduct{}
		newShopifyProduct.ProductShopifyGraphqlID = product.ID
		newShopifyProduct.VariantShopifyGraphqlID = product.ID
		newShopifyProduct.StoreID = storeID

		if product.Sku != nil {
			newShopifyProduct.Sku = *product.Sku
		}

		newShopifyProduct.InventoryPolicy = product.InventoryPolicy
		newShopifyProduct.InventoryItemShopifyGraphqlID = product.InventoryItem.ID
		newShopifyProduct.ProductStatus = product.Product.Status
		newShopifyProduct.ProductTitle = product.Product.Title
		newShopifyProduct.VariantTitle = product.Title
		newShopifyProduct.ProductVendor = product.Product.Vendor

		//handle product mappings
		dbProduct, _ := models.GetProductByClientIDAndSku(ctx, client.ID, newShopifyProduct.Sku)
		//if there is already a product with this sku, update the shopify product to point to it and return
		if dbProduct != nil {
			newShopifyProduct.ProductID = dbProduct.ID
			err = newShopifyProduct.Create(ctx)
			if err != nil {
				models.CreateSystemError(ctx, fmt.Sprintf("Error creating shopify product %s: %s", newShopifyProduct.VariantShopifyGraphqlID, err.Error()))
				return err
			}
			return nil
		}

		//create new product
		newProduct := models.Product{}
		newProduct.Sku = newShopifyProduct.Sku
		newProduct.ClientID = client.ID

		if product.Title != "Default Title" && product.Title != "" {
			newProduct.Name = product.Title
		} else {
			newProduct.Name = product.Product.Title
		}

		// TODO map image
		if product.Image != "" {
			// newProduct.ImageURL = product.Image
		}

		if product.Weight != nil && product.WeightUnit != "" {
			newProduct.Weight = *product.Weight
			newProduct.WeightUnit = product.WeightUnit
			newProduct.Grams = util.ConvertWeightToGrams(newProduct.Weight, newProduct.WeightUnit)
		}

		if product.Barcode != nil {
			newProduct.Barcode = *product.Barcode
		} else {
			newProduct.Barcode = util.GenerateRandomBarcode(12)
		}

		newProduct.PriceCurrency = product.Product.CurrencyCode
		newProduct.ValueCurrency = newProduct.PriceCurrency

		newProduct.Value = product.InventoryItem.UnitCost

		if product.InventoryItem.CountryCodeOfOrigin != nil {
			newProduct.CountryOfManufacture = *product.InventoryItem.CountryCodeOfOrigin
		}

		if product.InventoryItem.HarmonizedSystemCode != nil {
			newProduct.TariffCode = *product.InventoryItem.HarmonizedSystemCode
		}

		newProduct.Price = product.Price
		newProduct.CustomsValue = newProduct.Price

		//create product
		err = newProduct.Create(ctx)
		if err != nil {
			models.CreateSystemError(ctx, fmt.Sprintf("Error creating product for shopify product %s: %s", newShopifyProduct.VariantShopifyGraphqlID, err.Error()))
			return err
		}

		//create shopify product
		newShopifyProduct.ProductID = newProduct.ID
		err = newShopifyProduct.Create(ctx)
		if err != nil {
			models.CreateSystemError(ctx, fmt.Sprintf("Error creating shopify product %s: %s", newShopifyProduct.VariantShopifyGraphqlID, err.Error()))
			return err
		}

		return nil
	}

	return nil
}
