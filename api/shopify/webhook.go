package shopify

import (
	"context"

	"github.com/shipply-io/shipply-io-backend/util"
)

func SubscribeToAppUninstalledWebhook(shopName string, accessToken string) error {

	client := NewClient(shopName, accessToken)

	_, err := client.SubscribeToAppUninstallWebhook(context.Background(), util.ShopifyAppUninstalledWebhookURL)
	if err != nil {
		return err
	}

	return nil
}

func SubscribeToOrdersCreateWebhook(shopName string, accessToken string) error {

	client := NewClient(shopName, accessToken)

	_, err := client.SubscribeToOrdersCreateWebhook(context.Background(), util.ShopifyAppOrdersCreateWebhookURL)
	if err != nil {
		return err
	}

	return nil
}

func SubscribeToOrdersUpdateWebhook(shopName string, accessToken string) error {

	client := NewClient(shopName, accessToken)

	_, err := client.SubscribeToOrdersUpdateWebhook(context.Background(), util.ShopifyAppOrdersUpdateWebhookURL)
	if err != nil {
		return err
	}

	return nil
}

func SubscribeToProductsCreateWebhook(shopName string, accessToken string) error {

	client := NewClient(shopName, accessToken)

	_, err := client.SubscribeToProductsCreateWebhook(context.Background(), util.ShopifyAppProductsCreateWebhookURL)
	if err != nil {
		return err
	}

	return nil
}

func SubscribeToProductsUpdateWebhook(shopName string, accessToken string) error {

	client := NewClient(shopName, accessToken)

	_, err := client.SubscribeToProductsUpdateWebhook(context.Background(), util.ShopifyAppProductsUpdateWebhookURL)
	if err != nil {
		return err
	}

	return nil
}
