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
	"gorm.io/gorm"
)

func RetrieveOrders(shopName string, accessToken string) ([]gen.GetOrders_Orders_Edges_Node, error) {
	client := NewClient(shopName, accessToken)

	var cursor *string
	orders := make([]gen.GetOrders_Orders_Edges_Node, 0)

	for {
		retries := 0
		var getOrdersResponse *gen.GetOrders
		var err error

		for {
			getOrdersResponse, err = client.GetOrders(context.Background(), 10, cursor)
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

		for _, edge := range getOrdersResponse.Orders.Edges {
			orders = append(orders, edge.Node)
		}

		if !getOrdersResponse.Orders.PageInfo.HasNextPage {
			break
		}

		cursor = getOrdersResponse.Orders.PageInfo.EndCursor
	}

	return orders, nil
}

func RetrieveOrder(shopName string, accessToken string, orderID string) (*gen.GetOrderById_Order, error) {
	client := NewClient(shopName, accessToken)

	retries := 0
	var getOrderResponse *gen.GetOrderByID
	var err error

	for {
		getOrderResponse, err = client.GetOrderByID(context.Background(), orderID)
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

	return getOrderResponse.Order, nil
}

type OrderCreateWebhookRequest struct {
	AdminGraphqlAPIID string `json:"admin_graphql_api_id"`
}

func ParseAndValidateOrderCreateRequest(r *http.Request) (*string, error) {

	var errs []error

	var orderID string

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, ErrInvalidJSON
	}

	aux := &struct {
		AdminGraphqlAPIID json.RawMessage `json:"admin_graphql_api_id"`
	}{}

	err = json.Unmarshal(body, aux)
	if err != nil {
		return nil, ErrInvalidJSON
	}

	if aux.AdminGraphqlAPIID != nil {
		if err := json.Unmarshal(aux.AdminGraphqlAPIID, &orderID); err != nil {
			errs = append(errs, ErrInvalidShopifyGraphqlOrderID)
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return &orderID, nil
}

func ParseAndValidateOrderUpdateRequest(r *http.Request) (*string, []string) {

	var orderID string

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
		if err := json.Unmarshal(aux.AdminGraphqlAPIID, &orderID); err != nil {
			errs = append(errs, "invalid order id")
		}
	}

	if len(errs) > 0 {
		return nil, errs
	}

	return &orderID, nil
}

func ConvertShopifyAddressToModelAddress(srcAddress *ShopifyGraphqlModelAddress) *models.Address {
	if srcAddress == nil {
		return nil
	}

	address := models.Address{
		FirstName:  util.SafeDereferenceString(srcAddress.FirstName),
		LastName:   util.SafeDereferenceString(srcAddress.LastName),
		Company:    util.SafeDereferenceString(srcAddress.Company),
		Street1:    util.SafeDereferenceString(srcAddress.Address1),
		Street2:    util.SafeDereferenceString(srcAddress.Address2),
		City:       util.SafeDereferenceString(srcAddress.City),
		State:      util.SafeDereferenceString(srcAddress.Province),
		PostalCode: util.SafeDereferenceString(srcAddress.Zip),
		Country:    util.SafeDereferenceString(srcAddress.Country),
		Phone:      util.SafeDereferenceString(srcAddress.Phone),
	}

	return &address
}

func GetOrCreateShippingMethod(ctx context.Context, storeID int, shippingMethodName string) (*int, error) {
	shippingMethod, err := models.GetShippingMethodByStoreAndName(ctx, storeID, shippingMethodName)
	if err == nil && shippingMethod != nil {
		return &shippingMethod.ID, nil
	} else if err == gorm.ErrRecordNotFound {
		shippingMethod := models.ShippingMethod{
			Name:    shippingMethodName,
			StoreID: storeID,
			Mapped:  false,
		}
		err = shippingMethod.Create(ctx)
		if err != nil {
			return nil, err
		}
		return &shippingMethod.ID, nil
	} else {
		return nil, err
	}
}

type ShopifyGraphqlModelOrder struct {
	ID                string
	CanNotifyCustomer bool
	CancelReason      *string
	CancelledAt       *time.Time
	ClientIP          *string
	Closed            bool
	ClosedAt          *time.Time
	Confirmed         bool
	CurrencyCode      string
	CustomAttributes  []struct {
		Key   string
		Value *string
	}
	SubtotalPrice            float64
	TotalPrice               float64
	TotalTax                 float64
	TotalShipping            float64
	TotalDiscounts           float64
	TotalTip                 float64
	Currency                 string
	Email                    *string
	DisplayFinacialStatus    *string
	DisplayFulfillmentStatus string
	Edited                   bool
	Fulfillable              bool
	Name                     string
	Note                     *string
	Phone                    *string
	Refundable               bool
	Refunds                  []struct {
		ID string
	}
	RequiresShipping bool
	ShippingLine     *ShopifyGraphqlModelShippingLine
	RiskLevel        string
	Tags             []string
	PhysicalLocation *string
	BillingAddress   *ShopifyGraphqlModelAddress
	ShippingAddress  *ShopifyGraphqlModelAddress
	CreatedAt        time.Time
	DiscountCodes    []string
}

type ShopifyGraphqlModelShippingLine struct {
	Custom            bool
	Code              *string
	Price             string
	Title             string
	CarrierIdentifier *string
	Phone             *string
	Source            *string
	DiscountedPrice   float64
}

type ShopifyGraphqlModelAddress struct {
	Address1      *string
	Address2      *string
	City          *string
	Company       *string
	Country       *string
	CountryCodeV2 *string
	FirstName     *string
	LastName      *string
	Phone         *string
	Province      *string
	ProvinceCode  *string
	Zip           *string
}

func ConvertGetOrders_Orders_Edges_NodeToShopifyModelOrder(order gen.GetOrders_Orders_Edges_Node) (ShopifyGraphqlModelOrder, error) {

	var err error

	mOrder := ShopifyGraphqlModelOrder{}
	mOrder.ID = order.ID
	mOrder.CanNotifyCustomer = order.CanNotifyCustomer

	if order.CancelReason != nil {
		stringValue := string(*order.CancelReason)
		mOrder.CancelReason = &stringValue
	}

	if order.CancelledAt != nil {
		timeValue, err := util.ParseRFC3339Date(*order.CancelledAt)
		if err != nil {
			return mOrder, err
		}
		mOrder.CancelledAt = &timeValue
	}

	if order.ClientIP != nil {
		mOrder.ClientIP = order.ClientIP
	}

	mOrder.Closed = order.Closed

	if order.ClosedAt != nil {
		timeValue, err := util.ParseRFC3339Date(*order.ClosedAt)
		if err != nil {
			return mOrder, err
		}
		mOrder.ClosedAt = &timeValue
	}

	mOrder.Confirmed = order.Confirmed
	mOrder.CurrencyCode = order.CurrencyCode.String()

	if order.CustomAttributes != nil {
		for _, attribute := range order.CustomAttributes {
			mOrder.CustomAttributes = append(mOrder.CustomAttributes, struct {
				Key   string
				Value *string
			}{
				Key:   attribute.Key,
				Value: attribute.Value,
			})
		}
	}

	mOrder.SubtotalPrice, err = strconv.ParseFloat(order.SubtotalPriceSet.ShopMoney.Amount, 64)
	if err != nil {
		mOrder.SubtotalPrice = 0
	}

	mOrder.TotalPrice, err = strconv.ParseFloat(order.TotalPriceSet.ShopMoney.Amount, 64)
	if err != nil {
		mOrder.TotalPrice = 0
	}
	mOrder.Currency = order.TotalPriceSet.ShopMoney.CurrencyCode.String()

	mOrder.TotalTax, err = strconv.ParseFloat(order.TotalTaxSet.ShopMoney.Amount, 64)
	if err != nil {
		mOrder.TotalTax = 0
	}

	mOrder.TotalShipping, err = strconv.ParseFloat(order.TotalShippingPriceSet.ShopMoney.Amount, 64)
	if err != nil {
		mOrder.TotalShipping = 0
	}

	mOrder.TotalDiscounts, err = strconv.ParseFloat(order.TotalDiscountsSet.ShopMoney.Amount, 64)
	if err != nil {
		mOrder.TotalDiscounts = 0
	}

	mOrder.TotalTip, err = strconv.ParseFloat(order.TotalTipReceivedSet.ShopMoney.Amount, 64)
	if err != nil {
		mOrder.TotalTip = 0
	}

	if order.Email != nil {
		mOrder.Email = order.Email
	}

	if order.DisplayFinancialStatus != nil {
		mOrder.DisplayFinacialStatus = (*string)(order.DisplayFinancialStatus)
	}

	mOrder.DisplayFulfillmentStatus = order.DisplayFulfillmentStatus.String()
	mOrder.Edited = order.Edited
	mOrder.Fulfillable = order.Fulfillable
	mOrder.Name = order.Name

	if order.Note != nil {
		mOrder.Note = order.Note
	}

	if order.Phone != nil {
		mOrder.Phone = order.Phone
	}

	mOrder.Refundable = order.Refundable

	if order.Refunds != nil {
		for _, refund := range order.Refunds {
			mOrder.Refunds = append(mOrder.Refunds, struct {
				ID string
			}{
				ID: refund.ID,
			})
		}
	}

	mOrder.RequiresShipping = order.RequiresShipping

	if order.ShippingLine != nil {
		//init shipping line
		if mOrder.ShippingLine == nil {
			mOrder.ShippingLine = &ShopifyGraphqlModelShippingLine{}
		}
		mOrder.ShippingLine.Custom = order.ShippingLine.Custom
		mOrder.ShippingLine.Code = order.ShippingLine.Code
		mOrder.ShippingLine.Price = order.ShippingLine.Price
		mOrder.ShippingLine.Title = order.ShippingLine.Title
		mOrder.ShippingLine.CarrierIdentifier = order.ShippingLine.CarrierIdentifier
		mOrder.ShippingLine.Phone = order.ShippingLine.Phone
		mOrder.ShippingLine.Source = order.ShippingLine.Source
		mOrder.ShippingLine.DiscountedPrice, err = strconv.ParseFloat(order.ShippingLine.DiscountedPriceSet.ShopMoney.Amount, 64)
		if err != nil {
			return mOrder, err
		}
	}

	mOrder.RiskLevel = order.RiskLevel.String()

	if order.Tags != nil {
		mOrder.Tags = order.Tags
	}

	if order.PhysicalLocation != nil {
		mOrder.PhysicalLocation = &order.PhysicalLocation.Name
	}

	if order.BillingAddress != nil {
		mOrder.BillingAddress = &ShopifyGraphqlModelAddress{}
		if order.BillingAddress.Address1 != nil {
			mOrder.BillingAddress.Address1 = order.BillingAddress.Address1
		}
		if order.BillingAddress.Address2 != nil {
			mOrder.BillingAddress.Address2 = order.BillingAddress.Address2
		}
		if order.BillingAddress.City != nil {
			mOrder.BillingAddress.City = order.BillingAddress.City
		}
		if order.BillingAddress.Company != nil {
			mOrder.BillingAddress.Company = order.BillingAddress.Company
		}
		if order.BillingAddress.Country != nil {
			mOrder.BillingAddress.Country = (*string)(order.BillingAddress.CountryCodeV2)
		}
		if order.BillingAddress.FirstName != nil {
			mOrder.BillingAddress.FirstName = order.BillingAddress.FirstName
		}
		if order.BillingAddress.LastName != nil {
			mOrder.BillingAddress.LastName = order.BillingAddress.LastName
		}
		if order.BillingAddress.Phone != nil {
			mOrder.BillingAddress.Phone = order.BillingAddress.Phone
		}
		if order.BillingAddress.Province != nil {
			mOrder.BillingAddress.Province = order.BillingAddress.Province
		}
		if order.BillingAddress.ProvinceCode != nil {
			mOrder.BillingAddress.ProvinceCode = order.BillingAddress.ProvinceCode
		}
		if order.BillingAddress.Zip != nil {
			mOrder.BillingAddress.Zip = order.BillingAddress.Zip
		}
	}

	if order.ShippingAddress != nil {
		mOrder.ShippingAddress = &ShopifyGraphqlModelAddress{}
		if order.ShippingAddress.Address1 != nil {
			mOrder.ShippingAddress.Address1 = order.ShippingAddress.Address1
		}
		if order.ShippingAddress.Address2 != nil {
			mOrder.ShippingAddress.Address2 = order.ShippingAddress.Address2
		}
		if order.ShippingAddress.City != nil {
			mOrder.ShippingAddress.City = order.ShippingAddress.City
		}
		if order.ShippingAddress.Company != nil {
			mOrder.ShippingAddress.Company = order.ShippingAddress.Company
		}
		if order.ShippingAddress.Country != nil {
			mOrder.ShippingAddress.Country = (*string)(order.ShippingAddress.CountryCodeV2)
		}
		if order.ShippingAddress.FirstName != nil {
			mOrder.ShippingAddress.FirstName = order.ShippingAddress.FirstName
		}
		if order.ShippingAddress.LastName != nil {
			mOrder.ShippingAddress.LastName = order.ShippingAddress.LastName
		}
		if order.ShippingAddress.Phone != nil {
			mOrder.ShippingAddress.Phone = order.ShippingAddress.Phone
		}
		if order.ShippingAddress.Province != nil {
			mOrder.ShippingAddress.Province = order.ShippingAddress.Province
		}
		if order.ShippingAddress.ProvinceCode != nil {
			mOrder.ShippingAddress.ProvinceCode = order.ShippingAddress.ProvinceCode
		}
		if order.ShippingAddress.Zip != nil {
			mOrder.ShippingAddress.Zip = order.ShippingAddress.Zip
		}
	}

	mOrder.CreatedAt, err = util.ParseRFC3339Date(order.CreatedAt)
	if err != nil {
		return mOrder, err
	}

	mOrder.DiscountCodes = order.DiscountCodes

	return mOrder, nil
}

func ConvertGetOrderById_OrderToShopifyModelOrder(order gen.GetOrderById_Order) (ShopifyGraphqlModelOrder, error) {
	var err error

	mOrder := ShopifyGraphqlModelOrder{}
	mOrder.ID = order.ID
	mOrder.CanNotifyCustomer = order.CanNotifyCustomer

	if order.CancelReason != nil {
		stringValue := order.CancelReason.String()
		mOrder.CancelReason = &stringValue
	}

	if order.CancelledAt != nil {
		timeValue, err := util.ParseRFC3339Date(*order.CancelledAt)
		if err != nil {
			return mOrder, err
		}
		mOrder.CancelledAt = &timeValue
	}

	if order.ClientIP != nil {
		mOrder.ClientIP = order.ClientIP
	}

	mOrder.Closed = order.Closed

	if order.ClosedAt != nil {
		timeValue, err := util.ParseRFC3339Date(*order.ClosedAt)
		if err != nil {
			return mOrder, err
		}
		mOrder.ClosedAt = &timeValue
	}

	mOrder.Confirmed = order.Confirmed
	mOrder.CurrencyCode = order.CurrencyCode.String()

	if order.CustomAttributes != nil {
		for _, attribute := range order.CustomAttributes {
			mOrder.CustomAttributes = append(mOrder.CustomAttributes, struct {
				Key   string
				Value *string
			}{
				Key:   attribute.Key,
				Value: attribute.Value,
			})
		}
	}

	mOrder.SubtotalPrice, err = strconv.ParseFloat(order.SubtotalPriceSet.ShopMoney.Amount, 64)
	if err != nil {
		mOrder.SubtotalPrice = 0
	}

	mOrder.TotalPrice, err = strconv.ParseFloat(order.TotalPriceSet.ShopMoney.Amount, 64)
	if err != nil {
		mOrder.TotalPrice = 0
	}
	mOrder.Currency = order.TotalPriceSet.ShopMoney.CurrencyCode.String()

	mOrder.TotalTax, err = strconv.ParseFloat(order.TotalTaxSet.ShopMoney.Amount, 64)
	if err != nil {
		mOrder.TotalTax = 0
	}

	mOrder.TotalShipping, err = strconv.ParseFloat(order.TotalShippingPriceSet.ShopMoney.Amount, 64)
	if err != nil {
		mOrder.TotalShipping = 0
	}

	mOrder.TotalDiscounts, err = strconv.ParseFloat(order.TotalDiscountsSet.ShopMoney.Amount, 64)
	if err != nil {
		mOrder.TotalDiscounts = 0
	}

	mOrder.TotalTip, err = strconv.ParseFloat(order.TotalTipReceivedSet.ShopMoney.Amount, 64)
	if err != nil {
		mOrder.TotalTip = 0
	}

	if order.Email != nil {
		mOrder.Email = order.Email
	}

	if order.DisplayFinancialStatus != nil {
		mOrder.DisplayFinacialStatus = (*string)(order.DisplayFinancialStatus)
	}

	mOrder.DisplayFulfillmentStatus = order.DisplayFulfillmentStatus.String()
	mOrder.Edited = order.Edited
	mOrder.Fulfillable = order.Fulfillable
	mOrder.Name = order.Name

	if order.Note != nil {
		mOrder.Note = order.Note
	}

	if order.Phone != nil {
		mOrder.Phone = order.Phone
	}

	mOrder.Refundable = order.Refundable

	if order.Refunds != nil {
		for _, refund := range order.Refunds {
			mOrder.Refunds = append(mOrder.Refunds, struct {
				ID string
			}{
				ID: refund.ID,
			})
		}
	}

	mOrder.RequiresShipping = order.RequiresShipping

	if order.ShippingLine != nil {
		//init shipping line
		if mOrder.ShippingLine == nil {
			mOrder.ShippingLine = &ShopifyGraphqlModelShippingLine{}
		}
		mOrder.ShippingLine.Custom = order.ShippingLine.Custom
		mOrder.ShippingLine.Code = order.ShippingLine.Code
		mOrder.ShippingLine.Price = order.ShippingLine.Price
		mOrder.ShippingLine.Title = order.ShippingLine.Title
		mOrder.ShippingLine.CarrierIdentifier = order.ShippingLine.CarrierIdentifier
		mOrder.ShippingLine.Phone = order.ShippingLine.Phone
		mOrder.ShippingLine.Source = order.ShippingLine.Source
		mOrder.ShippingLine.DiscountedPrice, err = strconv.ParseFloat(order.ShippingLine.DiscountedPriceSet.ShopMoney.Amount, 64)
		if err != nil {
			return mOrder, err
		}
	}

	mOrder.RiskLevel = order.RiskLevel.String()

	if order.Tags != nil {
		mOrder.Tags = order.Tags
	}

	if order.PhysicalLocation != nil {
		mOrder.PhysicalLocation = &order.PhysicalLocation.Name
	}

	if order.BillingAddress != nil {
		mOrder.BillingAddress = &ShopifyGraphqlModelAddress{}
		if order.BillingAddress.Address1 != nil {
			mOrder.BillingAddress.Address1 = order.BillingAddress.Address1
		}
		if order.BillingAddress.Address2 != nil {
			mOrder.BillingAddress.Address2 = order.BillingAddress.Address2
		}
		if order.BillingAddress.City != nil {
			mOrder.BillingAddress.City = order.BillingAddress.City
		}
		if order.BillingAddress.Company != nil {
			mOrder.BillingAddress.Company = order.BillingAddress.Company
		}
		if order.BillingAddress.Country != nil {
			mOrder.BillingAddress.Country = (*string)(order.BillingAddress.CountryCodeV2)
		}
		if order.BillingAddress.FirstName != nil {
			mOrder.BillingAddress.FirstName = order.BillingAddress.FirstName
		}
		if order.BillingAddress.LastName != nil {
			mOrder.BillingAddress.LastName = order.BillingAddress.LastName
		}
		if order.BillingAddress.Phone != nil {
			mOrder.BillingAddress.Phone = order.BillingAddress.Phone
		}
		if order.BillingAddress.Province != nil {
			mOrder.BillingAddress.Province = order.BillingAddress.Province
		}
		if order.BillingAddress.ProvinceCode != nil {
			mOrder.BillingAddress.ProvinceCode = order.BillingAddress.ProvinceCode
		}
		if order.BillingAddress.Zip != nil {
			mOrder.BillingAddress.Zip = order.BillingAddress.Zip
		}
	}

	if order.ShippingAddress != nil {
		mOrder.ShippingAddress = &ShopifyGraphqlModelAddress{}
		if order.ShippingAddress.Address1 != nil {
			mOrder.ShippingAddress.Address1 = order.ShippingAddress.Address1
		}
		if order.ShippingAddress.Address2 != nil {
			mOrder.ShippingAddress.Address2 = order.ShippingAddress.Address2
		}
		if order.ShippingAddress.City != nil {
			mOrder.ShippingAddress.City = order.ShippingAddress.City
		}
		if order.ShippingAddress.Company != nil {
			mOrder.ShippingAddress.Company = order.ShippingAddress.Company
		}
		if order.ShippingAddress.Country != nil {
			mOrder.ShippingAddress.Country = (*string)(order.ShippingAddress.CountryCodeV2)
		}
		if order.ShippingAddress.FirstName != nil {
			mOrder.ShippingAddress.FirstName = order.ShippingAddress.FirstName
		}
		if order.ShippingAddress.LastName != nil {
			mOrder.ShippingAddress.LastName = order.ShippingAddress.LastName
		}
		if order.ShippingAddress.Phone != nil {
			mOrder.ShippingAddress.Phone = order.ShippingAddress.Phone
		}
		if order.ShippingAddress.Province != nil {
			mOrder.ShippingAddress.Province = order.ShippingAddress.Province
		}
		if order.ShippingAddress.ProvinceCode != nil {
			mOrder.ShippingAddress.ProvinceCode = order.ShippingAddress.ProvinceCode
		}
		if order.ShippingAddress.Zip != nil {
			mOrder.ShippingAddress.Zip = order.ShippingAddress.Zip
		}
	}

	mOrder.CreatedAt, err = util.ParseRFC3339Date(order.CreatedAt)
	if err != nil {
		return mOrder, err
	}

	mOrder.DiscountCodes = order.DiscountCodes

	return mOrder, nil
}
