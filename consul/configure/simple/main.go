package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

// consul实现远程配置
// consul地址是默认的127.0.0.1:8500

func main() {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}
	// 获得kv client
	kv := client.KV()

	p := &api.KVPair{
		Key:         "testkey/hahahh",
		CreateIndex: 1,
		Value:       []byte("errrrrrrrrr"),
	}
	info, err := kv.Put(p, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v \n", info)

	pair, _, err := kv.Get("testkey/hahahh", nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("KV: %v %v \n", pair.Key, string(pair.Value))
}
