package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/wurkhappy/WH-Agreements/models"
	"net/http"
)

func UpdatePayment(params map[string]interface{}, body []byte) ([]byte, error, int) {
	versionID := params["versionID"].(string)
	paymentID := params["paymentID"].(string)
	agreement, err := models.FindAgreementByVersionID(versionID)
	if err != nil {
		return nil, fmt.Errorf("Could not find agreement"), http.StatusBadRequest
	}
	payment := agreement.Payments.GetPayment(paymentID)
	var newPayment *models.Payment
	json.Unmarshal(body, &newPayment)
	payment.PaymentItems = newPayment.PaymentItems
	payment.UpdateAmountDue()
	agreement.Save()
	j, _ := json.Marshal(payment)
	return j, nil, http.StatusOK
}
