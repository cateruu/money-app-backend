package data

import (
	"context"
	"errors"
	"time"

	"github.com/cateruu/money-app-backend/internal/validator"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Expense struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Amount    float64   `json:"amount"`
	UserID    uuid.UUID `json:"userID"`
	Date      time.Time `json:"date"`
	Version   int       `json:"-"`
}

func ValidateExpense(v *validator.Validator, expense *Expense) {
	v.Check(expense.Name != "", "name", "must be provided")
	v.Check(len(expense.Name) <= 100, "name", "must not be longer than 100 characters")

	v.Check(expense.Type != "", "type", "must be provided")

	v.Check(expense.Amount != 0, "amount", "must be provided")

	v.Check(!expense.UserID.IsNil(), "userId", "must be provided")

	v.Check(expense.Date != time.Time{}, "date", "must be provided")
}

type ExpenseModel struct {
	DB *pgxpool.Pool
}

func (m ExpenseModel) Insert(expense *Expense) error {
	query := `INSERT INTO expenses (name, type, amount, date, user_id)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at, version`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{expense.Name, expense.Type, expense.Amount, expense.Date, expense.UserID}

	err := m.DB.QueryRow(ctx, query, args...).Scan(&expense.ID, &expense.CreatedAt, &expense.Version)
	return err
}

func (m ExpenseModel) GetByID(id uuid.UUID) (*Expense, error) {
	query := `SELECT id, created_at, name, type, amount, date, user_id, version
	FROM expenses 
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var expense Expense

	err := m.DB.QueryRow(ctx, query, id).Scan(
		&expense.ID,
		&expense.CreatedAt,
		&expense.Name,
		&expense.Type,
		&expense.Amount,
		&expense.Date,
		&expense.UserID,
		&expense.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &expense, nil
}

func (m ExpenseModel) Update(expense *Expense) error {
	query := `UPDATE expenses
	SET name = $1, type = $2, amount = $3, date = $4, version = version + 1
	WHERE id = $5 and version = $6
	RETURNING version`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{expense.Name, expense.Type, expense.Amount, expense.Date, expense.ID, expense.Version}

	err := m.DB.QueryRow(ctx, query, args...).Scan(&expense.Version)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	return nil
}

func (m ExpenseModel) Remove(id uuid.UUID) error {
	query := `DELETE FROM expenses
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.Exec(ctx, query, id)
	return err
}
