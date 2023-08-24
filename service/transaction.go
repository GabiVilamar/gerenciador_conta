package service

import (
	"errors"
	"gerenciador_conta/repository"
	"log"

	"github.com/labstack/echo/v4"
)

type TransactionService interface {
	CreateTranfer(c echo.Context, targetAccount int, sourceAccount *repository.ContaModel, value int) (repository.Operacao, error)
	CreateDeposit(c echo.Context, targetAccount int, value int) (repository.Operacao, error)
	CreateWithdrawal(c echo.Context, targetAccount int, value int) (repository.Operacao, error)
	// validTransaction(sourceAccount *repository.ContaModel, valueTransaction int) error
	// updateSourceBalance(c echo.Context, sourceAccount *repository.ContaModel, valueTransaction int) error
	// updateTargetBalance(c echo.Context, targetAccount *repository.ContaModel, valueTransaction int) error
}

type transactionService struct {
	rc repository.Conta
	ro repository.Operation
}

// type status struct {
// 	valid string
// 	invalid string
// }

func (ts *transactionService) validTransaction(sourceAccount *repository.ContaModel, valueTransaction int) error {
	if sourceAccount.Saldo < valueTransaction {
		return errors.New("saldo insuficiente")
	}

	return nil
}

func (ts *transactionService) updateSourceBalance(c echo.Context, sourceAccount *repository.ContaModel, valueTransaction int) error {
	err := ts.validTransaction(sourceAccount, valueTransaction)

	if err != nil {
		return errors.New("Erro ao atualizar saldo no banco!")
	}

	sourceAccount.Saldo = sourceAccount.Saldo - valueTransaction

	err = ts.rc.Save(c, *sourceAccount)
	if err != nil {
		return errors.New("Erro ao atualizar saldo no banco!")

	}
	return nil
}

func (ts *transactionService) updateTargetBalance(c echo.Context, targetAccount *repository.ContaModel, valueTransaction int) error {
	targetAccount.Saldo = targetAccount.Saldo + valueTransaction
	// targetAccount.Saldo = balanceTarget

	err := ts.rc.Save(c, *targetAccount)

	if err != nil {
		log.Fatal("Não foi possível atualizar o saldo.", err)

		return errors.New("Erro ao atualizar saldo no banco!")
	}
	return nil
}

func (ts *transactionService) CreateDeposit(c echo.Context, targetAccount int, value int) (repository.Operacao, error) {
	contaTarget, err := ts.rc.Get(targetAccount)

	if err != nil {
		return repository.Operacao{}, errors.New("Conta não encontrada")
	}

	contaTarget.Saldo = contaTarget.Saldo + value

	tx := repository.DB.Begin()
	c.Set("db_transaction", tx)

	err = ts.rc.Save(c, *contaTarget)

	if err != nil {
		log.Fatal("Não foi possível atualizar o saldo.", err)
		tx.Rollback()
		return repository.Operacao{}, errors.New("Erro ao salvar transação de depósito!")

	}

	operacao := repository.Operacao{
		Type:        1,
		ContaTarget: contaTarget,
		Value:       value,
	}

	err = ts.ro.Save(c, &operacao)
	if err != nil {
		log.Fatal("Não foi possível finalizar a operação.", err)
		tx.Rollback()

		return repository.Operacao{}, errors.New("Erro ao salvar operação!")
	}
	tx.Commit()

	return operacao, err
}

func (ts *transactionService) CreateWithdrawal(c echo.Context, targetAccount int, value int) (repository.Operacao, error) {
	contaTarget, err := ts.rc.Get(targetAccount)

	if err != nil {
		return repository.Operacao{}, errors.New("Conta não encontrada")
	}

	contaTarget.Saldo = contaTarget.Saldo - value

	tx := repository.DB.Begin()
	c.Set("db_transaction", tx)

	err = ts.rc.Save(c, *contaTarget)

	if err != nil {
		log.Fatal("Não foi possível atualizar o saldo.", err)
		tx.Rollback()
		return repository.Operacao{}, errors.New("Erro ao salvar transação de depósito!")

	}

	operacao := repository.Operacao{
		Type: 2,
		// ContaTargetID: &contaTarget.ID,
		ContaTarget: contaTarget,
		Value:       value,
	}

	err = ts.ro.Save(c, &operacao)
	if err != nil {
		log.Fatal("Não foi possível finalizar a operação.", err)
		tx.Rollback()

		return repository.Operacao{}, errors.New("Erro ao salvar operação!")
	}
	tx.Commit()

	return operacao, err
}

func NewTransactionService(rc repository.Conta, ro repository.Operation) TransactionService {
	return &transactionService{rc, ro}
}

func (ts *transactionService) CreateTranfer(c echo.Context, targetAccount int, sourceAccount *repository.ContaModel, value int) (repository.Operacao, error) {
	tx := repository.DB.Begin()
	c.Set("db_transaction", tx)

	contaTarget, err := ts.rc.Get(targetAccount)

	if err != nil {
		return repository.Operacao{}, errors.New("Conta não encontrada")
	}

	err = ts.updateSourceBalance(c, sourceAccount, value)
	if err != nil {
		tx.Rollback()
		return repository.Operacao{}, errors.New("Erro ao atualizar saldo no banco!")
	}

	err = ts.updateTargetBalance(c, contaTarget, value)
	if err != nil {
		tx.Rollback()
		return repository.Operacao{}, errors.New("Erro ao atualizar saldo no banco!")
	}

	// soma := func(saldo int, value int) int {
	// 	return saldo + value
	// }(contaTarget.Saldo, value)

	// contaTarget.Saldo = soma
	// err = ts.rc.Save(c, *contaTarget)

	// if err != nil {
	// 	log.Fatal("Não foi possível atualizar o saldo.", err)
	// 	tx.Rollback()

	// 	return repository.Operacao{}, errors.New("Erro ao salvar operação!")
	// }

	operacao := repository.Operacao{
		Type:        0,
		ContaSource: sourceAccount,
		ContaTarget: contaTarget,
		// ContaSourceID: &sourceAccount.ID,
		// ContaTargetID: &contaTarget.ID,
		Value: value,
	}

	err = ts.ro.Save(c, &operacao)

	if err != nil {
		log.Fatal("Não foi possível finalizar a operação.", err)
		tx.Rollback()

		return repository.Operacao{}, errors.New("Erro ao salvar operação!")
	}
	tx.Commit()
	// log.Printf("Conta target %d", soma)

	return operacao, err
}
