package main

import (
	manager "github.com/ct-zh/golib/middleware/zookeeper/zk"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// 配合register与write使用； zk实时获取当前服务列表;实时获取当前配置列表

var hosts = []string{"127.0.0.1:2181", "127.0.0.1:2182", "127.0.0.1:2183"}

func main() {
	man := manager.NewManager(hosts, "")
	man.GetConnect()
	defer man.Close()

	getService(man)
	getConfig(man)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func getService(man *manager.Manager) {
	zlist, err := man.GetServerListByPath("/service")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("list: %+v", zlist)

	ch, chanErr := man.WatchServerListByPath("/service")
	go func() {
		for {
			select {
			case e := <-chanErr:
				log.Println("chan err: ", e)
			case l := <-ch:
				log.Printf("watch list: %+v", l)
			}
		}
	}()
}

func getConfig(man *manager.Manager) {
	_, _, err := man.GetPathData("/rs_server_conf")
	if err != nil {
		log.Fatal(err)
	}

	dataChan, chanErr := man.WatchPathData("/rs_server_conf")
	go func() {
		select {
		case e := <-chanErr:
			log.Println("chan2 err: ", e)
		case d := <-dataChan:
			log.Printf("watch config: %+v", string(d))
		}
	}()
}
