package util

import "github.com/knadh/koanf/v2"

var (
	//ConfigPgAddr defines the default PostgreSQL address flag
	ConfigPgAddr string
	//ConfigPgUsername defines the default PostgreSQL username flag
	ConfigPgUsername string
	//ConfigPgPassword defines the default PostgreSQL password flag
	ConfigPgPassword string
	//ConfigPgDatabase defines the default PostgreSQL database flag
	ConfigPgDatabase string

	// ConfigAWSAccessKeyID defines the default S3 access key
	ConfigAWSAccessKeyID string
	// ConfigAWSSecretAccessKey defines the default S3 secret key
	ConfigAWSSecretAccessKey string

	// ConfigS3Region defines the default S3 region
	ConfigS3Region string
	// ConfigS3AttachmentBucket defines the default S3 bucket for attachments
	ConfigS3AttachmentBucket string
	// ConfigS3CDNBucket defines the default S3 bucket for cdn
	ConfigS3CDNBucket string

	//ConfigCDNHost defines the default CDN host
	ConfigCDNHost string

	//ConfigAuthSecretSigning defines the default secret used to sign JWT tokens
	ConfigAuthSecretSigning string
	//ConfigLocalPort
	ConfigLocalPort string

	//ConfigShipengineAPIKey
	// ConfigShipengineAPIKey string
	ConfigShipengineAPIKey string
	//ConfigShipengineAPIHost
	ConfigShipengineAPIHost string

	//ShopifyClientID
	ShopifyClientID string
	//ShopifyClientSecret
	ShopifyClientSecret string
	//ShopifyRedirectURL
	ShopifyRedirectURL string
	//ShopifyScope
	ShopifyScope string
	// ShopifyAppUninstalledWebhookURL
	ShopifyAppUninstalledWebhookURL string
	// ShopifyAppOrderCreatedWebhookURL
	ShopifyAppOrdersCreateWebhookURL string
	// ShopifyAppOrderUpdatedWebhookURL
	ShopifyAppOrdersUpdateWebhookURL string
	// ShopifyAppProductsCreateWebhookURL
	ShopifyAppProductsCreateWebhookURL string
	// ShopifyAppProductsUpdateWebhookURL
	ShopifyAppProductsUpdateWebhookURL string

	//SendgridAPIKey
	SendgridAPIKey string
	// SendgridFromEmail
	SendgridFromEmail string
	// SendgridResetPasswordTemplateID
	SendgridResetPasswordTemplateID string

	//FrontendURL
	FrontEndBaseURL string
)

// TODO: Delete all of this later, this is just a temporary solution.
// Proper way would be to use the context.Context type, store individual sub-structs of the config in the context, and then inject it on server start.
// So all the requests downstream have access to the config, instead of using global variables.
// Main issue with using global variables is of course it's a nightmare to test.
// But for now, this is the easiest way to get the config to the handlers.
// Since doing Contexts means we have to change the handler signatures, and then change the way we access these variables in the handlers.
// And there's a lot of handlers :)
func LoadConfig(k *koanf.Koanf) {
	ConfigLocalPort = k.String("server.port")
	FrontEndBaseURL = k.String("frontend.base_url")
	ConfigAuthSecretSigning = k.String("auth.secret")

	ConfigPgAddr = k.String("postgres.hostname")
	ConfigPgUsername = k.String("postgres.username")
	ConfigPgPassword = k.String("postgres.password")
	ConfigPgDatabase = k.String("postgres.database")

	ConfigAWSAccessKeyID = k.String("aws.access_key_id")
	ConfigAWSSecretAccessKey = k.String("aws.secret_access_key")

	ConfigS3Region = k.String("aws.region")
	ConfigS3AttachmentBucket = k.String("aws.attachment_bucket")
	ConfigS3CDNBucket = k.String("aws.cdn_bucket")

	ConfigCDNHost = k.String("aws.cdn")

	ConfigShipengineAPIKey = k.String("shipengine.api_key")
	ConfigShipengineAPIHost = k.String("shipengine.api_host")

	ShopifyClientID = k.String("shopify.client_id")
	ShopifyClientSecret = k.String("shopify.client_secret")
	ShopifyRedirectURL = k.String("shopify.redirect_url")
	ShopifyScope = k.String("shopify.scope")
	ShopifyAppUninstalledWebhookURL = k.String("shopify.webhooks.app_uninstalled")
	ShopifyAppOrdersCreateWebhookURL = k.String("shopify.webhooks.app_orders_create")
	ShopifyAppOrdersUpdateWebhookURL = k.String("shopify.webhooks.app_orders_update")
	ShopifyAppProductsCreateWebhookURL = k.String("shopify.webhooks.app_products_create")
	ShopifyAppProductsUpdateWebhookURL = k.String("shopify.webhooks.app_products_update")

	SendgridAPIKey = k.String("sendgrid.api_key")
	SendgridFromEmail = k.String("sendgrid.from_email")
	SendgridResetPasswordTemplateID = k.String("sendgrid.reset_password_template_id")
}
