# protobuf
> [protobuf教程](https://geektutu.com/post/quick-go-protobuf.html)

`protoc -I ./ --go_out=./ --micro_out=./ *.proto`


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
 