package models

import (
	"net/http"
)

func GetRequestingUser(r *http.Request) (*User, error) {
	userFromContext := r.Context().Value("user")
	if userFromContext == nil {
		return nil, ErrUserNotInContext
	}
	user := userFromContext.(User)
	return &user, nil
}
