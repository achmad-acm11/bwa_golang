package handler

import (
	"bwa_golang/transaction"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	transactionService transaction.Service
}

func NewTransactionHandler(t transaction.Service) *transactionHandler {
	return &transactionHandler{t}
}

func (t *transactionHandler) Index(c *gin.Context) {
	transactions, err := t.transactionService.GetAllTransaction()

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	c.HTML(http.StatusOK, "index_transaction.html", gin.H{"transactions": transactions})
}
