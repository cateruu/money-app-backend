package data

import (
	"time"

	"github.com/cateruu/money-app-backend/internal/validator"
	"github.com/gofrs/uuid/v5"
)

type Expense struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Amount    float32   `json:"amount"`
	UserID    uuid.UUID `json:"userID"`
	Version   int       `json:"-"`
}

func ValidateExpense(v *validator.Validator, expense *Expense) {
	v.Check(expense.Name != "", "name", "must be provided")
	v.Check(len(expense.Name) <= 100, "name", "must not be longer than 100 characters")

	v.Check(expense.Type != "", "type", "must be provided")

	v.Check(expense.CreatedAt != time.Time{}, "time", "must be provided")

	v.Check(expense.Amount != 0, "amount", "must be provided")

	v.Check(expense.UserID.IsNil(), "userId", "must be provided")
}
