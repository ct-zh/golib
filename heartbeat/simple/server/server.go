package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

// 长连接心跳机制、服务端
// 1. server启动监听服务，等待tcp连接；
// 2. 开启协程每隔15秒往CMap里的`cs.Wch`chan循环写`Req #push`信息(测试用)；
// 3. 建立连接后等待`Req_REGISTER`消息,收到后回复`Res_REGISTER`,并将`#`后面的内容作为uid新建一个cs写入CMap；
// 消息格式 []byte([信号]#[消息])
const (
	ReqRegister byte = 1 // client向server请求注册
	ResRegister byte = 2 // server回应register消息

	ReqHeartbeat byte = 3 // 心跳包
	ResHeartbeat byte = 4 // 回复心跳包

	Req byte = 5
	Res byte = 6
)

type CS struct {
	Rch chan []byte // 读chan
	Wch chan []byte // 写chan
	Dch chan bool   // delete chan
	u   string      // uid
}

func NewCs(uid string) *CS {
	return &CS{Rch: make(chan []byte), Wch: make(chan []byte), u: uid}
}

var CMap map[string]*CS

const (
	HOST string = "127.0.0.1"
	PORT int    = 12666
)

func main() {
	CMap = make(map[string]*CS)

	// 启动监听服务
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP: net.ParseIP(HOST), Port: PORT})
	if err != nil {
		log.Println("监听端口失败，error: ", err)
		return
	}

	log.Println("初始化成功，正在等待连接")
	go PushGRT() // 协程每隔15秒循环推送消息到每隔cs的Wch chan
	Server(listen)
}

func PushGRT() {
	// 推送Req#push到cs.Wch
	for {
		time.Sleep(15 * time.Second)
		for s, cs := range CMap {
			fmt.Println("push msg to user: ", s)
			cs.Wch <- []byte{Req, '#', 'p', 'u', 's', 'h', '!'}
		}
	}
}

func Server(listen *net.TCPListener) {
	for {
		// 等待接收连接
		conn, err := listen.AcceptTCP()
		if err != nil {
			fmt.Println("接受客户端连接异常:", err.Error())
			continue
		}
		fmt.Println("客户端连接来自:", conn.RemoteAddr().String())
		go Handler(conn)
	}
}

func Handler(conn net.Conn) {
	defer conn.Close()
	data := make([]byte, 128)
	var uid string
	var C *CS

	for {
		conn.Read(data)
		fmt.Println("客户端发来数据:", string(data))
		if data[0] == ReqRegister { // register消息
			conn.Write([]byte{ResRegister, '#', 'o', 'k'})
			uid = string(data[2:])
			C = NewCs(uid)
			CMap[uid] = C
			break
		} else {
			conn.Write([]byte{ResRegister, '#', 'e', 'r'})
		}
	}

	go WHandler(conn, C)
	go RHandler(conn, C)

	go Work(C)
	select {
	case <-C.Dch:
		fmt.Println("close handler goroutine")
	}
}

// 正常写数据
// 定时检测 conn die => goroutine die
func WHandler(conn net.Conn, C *CS) {
	// 读取业务Work 写入Wch的数据
	ticker := time.NewTicker(20 * time.Second)
	for {
		select {
		case d := <-C.Wch:
			conn.Write(d)
		case <-ticker.C:
			if _, ok := CMap[C.u]; !ok {
				fmt.Println("conn die, close WHandler")
				return
			}
		}
	}
}

// 读客户端数据 + 心跳检测
func RHandler(conn net.Conn, C *CS) {
	// 心跳ack
	// 业务数据 写入Wch
	for {
		data := make([]byte, 128)
		// setReadTimeout
		err := conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		if err != nil {
			fmt.Println(err)
		}
		if _, derr := conn.Read(data); derr == nil {
			fmt.Println(data)   // 读消息
			if data[0] == Res { // 收到client的ack
				fmt.Println("recv client data ack")
			} else if data[0] == Req { // 收到client的data，打印
				fmt.Println("recv client data")
				fmt.Println(data)
				conn.Write([]byte{Res, '#'}) // 回一个ack消息
				// C.Rch <- data
			}

			continue
		}

		// 发送心跳包
		conn.Write([]byte{ReqHeartbeat, '#'})
		fmt.Println("send ht packet")
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		if _, herr := conn.Read(data); herr == nil {
			// fmt.Println(string(data))
			fmt.Println("resv ht packet ack") // 收到心跳回复
		} else {
			delete(CMap, C.u) // 没有收到心跳回复，删除用户，直接返回
			fmt.Println("delete user!")
			return
		}
	}
}

func Work(C *CS) {
	time.Sleep(5 * time.Second)
	C.Wch <- []byte{Req, '#', 'h', 'e', 'l', 'l', 'o'}

	time.Sleep(15 * time.Second)
	C.Wch <- []byte{Req, '#', 'h', 'e', 'l', 'l', 'o'}
	// 从读ch读信息
	/*	ticker := time.NewTicker(20 * time.Second)
		for {
			select {
			case d := <-C.Rch:
				C.Wch <- d
			case <-ticker.C:
				if _, ok := CMap[C.u]; !ok {
					return
				}
			}

		}
	*/ // 往写ch写信息
}
