package main

import (
	"fmt"
	manager "github.com/ct-zh/golib/middleware/zookeeper/zk"
	"log"
	"time"
)

// 使用zk实现配置中心的功能

var hosts = []string{"127.0.0.1:2181", "127.0.0.1:2182", "127.0.0.1:2183"}

func main() {
	man := manager.NewManager(hosts, "")
	man.GetConnect()
	defer man.Close()
	i := 0

	for {
		conf := fmt.Sprintf("{name:" + fmt.Sprint(i) + "}")
		man.SetPathData("/rs_server_conf", []byte(conf))
		log.Println("zookeeper write: ", i)
		time.Sleep(5 * time.Second)
		i++

	}
}
