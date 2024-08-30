package v1

import (
	"errors"
	"net/http"

	"github.com/cateruu/money-app-backend/internal/data"
	"github.com/cateruu/money-app-backend/internal/validator"
	"github.com/cateruu/money-app-backend/pkg/httperror"
	"github.com/cateruu/money-app-backend/pkg/json"
)

func (h *Handler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.ReadJSON(w, r, &input)
	if err != nil {
		httperror.BadRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:  input.Name,
		Email: input.Email,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		httperror.ServerErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		httperror.FailedValidationResponse(w, r, v.Errors)
		return
	}

	err = h.Models.UserModel.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			httperror.FailedValidationResponse(w, r, v.Errors)
		default:
			httperror.ServerErrorResponse(w, r, err)
		}

		return
	}

	err = json.WriteJSON(w, http.StatusCreated, json.Envelope{"user": user}, nil)
	if err != nil {
		httperror.ServerErrorResponse(w, r, err)
	}
}
