package v1

import (
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

	expense := data.Expense{
		ID:        uuid.FromStringOrNil(idString),
		CreatedAt: time.Now(),
		Name:      "test",
		Type:      "test",
		UserID:    uuid.FromStringOrNil(idString),
		Version:   1,
	}

	err := json.WriteJSON(w, http.StatusOK, json.Envelope{"expense": expense}, nil)
	if err != nil {
		httperror.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) CreateExpenseHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name      string    `json:"name"`
		Type      string    `json:"type"`
		Amount    float32   `json:"amount"`
		CreatedAt time.Time `json:"createdAt"`
		UserID    uuid.UUID `json:"userId"`
	}

	err := json.ReadJSON(w, r, &input)
	if err != nil {
		httperror.BadRequestResponse(w, r, err)
		return
	}

	expense := &data.Expense{
		Name:      input.Name,
		Type:      input.Type,
		Amount:    input.Amount,
		CreatedAt: input.CreatedAt,
		UserID:    input.UserID,
	}

	validator := validator.New()
	if data.ValidateExpense(validator, expense); !validator.Valid() {
		httperror.FailedValidationResponse(w, r, validator.Errors)
		return
	}

	err = json.WriteJSON(w, http.StatusOK, json.Envelope{"movie": input}, nil)
	if err != nil {
		httperror.ServerErrorResponse(w, r, err)
		return
	}
}
