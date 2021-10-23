package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/stoovon/es-client/mapper"

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

func (e *EsClient) IndexPayments(payments ...*models.Payment) error {
	var wg sync.WaitGroup

	errors := make(chan error)
	wgDone:= make(chan bool)

	for i, payment := range payments {
		wg.Add(1)

		go func(i int, payment *models.Payment) {
			defer wg.Done()

			doc, err := json.Marshal(payment)
			if err != nil {
				errors <- fmt.Errorf("error encoding payment: %w", err)
				return
			}

			req := esapi.IndexRequest{
				Index:      "payments",
				DocumentID: payment.Id.String(),
				Body:       bytes.NewReader(doc),
				Refresh:    "true", // Force payment to appear in index; use for demo, but not for production
			}

			res, err := req.Do(context.Background(), e.client)
			if err != nil {
				errors <- fmt.Errorf("error getting response to index request: %w", err)
				return
			}

			defer func(Body io.ReadCloser) {
				err = Body.Close()

				if err != nil {
					errors <- fmt.Errorf("error closing payment index response: %w", err)
				}
			}(res.Body)

			if res.IsError() {
				errors <- fmt.Errorf("error indexing payment ID: %s", payment.Id.String())
				return
			}

			var result externalModels.IndexResponse
			if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
				errors <- fmt.Errorf("error parsing the payment index response body: %w", err)
				return
			}

			log.Printf("[%s] %s; version=%d", res.Status(), result.Result, result.Version)
		}(i, payment)
	}

	go func () {
		wg.Wait()
		close(wgDone)
	}()

	select {
	case <- wgDone:
		return nil
	case err := <- errors:
		close(errors)
		return fmt.Errorf("error inserting payment: %w", err)
	}
}

func (e *EsClient) FindPayments() error {
	var result externalModels.FindResponse

	var buf bytes.Buffer
	query := externalModels.ESQuery{
		Query: externalModels.Query{
			Match: map[string]string{
				"Beneficiary.Name": "Ada Lovelace",
				//"Amount.amount": 9999,
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return fmt.Errorf("error encoding find payment query: %w", err)
	}

	res, err := e.client.Search(
		e.client.Search.WithContext(context.Background()),
		e.client.Search.WithIndex("payments"),
		e.client.Search.WithBody(&buf),
		e.client.Search.WithTrackTotalHits(true),
		e.client.Search.WithPretty(),
	)
	if err != nil {
		return fmt.Errorf("error getting find payment response: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()

		if err != nil {
			err = fmt.Errorf("error closing find payment response: %w", err)
		}
	}(res.Body)

	if res.IsError() {
		var e map[string]interface{}

		if err = json.NewDecoder(res.Body).Decode(&e); err != nil {
			return fmt.Errorf("error parsing the response body: %w", err)
		} else {
			return fmt.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return fmt.Errorf("error parsing the response body: %w", err)
	}

	log.Printf(
		"[%s] %d hit(s); took %dms",
		res.Status(),
		result.Hits.Total.Value,
		result.Took,
	)

	for _, hit := range result.Hits.Hits {
		payment, err := mapper.PaymentMapper(hit.Source)
		if err != nil {
			return fmt.Errorf("error mapping the payment: %w", err)
		}

		var renderedPayment []byte
		if renderedPayment, err = json.MarshalIndent(payment, "", "\t"); err != nil {
			return fmt.Errorf("error rendering the payment: %w", err)
		}

		log.Printf(" * ID=%s\n%s", hit.Id, renderedPayment)
	}

	return err
}
