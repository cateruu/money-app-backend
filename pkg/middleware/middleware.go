package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/cateruu/money-app-backend/internal/data"
	"github.com/cateruu/money-app-backend/internal/validator"
	"github.com/cateruu/money-app-backend/pkg/httpcontext"
	"github.com/cateruu/money-app-backend/pkg/httperror"
)

type Handler struct {
	Models *data.Models
}

func NewHandler(models *data.Models) *Handler {
	return &Handler{
		Models: models,
	}
}

func (h *Handler) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				httperror.ServerErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			r = httpcontext.SetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			httperror.InvalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headerParts[1]

		v := validator.New()

		if data.ValidatePlaintextToken(v, token); !v.Valid() {
			httperror.FailedValidationResponse(w, r, v.Errors)
			return
		}

		user, err := h.Models.UserModel.GetForToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				httperror.InvalidAuthenticationTokenResponse(w, r)
			default:
				httperror.ServerErrorResponse(w, r, err)
			}

			return
		}

		r = httpcontext.SetUser(r, user)

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) ProtectedRoute(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := httpcontext.GetUser(r)

		if user.IsAnynonymous() {
			httperror.FailedAuthroizationResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
