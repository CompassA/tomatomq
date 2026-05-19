<!--
 * @Author: Tomato
 * @Date: 2026-05-16 18:49:37
-->
# 项目简介

fixme TOMATO todo

# 管理时设计

## 架构设计

## 领域模型
![admin_model.svg](./docs/output/admin_model.svg "管理时领域模型")

- **topic**:消息逻辑分组, 存放同种类型的消息  
- **message_queue**: 消息分区, 一个topic可以有多个消息分区 
- **broker**: 消息队列服务端, 负责处理消息收发网络请求  
- **db**: 消息存储端, tomato底层使用数据库存储消息  
- **broker_group**: 消息队列集群逻辑分组, topic的每个队列会绑定给相同分组的broker管理，队列的数据会存储在相同分组的db中  


## 数据模型
![admin_data_model.svg](./docs/output/admin_data_model.svg "管理时数据模型")

## 集群管理整体流程
![集群操作](./docs/output/BrokerGroup创建.svg "集群操作")

- 创建数据库资源
- 创建broker实例
- 创建topic
- broker注册逻辑
  - 使用etcd作为注册中心、存储topic的原数据与broker的机器信息。 
  - 每个broker启动时, 将自己的机器信息写入ectd, 并创建一个etcd租约, 将etcd租约与写入的ip信息绑定。 
  - ectd写入格式:  
  key:/broker/{broker分组}/{broker名称}   
  value: ip:port

# 通信层设计

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

# 安装/启动mysql
db_root_dir=/root/download/mysqldata
mkdir -p $db_root_dir/master/conf
mkdir -p $db_root_dir/master/data

podman run --name mysql-master -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root \
-v $db_root_dir/master/conf:/etc/mysql/conf.d                   \
-v $db_root_dir/master/data:/var/lib/mysql                      \
-d mysql:8.4.7

podman exec -it mysql-master bash
```

# sql
```sql
-- 消息库, 每个broker分组有几个broker就创建几个DB
CREATE DATABASE tomato_broker_0;

-- 
```
