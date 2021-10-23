package strategy

import (
	"github.com/Rhymond/go-money"
	"github.com/stoovon/es-client/models"
)

func UpdatePaymentStrategy(existingPayment *models.Payment) *models.Payment {
	result := existingPayment

	// A strategy will typically apply data from a message queue to Payment.
	result.Amount = *money.New(1000, money.GBP)
	result.Version = existingPayment.Version + 1

	return result
}