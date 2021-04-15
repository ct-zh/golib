package main

import (
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"time"
)

// zookeeper的基本curd操作

var hosts = []string{"127.0.0.1:2181", "127.0.0.1:2182", "127.0.0.1:2183"}

func main() {
	conn, _, err := zk.Connect(hosts, time.Second*5)
	if err != nil {
		panic(err)
	}

	path := "/test_tree"

	// 创建节点 flags: 节点类型； acl： 节点权限;
	str, err := conn.Create(path, []byte("tree"), 0, zk.WorldACL(zk.PermAll))
	if err != nil {
		panic(err)
	}
	log.Printf("root %+v", str)

	// 查询节点数据
	data, dStat, err := conn.Get(path)
	if err != nil {
		panic(err)
	}
	log.Printf("第一次查询： data: %v dStat: %+v", string(data), dStat)

	// 修改, 注意需要先查询，拿到版本号dStat.Version
	if _, err := conn.Set(path, []byte("new_content"), dStat.Version); err != nil {
		log.Println("update err", err)
	}

	// 删除,删除也需要版本号
	data, dStat, _ = conn.Get(path)
	log.Printf("修改后的数据： data: %v dStat: %+v", string(data), dStat)
	if err := conn.Delete(path, dStat.Version); err != nil {
		log.Println("Delete err", err)
		//return
	}

	// 设置子节点/test_tree/subnode, 子节点的父节点必须存在;
	if _, err := conn.Create(path, []byte("tree_content"), 0, zk.WorldACL(zk.PermAll)); err != nil {
		log.Println("create err", err)
	}
	if _, err := conn.Create(path+"/subnode", []byte("node_content"),
		0, zk.WorldACL(zk.PermAll)); err != nil {
		log.Println("create err", err)
	}

	// 获取子节点列表
	childNodes, _, err := conn.Children(path)
	if err != nil {
		log.Println("Children err", err)
	}
	log.Println("childNodes", childNodes)

	// 删除，不能删除存在子结点的父节点，必须先删除所有子结点
	_, dStat, _ = conn.Get(path + "/subnode")
	err = conn.Delete(path+"/subnode", dStat.Version)
	if err != nil {
		panic(err)
	}

	_, dStat, _ = conn.Get(path)
	err = conn.Delete(path, dStat.Version)
	if err != nil {
		panic(err)
	}
}
