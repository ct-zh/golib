# sentinel HA

1. 一共需要 sentinel 3台，redis 3台以上；
2. 一台机器起三个 sentinel，一台机器起三个 redis；  
3. 需要修改配置文件防止端口冲突;

## 流程
> 在`sentinel.conf`文件里将本次HA集群命令为`server1`
> 
> 注意对应主机的端口是否在防火墙里开放;

```bash
# 1. 发送sentinel脚本到 198机器上
scp -r sentinel/ caot@192.168.199.198:/home/caot/sentinel

# 2. 发送redis脚本到 197 机器上
scp -r redis/ caot@192.168.199.197:/home/caot/redis

# 3. 先执行197机器上,建立redis主从
cd ~/redis
docker-compose -f redis.yml up -d

# 检查建立情况
docker exec server1 redis-cli -p 12221 INFO REPLICATION
docker exec server2 redis-cli -p 12222 INFO REPLICATION

# 4. 在198机器上建立sentinel集群
cd ~/sentinel
/bin/bash server1.sh

# 检查sentine集群情况
# 查看master
docker exec -it sentinel1 redis-cli -p 26379 SENTINEL master server1

# 查看slave
docker exec -it sentinel1 redis-cli -p 26379 SENTINEL slaves server1

# 查看日志
docker logs -f sentinel1

```

## failover

```bash
# 在197机器上 关闭master节点
docker stop server1

# 检查slave的状态
docker exec server2 redis-cli -p 12222 INFO REPLICATION

```

正确情况下,过一段时间sentinel会重新设定新的master节点; 

### 出现问题
在写这个readme时这个HA是建立不成功的,每次sentinel在 12221这台master下线后又重新将其设置为master; 经过漫长的排查问题,发现docker compose会建立一个redis_default网络,相当于三个redis在一个子网,三个sentinel在另一个子网;

然后当sentinel从197机器上的master更新信息时, master传给sentinel的slaves信息是根据他的网络来的,也就是说:
- 对于sentinel,他只知道master的地址是 `192.168.199.197 12221`
- 对于master,他知道sentinel的地址是`192.168.199.198 26379`,他的slave地址是`172.20.0.2 12222`等等
- 于是master把slaves的地址传给sentinel
- sentinel根据master给的ip地址去访问,结果访问不到




