package transaction

import (
	"bwa_golang/project"
	"bwa_golang/user"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/veritrans/go-midtrans"
)

type Service interface {
	GetAllTransaction() ([]Transaction, error)
	GetTransactionByProjectId(input GetProjectIdUri) ([]Transaction, error)
	GetTransactionByUserId(id_user int) ([]Transaction, error)
	CreateTransaction(input CreateTransaction, user user.User) (Transaction, error)
	GetPaymentUrl(transaction Transaction, user user.User) (string, error)
	ResponsePayment(input ResponsePaymentInput) error
}

type service struct {
	repository        Repository
	repositoryProject project.Repository
}

func NewService(repository Repository, repositoryProject project.Repository) *service {
	return &service{repository, repositoryProject}
}

// Get All Transaction
func (s *service) GetAllTransaction() ([]Transaction, error) {
	transactions, err := s.repository.GetAll()

	if err != nil {
		return transactions, err
	}
	return transactions, nil
}

// Get List Transaction by Id_proejct service
func (s *service) GetTransactionByProjectId(input GetProjectIdUri) ([]Transaction, error) {
	// Get Data Project by Id
	dataProject, err := s.repositoryProject.GetOneById(input.Id)
	if err != nil {
		return []Transaction{}, err
	}
	// Check Project is have a User Login
	if input.Current_id_user != dataProject.Id_user {
		return []Transaction{}, errors.New("permission denied")
	}
	// Get Data List Transaction by Project
	transaction, err := s.repository.GetTransactionByProjectId(input.Id)
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

// Get List Transaction by Id_user Service
func (s *service) GetTransactionByUserId(id_user int) ([]Transaction, error) {
	// Get List Transaction User
	transaction, err := s.repository.GetTransactionByUserId(id_user)
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

// Create Transaction Service
func (s *service) CreateTransaction(input CreateTransaction, user user.User) (Transaction, error) {
	transaction := Transaction{
		Id_user:      input.Id_user,
		Id_project:   input.Id_project,
		Amount:       input.Amount,
		Created_date: time.Now(),
		Updated_date: time.Now(),
	}
	// Create Transaction
	transaction, err := s.repository.Create(transaction)
	if err != nil {
		return transaction, err
	}
	// Get Payment Url Midtrans
	paymentUrl, err := s.GetPaymentUrl(transaction, user)
	if err != nil {
		return transaction, err
	}
	// Update Transaction
	transaction.Payment_url = paymentUrl
	transaction, err = s.repository.Update(transaction)
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}

// Get Payment Url Service
func (s *service) GetPaymentUrl(t Transaction, user user.User) (string, error) {

	midclient := midtrans.NewClient()
	midclient.ServerKey = "SB-Mid-server-Jf-tUH-YfI90fwdCHbVcxlS6"
	midclient.ClientKey = "SB-Mid-client-ZCFWqbNrVydX8sZM"
	midclient.APIEnvType = midtrans.Sandbox

	// var snapGateway midtrans.SnapGateway

	snapGateway := midtrans.SnapGateway{
		Client: midclient,
	}

	snapReq := &midtrans.SnapReq{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(t.Id),
			GrossAmt: int64(t.Amount),
		},
		CustomerDetail: &midtrans.CustDetail{
			FName: user.Name,
			Email: user.Email,
		},
	}

	log.Println("GetToken:")
	snapTokenResp, err := snapGateway.GetToken(snapReq)
	if err != nil {
		return "", err
	}
	return snapTokenResp.RedirectURL, nil
}

func (s *service) ResponsePayment(input ResponsePaymentInput) error {
	transaction_id, _ := strconv.Atoi(input.OrderID)
	dataTransaction, err := s.repository.GetTransactionById(transaction_id)
	if err != nil {
		return err
	}

	if dataTransaction.Id == 0 {
		return err
	}

	if input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
		dataTransaction.Status = 1
	} else if input.TransactionStatus == "settlement" {
		dataTransaction.Status = 1
	} else if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" || input.TransactionStatus == "cancel" {
		dataTransaction.Status = 0
		return nil
	}
	dataTransaction, err = s.repository.Update(dataTransaction)

	if err != nil {
		return err
	}

	dataProject, err := s.repositoryProject.GetOneById(dataTransaction.Id_project)
	if err != nil {
		return err
	}
	dataProject.Backer_count += 1
	dataProject.Current_amount += dataTransaction.Amount
	dataProject, err = s.repositoryProject.Update(dataProject)
	if err != nil {
		return err
	}
	return nil
}
