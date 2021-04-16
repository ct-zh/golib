# zookeeper
## 安装、配置
> 参考官方文档: http://zookeeper.apache.org/doc/r3.6.0/zookeeperStarted.html
建议使用docker: `docker pull zookeeper`

## 操作
进入cli命令行: `zkCli.sh`, 使用`man`命令查看帮助;比较常用的命令有:`ls`, `ls2`, `get`, `stat`,`create`, `set`, `delete`等

## 节点
- 持久节点
- 临时节点: 会话失效,节点自动清理;
- 顺序节点: 节点创建, 自动分配序列号;

### 节点特性
- 同一级节点 key 名称是唯一的
- 创建节点时，必须要带上全路径
- session 关闭，临时节点清除
- 自动创建顺序节点
- watch 机制，监听节点变化
- delete 命令只能一层一层删除(新版本可以使用`deleteall`递归删除)

## 数据同步流程
zk使用zab协议来实现分布式数据一致性;zab协议分为:消息广播、故障恢复;
### 消息广播
2pc 事务提交: 
1. client发起写请求,如果是连接的follower,则将请求转发给leader;
2. leader将事务请求以 Proposal (提议)广播到所有 Follower 节点;(OBSERVING节点只负责复制数据,不参与消息广播)
3. 如果集群中有过半的Follower 服务器进行正确的 ACK 反馈, 那么Leader就会再次向所有的 Follower 服务器发送commit 消息;

### 故障恢复:选举原理
几个参数:
- 服务器 ID(myid)：编号越大在选举算法中权重越大;
- 事务 ID(zxid)：值越大说明数据越新，权重越大;
- 逻辑时钟(epoch-logicalclock)：同一轮投票过程中的逻辑时钟值是相同的，每投完一次值会增加

几个状态:
- LOOKING: 竞选状态
- FOLLOWING: 随从状态，同步 leader 状态，参与投票
- OBSERVING: 观察状态，同步 leader 状态，不参与投票
- LEADING: 领导者状态

1. 初始所有机器都投自己票;
2. 收到别人的票时,先比较票是否来源于LOOKING节点、票是否与自己时钟相同;
3. 处理投票: 比较zxid,最大的作为leader; 如果有相同的zxid, 则比较myid,最大的作为leader;
4. 统计投票,只要有过半机器接受了某个机器的票,则直接当选leader;
5. 新加入的节点直接变成follower;
6. failover过程: leader挂了,其他follower状态变成Looking, 并发起投票;流程与2、3相同;

## 分布式锁
### 排他锁(写锁)

**实现方式：**

利用 zookeeper 的同级节点的唯一性特性，在需要获取排他锁时，所有的客户端试图通过调用 create() 接口，在 **/exclusive_lock** 节点下创建临时子节点 **/exclusive_lock/lock**，最终只有一个客户端能创建成功，那么此客户端就获得了分布式锁。同时，所有没有获取到锁的客户端可以在 **/exclusive_lock** 节点上注册一个子节点变更的 watcher 监听事件，以便重新争取获得锁。

### 共享锁(读锁)

**实现方式：**

1. 客户端调用 create 方法创建类似定义锁方式的临时顺序节点。
2. 客户端调用 getChildren 接口来获取所有已创建的子节点列表。
3. 判断是否获得锁，对于读请求如果所有比自己小的子节点都是读请求或者没有比自己序号小的子节点，表明已经成功获取共享锁，同时开始执行度逻辑。对于写请求，如果自己不是序号最小的子节点，那么就进入等待。
4. 如果没有获取到共享锁，读请求向比自己序号小的最后一个写请求节点注册 watcher 监听，写请求向比自己序号小的最后一个节点注册watcher 监听。

## go zookeeper
> 依赖包: github.com/samuel/go-zookeeper/zk

