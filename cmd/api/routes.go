package main

import (
	"net/http"

	v1 "github.com/cateruu/money-app-backend/api/v1"
	"github.com/cateruu/money-app-backend/pkg/middleware"
)

func (app *app) routes() http.Handler {
	routes := v1.NewHandler(app.models)
	handler := middleware.NewHandler(app.models)

	router := http.NewServeMux()

	router.HandleFunc("POST /v1/expenses", handler.ProtectedRoute(routes.CreateExpenseHandler))
	router.HandleFunc("GET /v1/expenses/{id}", handler.ProtectedRoute(routes.GetExpenseHandler))
	router.HandleFunc("PATCH /v1/expenses/{id}", handler.ProtectedRoute(routes.UpdateExpenseHandler))
	router.HandleFunc("DELETE /v1/expenses/{id}", handler.ProtectedRoute(routes.DeleteExpenseHandler))

	router.HandleFunc("POST /v1/users", routes.RegisterUserHandler)

	router.HandleFunc("POST /v1/tokens/authentication", routes.GenerateTokenHandler)

	return handler.RecoverPanic(handler.Authenticate(router))
}
