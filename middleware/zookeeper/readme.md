# zookeeper
## 安装、配置
> 参考官方文档: http://zookeeper.apache.org/doc/r3.6.0/zookeeperStarted.html
建议使用docker: `docker pull zookeeper`

## 操作
进入cli命令行: `zkCli.sh`, 使用`man`命令查看帮助;



## 节点类型
- 持久节点
- 临时节点: 会话失效,节点自动清理;
- 顺序节点: 节点创建, 自动分配序列号;

## go zookeeper
> 依赖包: github.com/samuel/go-zookeeper/zk

