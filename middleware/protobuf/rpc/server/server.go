package main

import (
	"github.com/ct-zh/golib/middleware/protobuf/rpc/msg"
	"log"
	"net"
	"net/rpc"
)

type SayService struct {
}

func (s *SayService) Say(req *msg.MsgRequest, reply *msg.MsgResponse) error {
	log.Println("req: ", req.Name)
	reply.Msg = "hello: " + req.Name
	return nil
}

func main() {
	_ = rpc.RegisterName("SayService", new(SayService))

	listener, err := net.Listen("tcp", ":14444")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("conn error ", err.Error())
			continue
		}

		rpc.ServeConn(conn)
	}
}
