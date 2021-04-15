# consul
符合：CP， 一致性，容错性；
## 安装
建议使用docker: `docker pull consul`,见dockerhub: https://hub.docker.com/_/consul

开发模式: `consul agent -dev -client=<ip>`

访问管理后台: `http://localhost:8500`代表服务启动成功
> [文档]( https://www.consul.io/docs/agent)
## 介绍
功能:
    1. 服务发现;
    2. 运行状况检查;
    3. kv存储;
    4. 安全服务通信: 建立TLS连接;
    5. 多数据中心: 比如可以在阿里云、腾讯云等分别建立集群,就可以融合集群;

### 总体流程
1. consul分为client端和server端;
2. server选举产生leader,如果client访问到follower,follower会访问到leader再回传给client.server会通过lan gossip(port:8301)检查client服务的状态;
3. client端通过rpc的方式8300访问server;client端互相通过lan gossip协议同步数据;
4. 两个数据中心的server通过wan gossip同步数据;
5. server之间通过raft协议选举;

角色分为：consul, 服务提供producer, 服务调用consumer;
1. producer-> consul, 注册服务；
2. consul-> producer 健康检查；
3. consumer从consul拉去一个临时列表,读取到需要到服务，发送请求；

### gossip协议(八卦协议)
事件发生后,节点之间会立马互相传播这个事件,直到所有节点都知道. 

这里存在两个池Lan pool与Wan pool, lanpool主要功能是:
    1. 让client自动发现server节点, 减少配置量
    2. 分布式故障检测
    3. 快速传播广播事件
wanpool:
    1. 全局唯一;
    2. 不同数据中心的server都要加入wan pool
    3. 允许服务器执行跨数据中心的请求

有两种方式: 
    1. 推送方式（Push-based）
        1.1 网络中的某个节点随机选择N个节点作为数据接收对象
        1.2 该节点向其选中的N个节点传输相应的信息
        1.3 接收到信息的节点处理它接收到的数据
        1.4 接收到数据的节点再从第一步开始重复执行

    2. 拉取方式（Pull-based）
        2.1 某个节点周期性地选择随机N个节点询问有没有最新的信息
        2.2 收到请求的节点回复请求节点其最近未收到的信息

### raft选举协议
// to do

### consul常用端口
- 服务器RPC（默认8300）：由服务器用来处理来自其他代理的传入请求，仅限TCP。
- Serf LAN（默认8301）：用来处理局域网中的八卦。所有代理都需要，TCP和UDP。
- Serf WAN（默认8302）：被服务器用来在WAN上闲聊到其他服务器，TCP和UDP。从Consul 0.8开始，建议通过端口8302在LAN接口上为TCP和UDP启用服务器之间的连接，以及WAN加入泛滥功能。
- HTTP API（默认8500）：被客户用来与HTTP API交谈，仅限TCP。
DNS接口（默认8600）：用于解析DNS查询，TCP和UDP。
## 启动consul
启动容器:
```
docker run -d --name=dev-consul -e CONSUL_BIND_INTERFACE=eth0 consul
```
这样执行默认的参数为:`consul agent -data-dir=/consul/data -config-dir=/consul/config -bind=172.17.0.3 -dev -client 0.0.0.0`

```bash
consul agent -server -bootstrap-expect=3 -data-dir=/tmp/consul -node=10.200.110.90 -bind=10.200.110.90 -client=0.0.0.0 -datacenter=shenzhen -ui
```
### 单机执行consul
`consul agent -server -data-dir=/tmp/consul -bootstrap -advertise=<localIp>`
参数:
- bootstrap: 此标志用于控制服务器是否处于“引导”模式。在此模式下，每个数据中心运行的服务器不得超过一台
- advertise: 用于更改到群集中其他节点的地址; 

consul agent貌似必须设置bind参数或者adverise参数,不然会报错:
```
Multiple private IPv4 addresses found. Please configure one with 'bind' and/or 'advertise'.
```
## consul command命令:
查看端口占用: `lsof -i :8500`

1. 查看当前成员列表以及状态`consul members`, `-detailed`显示更详细的信息;
    docker: `docker exec -t <containerId> consul members`

2. 加入cluster, `consul join + <>`

## consul option命令:

### 常用:
- server： 以server身份启动。默认是client;
- client: 客户端模式;
- ui: 可以访问UI界面

- bootstrap: 此模式下，节点可以选举自己为leader，一个数据中心只能有一个此模式启动的节点。机群启动后，新启动的节点不建议使用这种模式。
- bootstrap-expect：集群要求的最少server数量，当低于这个数量，集群即失效。

- node：节点id，集群中的每个node必须有一个唯一的名称。默认情况下，Consul使用机器的hostname
- disable-host-node-id：不使用host信息生成node ID，适用于同一台服务器部署多个实例用于测试的情况。随机生成nodeID

- bind：绑定的内部通讯地址，默认0.0.0.0，即，所有的本地地址，会将第一个可用的ip地址散播到集群中，如果有多个可用的ipv4，则consul启动报错。


- datacenter 指定数据中心名称，默认是dc1

文件夹配置: 
- data-dir：状态数据存储文件夹，所有的节点都需要。文件夹位置需要不收consul节点重启影响，必须能够使用操作系统文件锁，unix-based系统下，文件夹文件权限为0600，注意做好账户权限控制，
- config-file：配置文件位置
- config-dir: 指定配置文件夹，Consul会加载其中的所有.json或者.hcl文件, 不会加载子文件夹;

副节点则需要以下参数:
- join: 启动时要加入的另一个代理的地址。可以多次指定此选项以指定要加入的多个代理。如果consur无法加入任何指定的地址，代理启动将失败。默认情况下，代理在启动时不会加入任何节点。请注意，在自动化consur集群部署时，使用`retry_join`可能更适合于帮助缓解节点启动争用情况。
- retry-join: 类似于-join，但允许在连接成功之前重试连接。一旦它成功加入到成员列表中的成员，它就再也不会尝试加入。agents将通过gossip维持他们的会员资格。这对于您知道地址最终将可用的情况非常有用。

## consul go
### 使用consul api实现配置中心
> 依赖包: github.com/hashicorp/consul/api


### go-micro实现服务注册
> 使用包: github.com/micro/go-plugins/registry/consul/v2
```
// 注册中心
consul.NewRegistry(func(options *registry.Options) {
    options.Addrs = []string{
        consulHost + ":8500",
    }
})
micro.NewService(micro.Registry(consulRegister))
```

### go-micro实现 配置中心
> 依赖包: github.com/micro/go-plugins/config/source/consul/v2
>
> github.com/micro/go-micro/v2/config 
```go
consulSource := consul.NewSource(
		//设置配置中心的地址
		consul.WithAddress(host+":"+strconv.FormatInt(port, 10)),
		//设置前缀，不设置默认前缀 /micro/config
		consul.WithPrefix(prefix),
		//是否移除前缀，这里是设置为true，表示可以不带前缀直接获取对应配置
		consul.StripPrefix(true),
	)
	config.NewConfig()
	conf.Load(consulSource)
```

