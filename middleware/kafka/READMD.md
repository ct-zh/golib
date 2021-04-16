# kafka
## 入门
> [官方文档 ](http://kafka.apache.org/documentation/)

> [小朋友也可以懂的Kafka入门教程](https://www.bilibili.com/video/BV1vx411f7hA)

- 跟普通消息队列一样,从producer处获得消息, 发给comsumer;
- 消息区分topic,也就是主题;
- 一个topic下面可能有很多个partition,也就是分区, 消息投递给主题时可以指定分区;
- 如果没指定分区, 分区器(一种算法)会根据消息键(消息的标记)投递到某个分区;
- 所以一个消息包含: 主题、分区、键、值;
- 消费者不是直接拿走消息,而是根据offset偏移量读取消息;
- 一个broker代表一个kafka服务, 多个borker组成集群;
- broker以复制的形式维持集群一致性;复制的是partition log

### 消费者组
Consumer Group: 消费者组，由多个 consumer 组成。消费者组内每个消费者负 责消费不同分区的数据，一个分区只能由一个组内消费者消费;消费者组之间互不影响。所 有的消费者都属于某个消费者组，即消费者组是逻辑上的一个订阅者。

### zookeeper的作用
Kafka 集群中有一个 broker 会被选举为 Controller，负责管理集群 broker 的上下线，所有topic的分区副本分配和 leader 选举等工作。Controller是通过zookeeper对其他节点进行控制的.

## 基本使用
cmd: 
> 删除:需要 server.properties 中设置 delete.topic.enable=true 否则只是标记删除;

```
// 查看topic列表 需要连接zoo1
kafka-topics.sh --zookeeper zoo1:2181 --list

// 查看topic详情
kafka-topics.sh --zookeeper zoo1:2181 --describe --topic firsttopic

// 创建topic, --topic指定topic名称, 
// replication-factor:副本数 ; partitions: 定义分区数
kafka-topics.sh --zookeeper zoo1:2181 --create --replication-factor 3 --partitions 1 --topic firsttopic

// 删除topic; 
kafka-topics.sh --zookeeper zoo1:2181 --delete --topic firsttopic

// 发送消息
kafka-console-producer.sh --broker-list kafka1:9092 --topic firsttopic

// 消费消息
// --from-beginning:会把主题中以往所有的数据都读取出来。
kafka-console-consumer.sh --zookeeper zoo1:2181 --topic firsttopic
kafka-console-consumer.sh --bootstrap-server kafka1:9092 --from-beginning --topic first
```


## 消息丢失、消息重复、消息顺序
### 消息丢失
sender设计,提供两个api:`SendMsg(msg bytes[])`(发送消息),`SendCallback()`;receiver设计,提供两个api:`RecvCallback(msg bytes[])`与`SendAck()`;流程如下:
- sender使用`SendMsg`发送消息给MQ;
- MQ消息落地,调用`SendCallback`将应答消息发给sender;(保证MQ必定收到消息)
- SendMsg在未收到Callback时timer会重发消息;一般采用指数退避的策略，先隔x秒重发，2x秒重发，4x秒重发，以此类推
- receiver使用`RecvCallback`从MQ收到消息;
- receiver使用`SendAck`向MQ申明已经收到消息了; MQ删除对应数据;MQ在未收到ack时会不停重发; (确保消息必达)
- *sender会重发很多次消息, receiver会收到很多重复消息, 因此接口的幂等一定要做好*
### 消息重复
做幂等设计, 一个是根据业务来去重,比如消息里面带有全局唯一的支付id、订单id之类的; 一个是CAS(compare and set): `update t set money=28 where uid=xx and money=100`
开启事务时`select for update`

ABA问题: 事务A取出money=100, 事务B扣款成80, 事务C加钱成100, 此时仍然满足`and money=100`条件; 解决办法:增加版本号字段, 每次修改都需要更新版本号, update从money的比对优化为版本号的比对;

### 消息顺序

## go kafka
> 依赖包: "github.com/Shopify/sarama"

