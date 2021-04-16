package main

import (
	"github.com/ct-zh/golib/middleware/protobuf/rpc/msg"
	"log"
	"net/rpc"
)

func main() {
	client, err := rpc.Dial("tcp", "localhost:14444")
	if err != nil {
		log.Fatal(err)
	}

	req := msg.MsgRequest{Name: "李四"}
	res := msg.MsgResponse{}

	err = client.Call("SayService.Say", &req, &res)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res.Msg)
}
