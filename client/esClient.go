package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/stoovon/es-client/externalModels"
	"github.com/stoovon/es-client/models"

	"github.com/elastic/go-elasticsearch/v8"
)

type EsClient struct {
	// Allows us to add extra logic around the lifecycle of the EsClient
	client *elasticsearch.Client
}

func NewESClient() (*EsClient, error) {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return nil, err
	}

	return &EsClient{
		client: es,
	}, nil
}

func (e *EsClient) ClusterInfo() error {
	var r map[string]interface{}

	res, err := e.client.Info()
	if err != nil {
		return fmt.Errorf("error getting cluster info response: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()

		if err != nil {
			err = fmt.Errorf("error closing cluster info response: %w", err)
		}
	}(res.Body)

	if res.IsError() {
		return fmt.Errorf("cluster info response error: %s", res.String())
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return fmt.Errorf("error parsing cluster info response body: %w", err)
	}

	log.Printf("Client: %s", elasticsearch.Version)
	log.Printf("Server: %s", r["version"].(map[string]interface{})["number"])

	return err
}

func (e *EsClient) IndexPayment(payment *models.Payment) error {
	doc, err := json.Marshal(payment)
	if err != nil {
		return fmt.Errorf("unable to marshal payment: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      "payments",
		DocumentID: payment.Id.String(),
		Body:       bytes.NewReader(doc),
		Refresh:    "true", // Force payment to appear in index; use for demo, but not for production
	}

	res, err := req.Do(context.Background(), e.client)
	if err != nil {
		return fmt.Errorf("error getting response to index request: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()

		err = fmt.Errorf("error closing payment index response: %w", err)
	}(res.Body)

	if res.IsError() {
		return fmt.Errorf("error indexing payment ID: %s", payment.Id.String())
	}

	var result externalModels.IndexResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return fmt.Errorf("error parsing the payment index response body: %w", err)
	}

	log.Printf("[%s] %s; version=%d", res.Status(), result.Result, result.Version)

	return err
}
