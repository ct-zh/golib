package main

import (
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch"
)

// elasticsearch的基本操作, 通过go-elasticsearch包
func main() {
	cfg := NewConfig()

	// step1. 连接es
	es := getClient(cfg)

	// 成功打印es的信息，代表连接成功
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	fmt.Println(res)
}

func getClient(cfg *Config) *elasticsearch.Client {
	esCfg := elasticsearch.Config{
		Addresses: []string{
			cfg.Address,
		},
	}
	es, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	return es
}
