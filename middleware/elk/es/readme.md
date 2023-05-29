# elasticsearch这一篇就够了

## 基本概念+使用入门
### 快速启动
> 5.29修改： 可以试一试compose docker-compose up -d
> [compose file](./demo/docker-compose.yml)

```shell
# 安装elasticsearch
docker run -p 9200:9200 -p 9300:9300 -d -e "discovery.type=single-node" docker.elastic.co/elasticsearch/elasticsearch:7.16.3

# 确定服务是否启动
curl -X GET "http://0.0.0.0:9200"

# 安装kibana
docker run --name kibana -p 5601:5601 -e ELASTICSEARCH_HOSTS=http://<elasticsearch_ip>:9200 -d kibana:7.16.3
# 其中，<elasticsearch_ip>是您的Elasticsearch IP地址,注意，这里应该填容器所绑定的网络地址,所以可能需要指定一下网络
docker run --name kibana -p 5601:5601 --network=NETWORK_NAME -e ELASTICSEARCH_HOSTS=http://<elasticsearch_ip>:9200 -d kibana:7.16.3

# 查看kibana是否启动，访问地址 0.0.0.0:5601
# kibana在7.x版本后自带了官方汉化，位置在kibana目录下的：
# `node_modules/x-pack/plugins/translations/translations/`
# 或者`x-pack/plugins/translations/translations/`
# 然后修改kibana配置文件kibana.yml：`i18n.locale: "zh-CN"`

# go导入包 github.com/elastic/go-elasticsearch

```


确定是否能使用es了：[see  connect连接es](./demo/main.go)


### 概念
Elasticsearch的基本概念有以下几个：

- 节点（Node）：运行了单个实例的ES主机称为节点，它是集群的一个成员，可以存储数据、参与集群索引及搜索操作。节点通过为其配置的ES集群名称确定其所要加入的集群。
- 集群（Cluster）：ES可以作为单机运行，也可以作为多台主机运行，多台主机组成的一个整体就是集群。集群有一个唯一的名称，用来区分不同的集群。
- 索引（Index）：索引是一类具有相似特征的文档的集合，类似于数据库中的数据库。索引有一个唯一的名称，用来标识不同的索引。
- 类型（Type）：类型是索引中的一个逻辑分类，类似于数据库中的表。类型有一个唯一的名称，用来标识不同的类型。
- 文档（Document）：文档是ES中存储和检索的基本单位，类似于数据库中的行。文档是一个JSON对象，包含了多个字段和值。
- 字段（Field）：字段是文档中的一个属性，类似于数据库中的列。字段有一个名称和一个类型，用来标识不同的字段。
- 映射（Mapping）：映射是对索引中类型和字段的定义，类似于数据库中的表结构。映射可以指定字段的类型、分析器、格式等属性。
- 分片（Shard）：分片是索引中数据的物理分割，用来实现水平扩展和负载均衡。每个分片都是一个完整的Lucene实例，可以在集群中任意节点上移动。
- 副本（Replica）：副本是分片的复制，用来实现高可用和容错。每个分片可以有多个副本，副本可以在集群中任意节点上移动。

其中，索引、类型、文档、字段等一般是开发人员角度所需要掌握的内容。而节点、集群、分片则是运维人员角度所需要掌握的内容。


#### 索引、文档



## 安全性

> see: [configuring-stack-security](https://www.elastic.co/guide/en/elasticsearch/reference/7.16/configuring-stack-security.html)

