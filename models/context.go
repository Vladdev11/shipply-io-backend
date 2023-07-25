package models

import "context"

// TODO reevaluate location of this file And change key

// Puts user model in context
func ContextWithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, "new-user", user)
}

// Gets user model from context
func UserFromContext(ctx context.Context) *User {
	if rv := ctx.Value("new-user"); rv != nil {
		return rv.(*User)
	}
	return nil
}
