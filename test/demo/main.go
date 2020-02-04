package main

import (
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v7"
)

func main() {

	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://192.168.1.102:9200",
		},
		// ...
	}
	es, err := elasticsearch.NewClient(cfg)
	fmt.Println(err)
	log.Println(elasticsearch.Version)
	_, err = es.Ping()
	log.Println(err)
}
