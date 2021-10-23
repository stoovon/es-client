package main

import (
	"log"
	"strings"

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
}

func getClusterInfo(es *client.EsClient) {
	if err := es.ClusterInfo(); err != nil {
		log.Fatalf("Error getting the cluster info: %s", err)
	}

	printSeparator("~")
}
