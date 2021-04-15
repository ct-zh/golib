# kafka
## 入门
> [小朋友也可以懂的Kafka入门教程](https://www.bilibili.com/video/BV1vx411f7hA)
- 跟普通消息队列一样,从producer处获得消息, 发给comsumer;
- 消息区分topic,也就是主题;
- 一个topic下面可能有很多个partition,也就是分区, 消息投递给主题时可以指定分区;
- 如果没指定分区, 分区器(一种算法)会根据消息键(消息的标记)投递到某个分区;
- 所以一个消息包含: 主题、分区、键、值;
- 消费者不是直接拿走消息,而是根据偏移量读取消息;
- 一个broker代表一个kafka服务, 多个borker组成集群;
- broker以复制的形式维持集群一致性;

