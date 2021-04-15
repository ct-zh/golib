package zk

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

type Manager struct {
	hosts      []string
	conn       *zk.Conn
	pathPrefix string
}

func NewManager(hosts []string, pathPrefix string) *Manager {
	return &Manager{hosts: hosts, pathPrefix: pathPrefix}
}

//连接zk服务器
func (z *Manager) GetConnect() error {
	conn, _, err := zk.Connect(z.hosts, 5*time.Second)
	if err != nil {
		return err
	}
	z.conn = conn
	return nil
}

//关闭服务
func (z *Manager) Close() {
	z.conn.Close()
	return
}

//获取配置
func (z *Manager) GetPathData(nodePath string) ([]byte, *zk.Stat, error) {
	return z.conn.Get(nodePath)
}

//更新配置
func (z *Manager) SetPathData(nodePath string, config []byte, version int32) (err error) {
	ex, _, _ := z.conn.Exists(nodePath)
	if !ex {
		z.conn.Create(nodePath, config, 0, zk.WorldACL(zk.PermAll))
		return nil
	}
	_, dStat, err := z.GetPathData(nodePath)
	if err != nil {
		return
	}
	_, err = z.conn.Set(nodePath, config, dStat.Version)
	if err != nil {
		fmt.Println("Update node error", err)
		return err
	}
	fmt.Println("SetData ok")
	return
}

//创建临时节点
func (z *Manager) RegisterServerPath(nodePath, host string) (err error) {
	ex, _, err := z.conn.Exists(nodePath)
	if err != nil {
		fmt.Println("Exists error", nodePath)
		return err
	}
	if !ex {
		//持久化节点，思考题：如果不是持久化节点会怎么样？
		_, err = z.conn.Create(nodePath, nil, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			fmt.Println("Create error", nodePath)
			return err
		}
	}
	//临时节点
	subNodePath := nodePath + "/" + host
	ex, _, err = z.conn.Exists(subNodePath)
	if err != nil {
		fmt.Println("Exists error", subNodePath)
		return err
	}
	if !ex {
		_, err = z.conn.Create(subNodePath, nil, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
		if err != nil {
			fmt.Println("Create error", subNodePath)
			return err
		}
	}
	return
}

//获取服务列表
func (z *Manager) GetServerListByPath(path string) (list []string, err error) {
	list, _, err = z.conn.Children(path)
	return
}

//watch机制，服务器有断开或者重连，收到消息
func (z *Manager) WatchServerListByPath(path string) (chan []string, chan error) {
	conn := z.conn
	snapshots := make(chan []string)
	errors := make(chan error)
	go func() {
		for {
			snapshot, _, events, err := conn.ChildrenW(path)
			if err != nil {
				errors <- err
			}
			snapshots <- snapshot
			select {
			case evt := <-events:
				if evt.Err != nil {
					errors <- evt.Err
				}
				fmt.Printf("ChildrenW Event Path:%v, Type:%v\n", evt.Path, evt.Type)
			}
		}
	}()

	return snapshots, errors
}

//watch机制，监听节点值变化
func (z *Manager) WatchPathData(nodePath string) (chan []byte, chan error) {
	conn := z.conn
	snapshots := make(chan []byte)
	errors := make(chan error)

	go func() {
		for {
			dataBuf, _, events, err := conn.GetW(nodePath)
			if err != nil {
				errors <- err
				return
			}
			snapshots <- dataBuf
			select {
			case evt := <-events:
				if evt.Err != nil {
					errors <- evt.Err
					return
				}
				fmt.Printf("GetW Event Path:%v, Type:%v\n", evt.Path, evt.Type)
			}
		}
	}()
	return snapshots, errors
}
