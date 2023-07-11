package util

import (
	"context"

	"gorm.io/gorm"
)

type contextKey int

const (
	dbKey contextKey = iota
	authSecretKey
	frontendBaseURLKey
	shopifyKey
	cdnKey
)

func ContextWithCDN(ctx context.Context, cdn string) context.Context {
	return context.WithValue(ctx, cdnKey, cdn)
}

func CDNFromContext(ctx context.Context) string {
	if rv := ctx.Value(cdnKey); rv != nil {
		return rv.(string)
	}
	return ""
}

func ContextWithAuthSecret(ctx context.Context, secret string) context.Context {
	return context.WithValue(ctx, authSecretKey, secret)
}

func AuthSecretFromContext(ctx context.Context) string {
	if rv := ctx.Value(authSecretKey); rv != nil {
		return rv.(string)
	}
	return ""
}

func ContextWithFrontendBaseURL(ctx context.Context, url string) context.Context {
	return context.WithValue(ctx, frontendBaseURLKey, url)
}

func FrontendBaseURLFromContext(ctx context.Context) string {
	if rv := ctx.Value(frontendBaseURLKey); rv != nil {
		return rv.(string)
	}
	return ""
}

func ContextWithDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, dbKey, db)
}

func DBFromContext(ctx context.Context) *gorm.DB {
	if rv := ctx.Value(dbKey); rv != nil {
		return rv.(*gorm.DB)
	}
	return nil
}
