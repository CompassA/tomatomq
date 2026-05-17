<!--
 * @Author: Tomato
 * @Date: 2026-05-16 18:49:37
-->
# broker设计


## 元数据设计

使用etcd作为注册中心、存储topic的原数据与broker的机器信息


### broker注册

每个broker启动时, 将自己的机器信息写入ectd, 并创建一个etcd租约, 将etcd租约与写入的ip信息绑定。

ectd写入格式:
key:/broker/{broker分组}/{broker名称} 
value: ip:port


broker每30s定期续租, 若broker宕机, 租约将过期, etcd将删除对应broker机器的ip


### topic队列与broker的绑定关系

用户向admin服务提交topic创建请求, topic创建参数包括topic名称、总队列数、mq分组名称、

admin处理逻辑:

1. 生成topic队列id
2. 获取对应mq分组中的所有broker, 按平均分配策略, 为每个队列分配一个broker
3. 将分配结果持久化写入etcd


etcd写入格式
key: /topic/{队列id}
value: { "g": "队列分组", "n": "队列名称"}


创建完成后, topic的队列即永远和一个broker实例绑定。

## TCP层

核心领域模型

Session: 代表一个客户端连接。
Server: 代表broker服务端

## 协议层


# 本地启动
```bash
# 安装/启动etcd
NODE1=127.0.0.1
DATA_DIR=/root/download/etcddata
REGISTRY=quay.io/coreos/etcd
ETCD_VERSION=v3.6.11
IMAGE_NAME=my_etcd

setenforce 0
 
podman run -it -d -p 2379:2379 -p 2380:2380 \
  --user root \
  --volume=${DATA_DIR}:/etcd-data \
  --name ${IMAGE_NAME} ${REGISTRY}:${ETCD_VERSION} \
  /usr/local/bin/etcd \
  --data-dir=/etcd-data --name node1 \
  --initial-advertise-peer-urls http://${NODE1}:2380 --listen-peer-urls http://0.0.0.0:2380 \
  --advertise-client-urls http://${NODE1}:2379 --listen-client-urls http://0.0.0.0:2379 \
  --initial-cluster node1=http://${NODE1}:2380

podman exec ${IMAGE_NAME} etcdctl endpoint health
```
