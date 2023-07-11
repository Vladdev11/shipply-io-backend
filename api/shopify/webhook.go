package shopify

import (
	"context"
)

func SubscribeToAppUninstalledWebhook(ctx context.Context, shopName string, accessToken string) error {

	client := NewClient(shopName, accessToken)

	_, err := client.SubscribeToAppUninstallWebhook(context.Background(), FromContext(ctx).Webhooks["app_uninstalled"])
	if err != nil {
		return err
	}

	return nil
}

func SubscribeToOrdersCreateWebhook(ctx context.Context, shopName string, accessToken string) error {

	client := NewClient(shopName, accessToken)

	_, err := client.SubscribeToOrdersCreateWebhook(context.Background(), FromContext(ctx).Webhooks["orders_create"])
	if err != nil {
		return err
	}

	return nil
}

func SubscribeToOrdersUpdateWebhook(ctx context.Context, shopName string, accessToken string) error {

	client := NewClient(shopName, accessToken)

	_, err := client.SubscribeToOrdersUpdateWebhook(context.Background(), FromContext(ctx).Webhooks["orders_update"])
	if err != nil {
		return err
	}

	return nil
}

func SubscribeToProductsCreateWebhook(ctx context.Context, shopName string, accessToken string) error {

	client := NewClient(shopName, accessToken)

	_, err := client.SubscribeToProductsCreateWebhook(context.Background(), FromContext(ctx).Webhooks["products_create"])
	if err != nil {
		return err
	}

	return nil
}

func SubscribeToProductsUpdateWebhook(ctx context.Context, shopName string, accessToken string) error {

	client := NewClient(shopName, accessToken)

	_, err := client.SubscribeToProductsUpdateWebhook(context.Background(), FromContext(ctx).Webhooks["products_create"])
	if err != nil {
		return err
	}

	return nil
}
