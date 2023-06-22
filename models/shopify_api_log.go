package models

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/shipply-io/shipply-io-backend/util"
)

type ShopifyApiLog struct {
	ID         int64
	EventType  string      // webhook or api_request
	Headers    http.Header `gorm:"type:jsonb"`
	Body       string
	Endpoint   string
	ShopDomain string
	StatusCode int
	RequestURL string
	Method     string
	CreatedAt  time.Time
}

func (log *ShopifyApiLog) Create() error {
	return PGDB.Create(log).Error
}

func LogShopifyWebhookRequest(r *http.Request) error {

	bodyCopy := &bytes.Buffer{}
	bodyReader := io.TeeReader(r.Body, bodyCopy)

	log := ShopifyApiLog{
		EventType:  "webhook",
		Headers:    r.Header,
		Body:       bodyCopy.String(),
		Endpoint:   r.URL.Path,
		ShopDomain: r.Header.Get("X-Shopify-Shop-Domain"),
		RequestURL: r.URL.String(),
		Method:     r.Method,
		CreatedAt:  time.Now(),
	}

	r.Body = ioutil.NopCloser(bodyReader)
	return log.Create()
}

func LogShopifyInstallRequest(r *http.Request) error {

	shop, err := util.GetStringQueryParam(r, "shop")
	if err != nil {
		return err
	}

	log := ShopifyApiLog{
		ShopDomain: shop,
		EventType:  "install",
		Headers:    r.Header,
		Endpoint:   r.URL.Path,
		Method:     r.Method,
	}
	return log.Create()
}

func LogShopifyAPIEvent(log *ShopifyApiLog) error {
	return log.Create()
}
