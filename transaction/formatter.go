package transaction

import "time"

type TransactionFormatter struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Amount     int       `json:"amount"`
	Created_at time.Time `json:"created_at"`
}
type TransactionUserFormatter struct {
	Id          int              `json:"id"`
	Name        string           `json:"name"`
	Payment_url string           `json:"payment_uri"`
	Amount      int              `json:"amount"`
	Status      string           `json:"status"`
	Created_at  time.Time        `json:"created_at"`
	Project     ProjectFormatter `json:"project"`
}
type ProjectFormatter struct {
	Name      string `json:"name"`
	Image_url string `json:"image_url"`
}

func Formatter(transaction Transaction) TransactionFormatter {

	transactionFormat := TransactionFormatter{
		Id:         transaction.Id,
		Name:       transaction.User.Name,
		Amount:     transaction.Amount,
		Created_at: transaction.Created_date,
	}
	return transactionFormat
}
func FormatterUser(transaction Transaction) TransactionUserFormatter {
	status := "pending"
	if transaction.Status == 1 {
		status = "paid"
	}

	transactionFormat := TransactionUserFormatter{
		Id:          transaction.Id,
		Name:        transaction.User.Name,
		Payment_url: transaction.Payment_url,
		Amount:      transaction.Amount,
		Status:      status,
		Created_at:  transaction.Created_date,
		Project: ProjectFormatter{
			Name: transaction.Project.Project_name,
		},
	}
	if len(transaction.Project.Images) > 0 {
		transactionFormat.Project.Image_url = transaction.Project.Images[0].Image
	}

	return transactionFormat
}
func FormatterAll(transaction []Transaction) []TransactionFormatter {
	var transactions []TransactionFormatter
	if len(transaction) == 0 {
		return transactions
	}
	for _, v := range transaction {
		transactionFormat := Formatter(v)
		transactions = append(transactions, transactionFormat)

	}
	return transactions
}
func FormatterUserAll(transaction []Transaction) []TransactionUserFormatter {
	var transactions []TransactionUserFormatter
	if len(transaction) == 0 {
		return transactions
	}
	for _, v := range transaction {
		transactionFormat := FormatterUser(v)
		transactions = append(transactions, transactionFormat)

	}
	return transactions
}
