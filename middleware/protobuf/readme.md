# protobuf
protobuf 即 Protocol Buffers，是一种轻便高效的结构化数据存储格式，与语言、平台无关，可扩展可序列化。protobuf 性能和效率大幅度优于 JSON、XML 等其他的结构化数据格式。protobuf 是以二进制方式存储的，占用空间小，但也带来了可读性差的缺点。

[官方文档](https://developers.google.com/protocol-buffers/docs/overview)



## 错误
报错
```
github.com\coreos\etcd@v3.3.22+incompatible\clientv3\balancer\resolver\endpoint\endpoint.go:114:78: undefined: resolver.BuildOption
...
```
大概是说原因是google.golang.org/grpc 1.26后的版本是不支持clientv3的。也就是说要把这个改成1.26版本的就可以了。
具体操作方法是在go.mod里加上：
```
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
```

## reference
> - [protobuf教程](https://geektutu.com/post/quick-go-protobuf.html)
