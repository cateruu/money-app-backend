package main

import (
	"net/http"

	v1 "github.com/cateruu/money-app-backend/api/v1"
	"github.com/cateruu/money-app-backend/internal/data"
	"github.com/cateruu/money-app-backend/pkg/middleware"
)

func (app *app) routes(models *data.Models) http.Handler {
	routes := v1.NewHandler(models)
	handler := middleware.NewHandler(models)
	router := http.NewServeMux()

	router.HandleFunc("POST /v1/expenses", routes.CreateExpenseHandler)
	router.HandleFunc("GET /v1/expenses/{id}", routes.GetExpenseHandler)

	router.HandleFunc("POST /v1/users", routes.RegisterUserHandler)

	router.HandleFunc("POST /v1/tokens/authentication", routes.GenerateTokenHandler)

	return handler.RecoverPanic(handler.Authenticate(router))
}
