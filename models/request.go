package models

import (
	"errors"
	"net/http"
)

func GetRequestingUser(r *http.Request) (*User, error) {
	userFromContext := r.Context().Value("user")
	if userFromContext == nil {
		return nil, errors.New("user not found in request context")
	}
	user := userFromContext.(User)
	return &user, nil
}
