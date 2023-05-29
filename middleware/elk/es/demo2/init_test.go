package demo2

import (
	"fmt"
	"testing"

	"github.com/elastic/go-elasticsearch"
)

// 初始化es客户端

var client *elasticsearch.Client

func TestMain(t *testing.M) {
	cfg := NewConfig()
	if cfg.Address == "" {
		panic("address is empty")
	}

	esCfg := elasticsearch.Config{
		Addresses: []string{
			cfg.Address,
		},
	}
	var err error
	client, err = elasticsearch.NewClient(esCfg)
	if err != nil {
		panic(fmt.Errorf("error creating the client: %s", err))
	}
	t.Run()
}
