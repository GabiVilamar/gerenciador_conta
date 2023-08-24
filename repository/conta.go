package repository

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Conta interface {
	Get(id int) (*ContaModel, error)
	GetAll() ([]ContaModel, error)
	Save(echo.Context, ContaModel) error
	Remove(id int) error
}

type conta struct {
	db *gorm.DB
}

func NewConta() Conta {
	return &conta{db: DB}
}

func (c *conta) Get(id int) (*ContaModel, error) {
	conta := &ContaModel{}
	result := c.db.First(conta, id)

	return conta, result.Error
}

func (c *conta) GetAll() ([]ContaModel, error) {
	var contas []ContaModel
	result := c.db.Find(&contas)

	return contas, result.Error
}

func (c *conta) Save(ctx echo.Context, conta ContaModel) error {
	db := c.db
	if ctx.Get("db_transaction") != nil {
		db = ctx.Get("db_transaction").(*gorm.DB)
	}
	result := db.Save(&conta)

	return result.Error
}

func (c *conta) Remove(id int) error {

	err := c.db.Delete(&ContaModel{}, id).Error
	return err
}
