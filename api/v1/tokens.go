package v1

import (
	"errors"
	"net/http"
	"time"

	"github.com/cateruu/money-app-backend/internal/data"
	"github.com/cateruu/money-app-backend/internal/validator"
	"github.com/cateruu/money-app-backend/pkg/httperror"
	"github.com/cateruu/money-app-backend/pkg/json"
)

func (h *Handler) GenerateTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.ReadJSON(w, r, &input)
	if err != nil {
		httperror.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		httperror.FailedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := h.Models.UserModel.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			httperror.InvalidCredentialsResponse(w, r)
		default:
			httperror.ServerErrorResponse(w, r, err)
		}

		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		httperror.ServerErrorResponse(w, r, err)
		return
	}

	if !match {
		httperror.InvalidCredentialsResponse(w, r)
		return
	}

	token, err := h.Models.TokenModel.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		httperror.ServerErrorResponse(w, r, err)
		return
	}

	err = json.WriteJSON(w, http.StatusCreated, json.Envelope{"token": token}, nil)
	if err != nil {
		httperror.ServerErrorResponse(w, r, err)
	}
}
