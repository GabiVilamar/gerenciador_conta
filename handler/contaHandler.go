package handler

import (
	"gerenciador_conta/repository"
	"gerenciador_conta/service"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type contaHandler struct {
	rc repository.Conta
	ro repository.Operation
	ts service.TransactionService
}

func (ch *contaHandler) Create(c echo.Context) error {
	conta := new(repository.ContaModel)

	err := c.Bind(&conta)

	if err != nil {
		log.Fatal("Erro ao criar conta.", err)
		return c.String(http.StatusBadRequest, "bad request")
	}

	err = ch.rc.Save(c, *conta)
	if err != nil {
		log.Fatal("Não foi possível criar esse registro.", err)
		return c.String(http.StatusBadRequest, "bad request")
	}

	return c.NoContent(http.StatusNoContent)

}

func (ch *contaHandler) GetByID(c echo.Context) error {
	var id string
	id = c.Param("id")

	formattedID, err := strconv.Atoi(id)

	if err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusNotFound, data)
	}

	response, err := ch.rc.Get(formattedID)

	if err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusNotFound, data)
	}

	return c.JSON(http.StatusOK, response)
}

func (ch *contaHandler) GetAll(c echo.Context) error {
	contas, err := ch.rc.GetAll()

	if err != nil {
		log.Fatal("Erro ao buscar contas", err)
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, contas)

}

func (ch *contaHandler) Remove(c echo.Context) error {
	id := c.Param("id")
	formattedID, err := strconv.Atoi(id)

	if err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	_, err = ch.rc.Get(formattedID)

	if err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusNotFound, data)
	}

	err = ch.rc.Remove(formattedID)

	if err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	response := map[string]interface{}{
		"message": "A conta foi removida com sucesso!",
	}
	return c.JSON(http.StatusOK, response)

}

func (ch *contaHandler) Update(c echo.Context) error {
	var id string
	id = c.Param("id")
	formattedID, err := strconv.Atoi(id)

	if err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	conta, err := ch.rc.Get(formattedID)

	if err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusNotFound, data)

	}

	err = c.Bind(&conta)

	if err != nil {
		log.Fatal("Erro ao atualizar conta.", err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err = ch.rc.Save(c, *conta)

	if err != nil {
		log.Fatal("Não foi possível atualizar esse registro.", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (ch *contaHandler) Transfer(c echo.Context) error {

	id := c.Param("id")
	contaSourceID, err := strconv.Atoi(id)

	if err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	payload := new(DataTranferObject)
	err = c.Bind(payload)

	if err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusBadRequest, data)

	}

	contaSource, err := ch.rc.Get(int(contaSourceID))
	if err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusNotFound, data)
	}

	log.Print("Conta source", contaSource)

	// Início da operacao de transferência

	switch payload.Type {
	case 0:

		extract, err := ch.ts.CreateTranfer(c, payload.ContaTarget, contaSource, payload.Value)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		return c.JSON(http.StatusOK, extract)

	case 1:
		extract, err := ch.ts.CreateDeposit(c, payload.ContaTarget, payload.Value)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		return c.JSON(http.StatusOK, extract)

	case 2:
		extract, err := ch.ts.CreateWithdrawal(c, payload.ContaTarget, payload.Value)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		return c.JSON(http.StatusOK, extract)
	}

	return c.JSON(http.StatusBadRequest, "Operação não permitida!")
}

func (ch *contaHandler) GetAllTransaction(c echo.Context) error {
	transactions, err := ch.ro.GetAllTransaction()

	if err != nil {
		log.Fatal("Erro ao buscar transações", err)
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, transactions)

}
