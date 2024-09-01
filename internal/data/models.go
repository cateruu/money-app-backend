package data

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	UserModel    UserModel
	TokenModel   TokenModel
	ExpenseModel ExpenseModel
}

func NewModels(db *pgxpool.Pool) Models {
	return Models{
		UserModel:    UserModel{DB: db},
		TokenModel:   TokenModel{DB: db},
		ExpenseModel: ExpenseModel{DB: db},
	}
}
