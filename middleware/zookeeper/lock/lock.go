package lock

import (
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"time"
)

// 分布式锁

// 利用 zookeeper 的同级节点的唯一性特性，在需要获取排他锁时，
// 所有的客户端试图通过调用 create() 接口，
// 在 /exclusive_lock节点下创建临时子节点/exclusive_lock/lock，
// 最终只有一个客户端能创建成功，那么此客户端就获得了分布式锁。
// 同时，所有没有获取到锁的客户端可以在/exclusive_lock节点上
// 注册一个子节点变更的 watcher 监听事件，以便重新争取获得锁。

// 定义锁
const LockPath = "/exclusive_lock/lock"

var hosts = []string{"127.0.0.1:2181", "127.0.0.1:2182", "127.0.0.1:2183"}

func zkLock() {
	conn, _, err := zk.Connect(hosts, time.Second)
	if err != nil {
		panic(err)
	}

	lock := zk.NewLock(conn, LockPath, zk.WorldACL(zk.PermAll))
	err = lock.Lock()

	// 已经有锁了
	if err != nil {
		panic(err)
	}

	log.Println("i get lock ")
	time.Sleep(time.Second * 10)

	defer lock.Unlock()
	log.Println("unlock !")
}
