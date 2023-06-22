package models

import "github.com/shipply-io/shipply-io-backend/util"

type Marketplace struct {
	ID                int
	Name              string
	ThumbnailURL      string
	SmallThumbnailURL string
}

type MarketplaceReturnJSON struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	ThumbnailURL      string `json:"thumbnail_url"`
	SmallThumbnailURL string `json:"small_thumbnail_url"`
}

func (m *Marketplace) ConvertToReturnJSON() *MarketplaceReturnJSON {

	if m == nil {
		return nil
	}

	return &MarketplaceReturnJSON{
		ID:                m.ID,
		Name:              m.Name,
		ThumbnailURL:      m.ThumbnailURL,
		SmallThumbnailURL: m.SmallThumbnailURL,
	}
}

var Marketplaces = map[int]Marketplace{
	util.ShopifyMarketplaceID: {
		ID:                util.ShopifyMarketplaceID,
		Name:              "Shopify",
		ThumbnailURL:      "https://pdimage.petparty.co/warehance/shopify.jpg",
		SmallThumbnailURL: "https://pdimage.petparty.co/warehance/shopify-150px.jpg",
	},
}
