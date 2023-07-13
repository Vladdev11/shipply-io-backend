package ShipengineAPI

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

type contextKey int

const (
	shipengineKey contextKey = iota
)

type APIClient struct {
	httpClient *http.Client
	apiHost    string
	apiKey     string
}

func ContextWithShipengineClient(ctx context.Context, host string, key string) (context.Context, error) {
	u, err := url.Parse(host)
	if err != nil {
		return ctx, err
	}
	client := &APIClient{
		httpClient: &http.Client{
			Timeout: time.Second * 10,
			Transport: &http.Transport{
				MaxIdleConns:          100,
				IdleConnTimeout:       time.Second * 90,
				DisableCompression:    true,
				ResponseHeaderTimeout: time.Second * 10,
			},
		},
		apiHost: u.Host,
		apiKey:  key,
	}
	return context.WithValue(ctx, shipengineKey, client), nil
}

func (c *APIClient) Do(req *http.Request) (*http.Response, error) {
	// We use the API host and key from the client instead of the request
	// Request only has the path and body
	req.URL.Host = c.apiHost
	req.Header.Set("API-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-us")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	return c.httpClient.Do(req)
}

func ShipengineClientFromContext(ctx context.Context) *APIClient {
	if rv := ctx.Value(shipengineKey); rv != nil {
		return rv.(*APIClient)
	}
	return nil
}
