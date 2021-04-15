package main

import (
	"log"
	"net"
)

// 长连接心跳机制，客户端
//
const (
	ReqRegister byte = 1 // 1 --- c register cid
	ResRegister byte = 2 // 2 --- s response

	ReqHeartbeat byte = 3 // 3 --- s send heartbeat req
	ResHeartbeat byte = 4 // 4 --- c send heartbeat res

	Req byte = 5 // 5 --- cs send data
	Res byte = 6 // 6 --- cs send ack
)

var Dch chan bool
var Rch chan []byte
var Wch chan []byte

func main() {
	Dch = make(chan bool)
	Rch = make(chan []byte)
	Wch = make(chan []byte)
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:12666")
	conn, err := net.DialTCP("tcp", nil, addr)
	//	conn, err := net.Dial("tcp", "127.0.0.1:6666")
	if err != nil {
		log.Println("连接服务端失败:", err.Error())
		return
	}
	log.Println("已连接服务器")
	defer conn.Close()
	go Handler(conn)
	select {
	case <-Dch:
		log.Println("关闭连接")
	}
}

func Handler(conn *net.TCPConn) {
	// 直到register ok
	data := make([]byte, 128)
	for {
		conn.Write([]byte{ReqRegister, '#', '2'})
		conn.Read(data)
		//		fmt.Println(string(data))
		if data[0] == ResRegister {
			break
		}
	}
	//	fmt.Println("i'm register")
	go RHandler(conn)
	go WHandler(conn)
	go Work()
}

// 循环读消息，直到Dch收到数据
func RHandler(conn *net.TCPConn) {
	for {
		// 心跳包,回复ack
		data := make([]byte, 128)
		i, _ := conn.Read(data)
		if i == 0 {
			Dch <- true
			return
		}
		if data[0] == ReqHeartbeat { // 心跳包
			log.Println("recv ht pack")
			conn.Write([]byte{ResRegister, '#', 'h'})
			log.Println("send ht pack ack")
		} else if data[0] == Req {
			log.Println("recv data pack")
			log.Printf("%v\n", string(data[2:]))
			Rch <- data[2:]
			conn.Write([]byte{Res, '#'})
		}
	}
}

func WHandler(conn net.Conn) {
	for {
		select {
		case msg := <-Wch:
			log.Println((msg[0]))
			log.Println("send data after: " + string(msg[1:]))
			conn.Write(msg)
		}
	}

}

func Work() {
	for {
		select {
		case msg := <-Rch:
			log.Println("work recv " + string(msg))
			Wch <- []byte{Req, '#', 'x', 'x', 'x', 'x', 'x'}
		}
	}
}
