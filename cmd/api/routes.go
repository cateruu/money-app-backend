package main

import (
	"net/http"

	v1 "github.com/cateruu/money-app-backend/api/v1"
	"github.com/cateruu/money-app-backend/pkg/middleware"
)

func (app *app) routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("POST /v1/expenses", v1.CreateExpenseHandler)
	router.HandleFunc("GET /v1/expenses/{id}", v1.GetExpenseHandler)

	return middleware.RecoverPanic(router)
}
