package repository

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Operation interface {
	Save(echo.Context, *Operacao) error
	GetAllTransaction() ([]Operacao, error)
}

type operation struct {
	db *gorm.DB
}

func NewOperation() Operation {
	return &operation{db: DB}
}

func (o *operation) Save(ctx echo.Context, operation *Operacao) error {
	db := o.db
	if ctx.Get("db_transaction") != nil {
		db = ctx.Get("db_transaction").(*gorm.DB)
	}
	result := db.Save(operation)

	return result.Error
}

func (o *operation) GetAllTransaction() ([]Operacao, error) {
	var transactions []Operacao
	result := o.db.Find(&transactions)

	return transactions, result.Error
}
