package helper

import (
	"github.com/go-playground/validator/v10"
	"github.com/leekchan/accounting"
)

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Meta struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Code    int    `json:"code"`
}

func APIResponse(message string, status string, code int, data interface{}) Response {
	meta := Meta{
		Message: message,
		Status:  status,
		Code:    code,
	}
	response := Response{
		Meta: meta,
		Data: data,
	}
	return response
}

func FormatValidationError(err error) (errors []string) {

	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Error())
	}
	return
}

func FormatIDR(number int) string {
	moneyFormat := accounting.Accounting{Symbol: "Rp", Precision: 0, Thousand: ".", Decimal: ","}
	return moneyFormat.FormatMoney(number)
}
