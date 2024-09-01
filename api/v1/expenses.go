package v1

import (
	"errors"
	"net/http"
	"time"

	"github.com/cateruu/money-app-backend/internal/data"
	"github.com/cateruu/money-app-backend/internal/validator"
	"github.com/cateruu/money-app-backend/pkg/httperror"
	"github.com/cateruu/money-app-backend/pkg/json"
	"github.com/gofrs/uuid/v5"
)

func (h *Handler) GetExpenseHandler(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")

	expense, err := h.Models.ExpenseModel.GetByID(uuid.FromStringOrNil(idString))
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			httperror.NotFoundResponse(w, r)
		default:
			httperror.ServerErrorResponse(w, r, err)
		}

		return
	}

	err = json.WriteJSON(w, http.StatusOK, json.Envelope{"expense": expense}, nil)
	if err != nil {
		httperror.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) CreateExpenseHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name   string    `json:"name"`
		Type   string    `json:"type"`
		Amount float64   `json:"amount"`
		Date   time.Time `json:"date"`
		UserID uuid.UUID `json:"userId"`
	}

	err := json.ReadJSON(w, r, &input)
	if err != nil {
		httperror.BadRequestResponse(w, r, err)
		return
	}

	expense := &data.Expense{
		Name:   input.Name,
		Type:   input.Type,
		Amount: input.Amount,
		Date:   input.Date,
		UserID: input.UserID,
	}

	validator := validator.New()
	if data.ValidateExpense(validator, expense); !validator.Valid() {
		httperror.FailedValidationResponse(w, r, validator.Errors)
		return
	}

	err = h.Models.ExpenseModel.Insert(expense)
	if err != nil {
		httperror.ServerErrorResponse(w, r, err)
		return
	}

	err = json.WriteJSON(w, http.StatusOK, json.Envelope{"expense": expense}, nil)
	if err != nil {
		httperror.ServerErrorResponse(w, r, err)
		return
	}
}

func (h *Handler) UpdateExpenseHandler(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")

	var input struct {
		Name   *string    `json:"name"`
		Type   *string    `json:"type"`
		Amount *float64   `json:"amount"`
		Date   *time.Time `json:"date"`
	}

	err := json.ReadJSON(w, r, &input)
	if err != nil {
		httperror.ServerErrorResponse(w, r, err)
		return
	}

	expense, err := h.Models.ExpenseModel.GetByID(uuid.FromStringOrNil(idString))
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			httperror.NotFoundResponse(w, r)
		default:
			httperror.ServerErrorResponse(w, r, err)
		}

		return
	}

	if input.Name != nil {
		expense.Name = *input.Name
	}

	if input.Type != nil {
		expense.Type = *input.Type
	}

	if input.Amount != nil {
		expense.Amount = *input.Amount
	}

	if input.Date != nil {
		expense.Date = *input.Date
	}

	err = h.Models.ExpenseModel.Update(expense)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			httperror.NotFoundResponse(w, r)
		default:
			httperror.ServerErrorResponse(w, r, err)
		}

		return
	}

	err = json.WriteJSON(w, http.StatusOK, json.Envelope{"expense": expense}, nil)
	if err != nil {
		httperror.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) DeleteExpenseHandler(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")

	err := h.Models.ExpenseModel.Remove(uuid.FromStringOrNil(idString))
	if err != nil {
		httperror.ServerErrorResponse(w, r, err)
		return
	}

	err = json.WriteJSON(w, http.StatusOK, json.Envelope{"success": true}, nil)
	if err != nil {
		httperror.ServerErrorResponse(w, r, err)
	}
}
