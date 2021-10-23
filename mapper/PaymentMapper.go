package mapper

import (
	"fmt"

	"github.com/stoovon/es-client/externalModels"

	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
	"github.com/stoovon/es-client/models"
)

func moneyMapper(amount externalModels.Amount) *money.Money {
	return money.New(
		amount.Amount,
		amount.Currency,
	)
}

func PaymentMapper(payment externalModels.Payment) (*models.Payment, error) {
	id, err := uuid.Parse(payment.Id)
	if err != nil {
		return nil, fmt.Errorf("unable to parse id: %w", err)
	}

	return &models.Payment{
		Amount:      *moneyMapper(payment.Amount),
		Beneficiary: payment.Beneficiary,
		Debtor:      payment.Debtor,
		Id:          id,
	}, nil
}
