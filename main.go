package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
	"github.com/stoovon/es-client/externalModels"
	"github.com/stoovon/es-client/models"

	"github.com/stoovon/es-client/client"
)

func printSeparator(sep string) {
	log.Println(strings.Repeat(sep, 37))
}

func main() {
	log.SetFlags(0)

	var es *client.EsClient
	var err error

	if es, err = client.NewESClient(); err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	getClusterInfo(es)
	createPayment(es)
	findPayments(es)
}

func getClusterInfo(es *client.EsClient) {
	if err := es.ClusterInfo(); err != nil {
		log.Fatalf("Error getting the cluster info: %s", err)
	}

	printSeparator("~")
}

func createPayment(es *client.EsClient) {
	var payment *models.Payment
	var payment2 *models.Payment
	var err error

	if payment, err = newPayment(); err != nil {
		log.Fatalf("Unable to instantiate new payment: %v", err)
	}

	if payment2, err = newPayment(); err != nil {
		log.Fatalf("Unable to instantiate new payment: %v", err)
	}

	if err = es.IndexPayments([]*models.Payment{payment, payment2}...); err != nil {
		log.Fatalf("Unable to index new payment: %v", err)
	}

	printSeparator("-")
}

func newPayment() (*models.Payment, error) {
	var id uuid.UUID
	var err error

	if id, err = uuid.NewRandom(); err != nil {
		return nil, fmt.Errorf("unable to create payment Id: %w", err)
	}

	return &models.Payment{
		Amount: *money.New(1337, money.GBP),
		Beneficiary: externalModels.Party{
			Name:          "Ada Lovelace",
			AccountNumber: "12345678",
			BankId:        "111111",
		},
		Debtor: externalModels.Party{
			Name:          "Charles Babbage",
			AccountNumber: "12345678",
			BankId:        "222222",
		},
		Id: id,
	}, nil
}

func findPayments(es *client.EsClient) {
	if err := es.FindPayments(); err != nil {
		log.Fatalf("Unable to find payments: %v", err)
	}
	printSeparator("=")
}
