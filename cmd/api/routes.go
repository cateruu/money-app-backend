package main

import (
	"net/http"

	v1 "github.com/cateruu/money-app-backend/api/v1"
	"github.com/cateruu/money-app-backend/internal/data"
	"github.com/cateruu/money-app-backend/pkg/middleware"
)

func (app *app) routes(models *data.Models) http.Handler {
	handler := v1.NewHandler(models)
	router := http.NewServeMux()

	router.HandleFunc("POST /v1/expenses", handler.CreateExpenseHandler)
	router.HandleFunc("GET /v1/expenses/{id}", handler.GetExpenseHandler)

	router.HandleFunc("POST /v1/users", handler.RegisterUserHandler)

	return middleware.RecoverPanic(router)
}
