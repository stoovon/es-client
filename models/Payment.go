package models

import (
	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
	"github.com/stoovon/es-client/externalModels"
)

type Payment struct {
	Amount      money.Money
	Beneficiary externalModels.Party
	Debtor      externalModels.Party
	Id          uuid.UUID
	Version     int
}
