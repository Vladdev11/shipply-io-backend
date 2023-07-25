package shopify

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

var (
	//ErrInvalidJSON
	ErrInvalidJSON = errors.New("failed to parse json")
	//ErrInvalidShopifyGraphqlOrderID
	ErrInvalidShopifyGraphqlOrderID = errors.New("invalid graphql shopify order id")
	//ErrInvalidShopifyGraphqlProductID
	ErrInvalidShopifyGraphqlProductID = errors.New("invalid graphql shopify product id")
)

// GraphQLResponse represents a generic GraphQL response
type GraphQLResponse struct {
	NetworkErrors interface{}    `json:"networkErrors"`
	GraphQLErrors []GraphQLError `json:"graphqlErrors"`
}

// GraphQLError represents an individual GraphQL error
type GraphQLError struct {
	Message    string     `json:"message"`
	Extensions GraphQLEXT `json:"extensions"`
}

// GraphQLEXT represents the extensions field in the GraphQLError struct
type GraphQLEXT struct {
	Code          string `json:"code"`
	Cost          int    `json:"cost"`
	Documentation string `json:"documentation"`
	MaxCost       int    `json:"maxCost"`
}

const (
	ShopifyRateLimitError          = 1
	ShopifyQueryLimitExceededError = 2
)

func DetermineShopifyGraphqlError(data error) (int, error) {
	var response GraphQLResponse
	err := json.Unmarshal([]byte(fmt.Sprintf("%s", data)), &response)
	if err != nil {
		return 0, errors.New("failed to parse shopify graphql error response")
	}

	//determine error type
	if len(response.GraphQLErrors) > 0 {
		for _, err := range response.GraphQLErrors {
			//if the error is a rate limit error
			if strings.Contains(strings.ToLower(err.Message), "throttled") {
				return ShopifyRateLimitError, nil
			}

			//if the error is a query limit exceeded error
			if strings.Contains(strings.ToLower(err.Message), "exceeds the single query max cost limit") {
				return ShopifyQueryLimitExceededError, nil
			}
		}
	}

	return 0, errors.New("unknown shopify graphql error")
}

const (
	MaxRetries          = 20
	InitialWaitDuration = time.Second * 2
	JitterFactor        = 0.5
)

func ExponentialBackoffWithJitter(retries int, initialDuration time.Duration, jitterFactor float64) time.Duration {
	duration := float64(initialDuration) * math.Pow(2, float64(retries))
	jitter := jitterFactor * float64(initialDuration) * (rand.Float64()*2 - 1)
	return time.Duration(duration + jitter)
}
