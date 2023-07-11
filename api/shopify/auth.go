package shopify

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/Yamashou/gqlgenc/clientv2"
	"github.com/shipply-io/shipply-io-backend/api/shopify/gen"
	"github.com/shipply-io/shipply-io-backend/models"
)

type contextKey int

const (
	shopifyAppConfigKey contextKey = iota
)

type ShopifyAppConfig struct {
	ID          string
	Secret      string
	Scope       string
	RedirectURL string
	Webhooks    map[string]string
}

func WithContext(ctx context.Context, sac *ShopifyAppConfig) context.Context {
	return context.WithValue(ctx, shopifyAppConfigKey, sac)
}

func FromContext(ctx context.Context) *ShopifyAppConfig {
	if rv := ctx.Value(shopifyAppConfigKey); rv != nil {
		return rv.(*ShopifyAppConfig)
	}
	return nil
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// Verify a message against a message HMAC
func (sac *ShopifyAppConfig) VerifyMessage(message string, messageMAC string) bool {
	mac := hmac.New(sha256.New, []byte(sac.Secret))
	mac.Write([]byte(message))
	expectedMAC := mac.Sum(nil)

	// shopify HMAC is in hex so it needs to be decoded
	actualMac, _ := hex.DecodeString(messageMAC)

	return hmac.Equal(actualMac, expectedMAC)
}

func (sac *ShopifyAppConfig) VerifyOAuthCallback(u *url.URL, appNonce string) error {

	q := u.Query()
	messageMAC := q.Get("hmac")

	nonce := q.Get("state")
	if nonce != appNonce {
		return errors.New("invalid nonce")
	}

	q.Del("hmac")
	q.Del("signature")

	message, err := url.QueryUnescape(q.Encode())
	if err != nil {
		return err
	}

	if !sac.VerifyMessage(message, messageMAC) {
		return errors.New("failed to verify hmac")
	}

	return nil

}

func isValidShopName(shopName string) bool {
	if !strings.HasSuffix(shopName, ".myshopify.com") {
		return false
	}
	hostname := shopName[:len(shopName)-15] // remove .myshopify.com suffix
	regex := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\-]*$`)
	return regex.MatchString(hostname)
}

func (sac *ShopifyAppConfig) AuthorizeUrl(shopName string, state string) string {
	shopUrl, _ := url.Parse("https://" + shopName)
	shopUrl.Path = "/admin/oauth/authorize"
	query := shopUrl.Query()
	query.Set("client_id", sac.ID)
	query.Set("redirect_uri", sac.RedirectURL)
	query.Set("scope", sac.Scope)
	query.Set("state", state)
	shopUrl.RawQuery = query.Encode()
	return shopUrl.String()
}

func (sac *ShopifyAppConfig) VerifyAuthorizationURL(u *url.URL) (bool, error) {
	q := u.Query()
	messageMAC := q.Get("hmac")

	// Remove hmac and signature and leave the rest of the parameters alone.
	q.Del("hmac")
	q.Del("signature")

	message, err := url.QueryUnescape(q.Encode())

	return sac.VerifyMessage(message, messageMAC), err
}

func (sac *ShopifyAppConfig) GetAccessToken(ctx context.Context, shopName string, code string) (string, error) {

	// Create a POST request to the Shopify access token endpoint
	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/admin/oauth/access_token", shopName), nil)
	if err != nil {
		return "", err
	}

	// Set the query parameters for the request
	q := req.URL.Query()
	q.Add("client_id", sac.ID)
	q.Add("client_secret", sac.Secret)
	q.Add("code", code)
	req.URL.RawQuery = q.Encode()

	// Send the request to Shopify
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body into a []byte
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Unmarshal the response body into an AccessTokenResponse struct
	var tokenResp AccessTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", err
	}

	models.LogShopifyAPIEvent(ctx, &models.ShopifyApiLog{
		ShopDomain: shopName,
		EventType:  "api_response",
		StatusCode: resp.StatusCode,
		RequestURL: req.URL.String(),
		Endpoint:   "access_token",
		Headers:    req.Header,
		Body:       string(body),
		Method:     req.Method,
	})

	// Return the access token
	return tokenResp.AccessToken, nil
}

func (sac *ShopifyAppConfig) VerifyWebhookRequest(r *http.Request) bool {
	shopifySha256 := r.Header.Get("X-Shopify-Hmac-Sha256")
	actualMac, err := base64.StdEncoding.DecodeString(shopifySha256)
	if err != nil {
		return false
	}

	// Convert the base64-decoded HMAC to hexadecimal
	hexActualMac := hex.EncodeToString(actualMac)

	requestBody, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(requestBody))

	return sac.VerifyMessage(string(requestBody), hexActualMac)
}

func UninstallApp(shopName string, accessToken string) error {

	req, err := http.NewRequest("DELETE", fmt.Sprintf("https://%s/admin/api_permissions/current.json", shopName), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Length", "0")
	req.Header.Set("X-Shopify-Access-Token", accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func NewClient(shopName string, accessToken string) *gen.Client {
	url := fmt.Sprintf("https://%s/admin/api/2023-01/graphql.json", shopName)
	return &gen.Client{
		Client: clientv2.NewClient(http.DefaultClient, url, nil, func(ctx context.Context, req *http.Request, gqlInfo *clientv2.GQLRequestInfo, res interface{}, next clientv2.RequestInterceptorFunc) error {
			req.Header.Set("X-Shopify-Access-Token", accessToken)
			return next(ctx, req, gqlInfo, res)
		}),
	}
}
