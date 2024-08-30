package v1

import "github.com/cateruu/money-app-backend/internal/data"

type Handler struct {
	Models *data.Models
}

func NewHandler(models *data.Models) *Handler {
	return &Handler{
		Models: models,
	}
}
