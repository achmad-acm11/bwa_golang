package handler

import (
	"bwa_golang/helper"
	"bwa_golang/transaction"
	"bwa_golang/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	transactionService transaction.Service
}

func NewTransactionHandler(t transaction.Service) *transactionHandler {
	return &transactionHandler{t}
}

// Handler for Get List Transaction Id_project
func (t *transactionHandler) GetTransactionByProjectId(c *gin.Context) {
	// Get Id_project URI
	var uri_id transaction.GetProjectIdUri
	err := c.ShouldBindUri(&uri_id)

	if err != nil {
		response := helper.APIResponse("Internal Server Error", "error", http.StatusBadRequest, gin.H{"errors": err.Error()})
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// Get Id_user Login
	currentUser := c.MustGet("currentUser").(user.User)
	uri_id.Current_id_user = currentUser.Id

	// Get List Transaction by Project Id service
	dataTrasaction, err := t.transactionService.GetTransactionByProjectId(uri_id)
	if err != nil {
		response := helper.APIResponse("Internal Service Error", "error", http.StatusBadRequest, gin.H{"errors": err.Error()})
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Get Data Transaction Success", "success", http.StatusOK, transaction.FormatterAll(dataTrasaction))
	c.JSON(http.StatusOK, response)
}

// Handler for Get List Transaction Id_user
func (t *transactionHandler) GetTransactionByUserId(c *gin.Context) {
	// Get Id_user Login
	currentUser := c.MustGet("currentUser").(user.User)
	// Get Data Transaction User Service
	dataTrasaction, err := t.transactionService.GetTransactionByUserId(currentUser.Id)
	if err != nil {
		response := helper.APIResponse("Internal Service Error", "error", http.StatusBadRequest, gin.H{"errors": err.Error()})
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Get Data Transaction Success", "success", http.StatusOK, transaction.FormatterUserAll(dataTrasaction))
	c.JSON(http.StatusOK, response)
}

// Handler for Create Transaction Project
func (t *transactionHandler) CreateTransaction(c *gin.Context) {
	// Format Input Data Project & Validation
	var input transaction.CreateTransaction

	err := c.ShouldBindJSON(&input)
	if err != nil {
		response := helper.APIResponse("Internal Server Error", "error", http.StatusBadRequest, helper.FormatValidationError(err))
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// Get Id_user Login
	currentUser := c.MustGet("currentUser").(user.User)
	input.Id_user = currentUser.Id
	// Create Transaction Service
	transaction, err := t.transactionService.CreateTransaction(input, currentUser)
	if err != nil {
		response := helper.APIResponse("Internal Server Error", "error", http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}
	transaction.User = currentUser
	response := helper.APIResponse("Success Create Transaction", "success", http.StatusOK, transaction)
	c.JSON(http.StatusOK, response)
}

// Handler for Response Payment
func (t *transactionHandler) ResponsePayment(c *gin.Context) {
	var inputResponsePayment transaction.ResponsePaymentInput
	err := c.ShouldBindJSON(&inputResponsePayment)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}
	err = t.transactionService.ResponsePayment(inputResponsePayment)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}
	c.JSON(http.StatusOK, "Success Payment")
}
