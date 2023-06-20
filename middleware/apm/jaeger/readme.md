# jaeger

## 组成

- Jaeger Client： 符合OpenTracing标准的SDK。应用程序通过API写入数据， client library把trace信息按照采样策略传递给jaeger-agent。
- Agent：监听在UDP端口上接收span数据的网络守护进程，会将数据批量发送给collector。作为基础组件部署到所有的宿主机上。Agent将client library和collector解耦，为client library屏蔽了路由和发现collector的细节。
- Collector：接收jaeger-agent发送来的数据，然后将数据写入后端存储。Collector被设计成无状态的组件，因此用户可以运行任意数量的Collector;
- Data Store：后端存储被设计成一个可插拔的组件，支持数据写入cassandra， elastic search;
- 