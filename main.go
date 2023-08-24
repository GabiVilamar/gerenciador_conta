package main

import (
	"fmt"
	"gerenciador_conta/handler"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	// g := e.Group("meu_banco")
	e.GET("/", hello)
	e.GET("/contas", handler.Conta.GetAll)
	e.GET("/contas/:id", handler.Conta.GetByID)
	e.POST("/contas", handler.Conta.Create)
	e.PUT("/contas/:id", handler.Conta.Update)
	e.DELETE("/contas/:id", handler.Conta.Remove)
	e.POST("/contas/:id/transaction", handler.Conta.Transfer)
	e.GET("/contas/transactions", handler.Conta.GetAllTransaction)
	// e.GET("/healthcheck", handler)

	// /contas?id=<contasID>

	// Start server
	port := 1350
	addr := fmt.Sprintf(":%d", port)
	// log.Info(fmt.Sprintf("Rodando na porta %d", port))
	e.Logger.Fatal(e.Start(addr))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func saldo(c echo.Context) error {
	return c.String(http.StatusOK, "Opaa")
}

// **Bug pra corrigir**
// Quando atualizo um valor num registro da tabela, o campo created buga e fica com a data errada

// Criar estrutura de payload

// Ferramenta de hot reload pra recarregar server

// log middleware implementado do 0, sem que seja pelo echo

// criar lógica de transação do handler para o service
