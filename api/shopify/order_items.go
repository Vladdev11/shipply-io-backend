package shopify

import (
	"context"
	"errors"
	"time"

	"github.com/shipply-io/shipply-io-backend/api/shopify/gen"
)

func RetrieveOrderItemsByOrderID(shopName string, accessToken string, orderID string) ([]gen.GetLineItemsByOrderID_Order_LineItems_Edges_Node, error) {
	client := NewClient(shopName, accessToken)

	var cursor *string
	orderItems := make([]gen.GetLineItemsByOrderID_Order_LineItems_Edges_Node, 0)

	for {
		retries := 0
		var getLineItemsResponse *gen.GetLineItemsByOrderID
		var err error

		for {
			getLineItemsResponse, err = client.GetLineItemsByOrderID(context.Background(), orderID, 10, cursor)
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

		for _, edge := range getLineItemsResponse.Order.LineItems.Edges {
			orderItems = append(orderItems, edge.Node)
		}

		if !getLineItemsResponse.Order.LineItems.PageInfo.HasNextPage {
			break
		}

		cursor = getLineItemsResponse.Order.LineItems.PageInfo.EndCursor
	}

	return orderItems, nil
}
