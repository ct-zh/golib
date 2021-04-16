package main

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
)

// 获取当前所有服务

func main() {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}

	services, err := client.Agent().Services()
	if err != nil {
		panic(err)
	}

	for _, service := range services {
		s := service
		data, err := json.Marshal(s)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v \n", string(data))
	}
}
