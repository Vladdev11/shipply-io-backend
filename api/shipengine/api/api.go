package ShipengineAPI

import (
	"net/http"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
)

var apiClient = NewAPIClient()

type APIClient struct {
	httpClient *http.Client
	apiHost    string
}

func NewAPIClient() *APIClient {
	return &APIClient{
		httpClient: &http.Client{
			Timeout: time.Second * 10,
			Transport: &http.Transport{
				MaxIdleConns:          100,
				IdleConnTimeout:       time.Second * 90,
				DisableCompression:    true,
				ResponseHeaderTimeout: time.Second * 10,
			},
		},
		apiHost: util.ConfigShipengineAPIHost,
	}
}

func (c *APIClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("API-Key", util.ConfigShipengineAPIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-us")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	return c.httpClient.Do(req)
}

func (c *APIClient) GetApiHost() string {
	return c.apiHost
}
