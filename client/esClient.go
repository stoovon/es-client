package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

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
