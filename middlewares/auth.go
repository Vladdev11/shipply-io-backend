package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/shipply-io/shipply-io-backend/models"
	"github.com/shipply-io/shipply-io-backend/util"
)

// AuthMiddleware is the middleware for authentication we run on all protected routes
func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			//get the authorization header confirm it is a bearer token
			authHeader := r.Header.Get("Authorization")
			tokenString := authHeader[7:]
			if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
				util.ErrorResponse(w, "No Token Provided", http.StatusBadRequest)
				return
			}

			token, err := jwt.Parse(string(tokenString), func(token *jwt.Token) (interface{}, error) {
				//check if the token is valid
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				//return the secret
				return []byte(secret), nil
			})

			if err != nil {
				util.ErrorResponse(w, strings.ToLower(err.Error()), http.StatusBadRequest)
				return
			}

			if !token.Valid {
				util.ErrorResponse(w, "Invalid Token", http.StatusBadRequest)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				util.ErrorResponse(w, "Invalid Token", http.StatusBadRequest)
				return
			}

			userID := int(claims["user_id"].(float64))
			user, err := models.GetUserByID(ctx, userID)
			if err != nil {
				util.ErrorResponse(w, "failed to find user", http.StatusBadRequest)
				return
			}
			// TODO clean this once all handlers are updated to new user context version
			// add the user to the context (old version -- remove once possible)
			ctx = context.WithValue(r.Context(), "user", user)

			// add the user to the context
			if user.IsClientRole() {
				err = user.GetClient(ctx)
				if err != nil {
					util.ErrorResponse(w, "failed to get client", http.StatusBadRequest)
					return
				}
			} else {
				err = user.GetOrganization(ctx)
				if err != nil {
					util.ErrorResponse(w, "failed to get organization", http.StatusBadRequest)
					return
				}
			}
			ctx = models.ContextWithUser(ctx, &user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
