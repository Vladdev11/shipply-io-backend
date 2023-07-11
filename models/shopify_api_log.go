package models

import (
	"bytes"
	"context"
	"io"
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

func (log *ShopifyApiLog) Create(ctx context.Context) error {
	return util.DBFromContext(ctx).Create(log).Error
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

	r.Body = io.NopCloser(bodyReader)
	return log.Create(r.Context())
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
	return log.Create(r.Context())
}

func LogShopifyAPIEvent(ctx context.Context, log *ShopifyApiLog) error {
	return log.Create(ctx)
}
