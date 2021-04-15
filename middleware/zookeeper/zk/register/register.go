package main

import (
	manager "github.com/ct-zh/golib/middleware/zookeeper/zk"
	"log"
	"strconv"
	"time"
)

// zookeeper实现服务注册原理，创建临时节点，会话结束自动删除;

var hosts = []string{"127.0.0.1:2181", "127.0.0.1:2182", "127.0.0.1:2183"}

func main() {
	man := manager.NewManager(hosts, "")
	man.GetConnect()
	defer man.Close()

	i := 0
	for {
		man.RegisterServerPath("/service", strconv.Itoa(i))
		log.Println("zookeeper register: ", i)
		time.Sleep(time.Second * 5)
		i++
	}
}
