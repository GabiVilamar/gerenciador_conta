package handler

import (
	"gerenciador_conta/repository"
	"gerenciador_conta/service"
)

var Conta *contaHandler

func init() {
	rc := repository.NewConta()
	ro := repository.NewOperation()
	ts := service.NewTransactionService(rc, ro)
	Conta = &contaHandler{rc, ro, ts}
}
