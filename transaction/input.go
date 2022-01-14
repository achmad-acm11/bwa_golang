package transaction

type GetProjectIdUri struct {
	Id              int `uri:"id" binding:"required"`
	Current_id_user int
}
type CreateTransaction struct {
	Id_user    int
	Id_project int `json:"id_project" binding:"required"`
	Amount     int `json:"amount" binding:"required"`
}

type ResponsePaymentInput struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	PaymentType       string `json:"payment_type"`
	FraudStatus       string `json:"fraud_status"`
}
