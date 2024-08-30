package httpcontext

import (
	"context"
	"net/http"

	"github.com/cateruu/money-app-backend/internal/data"
)

type ContextKey string

const UserContextKey = ContextKey("user")

func SetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContextKey, user)

	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(UserContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}
