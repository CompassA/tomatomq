-- podman exec -it mysql-master bash
-- mysql -uroot -proot
-- 进入mysql终端,创建数据库

CREATE DATABASE tomato_mq_admin;
USE tomato_mq_admin;

-- 数据库信息表
CREATE TABLE `tomato_mq_db`(
    `id` BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
    `db_guid` VARCHAR(512) NOT NULL COMMENT 'DB唯一标识, 目前存DB的名称',
    `db_dsn` VARCHAR(2048) NOT NULL COMMENT '数据库资源连接信息',
    `broker_group` VARCHAR(128) NOT NULL COMMENT '绑定的broker分组',
    `gmt_create` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `gmt_modified` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   
    UNIQUE KEY(`db_guid`),
    INDEX `idx_broker_group`(`broker_group`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- topic信息表
CREATE TABLE `tomato_mq_topic`(
    `id` BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
    `name` VARCHAR(64) NOT NULL COMMENT 'topic名称',
    `broker_group` VARCHAR(128) NOT NULL COMMENT '绑定的消息队列分组',
    `msg_type` INT(11) NOT NULL COMMENT '消息类型: 0=普通消息,1=顺序消息,2=事务消息,3=定时消息',
    `msg_queue_num` INT(11) NOT NULL COMMENT '底层队列数',
    `status` VARCHAR(32) NOT NULL COMMENT 'topic创建状态: CREATING=创建中, ONLINE=生效, OFFLINE=下线',
    `gmt_create` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `gmt_modified` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY `uni_name`(`name`),
    INDEX `idx_broker_group`(`broker_group`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- topic底层队列信息表
CREATE TABLE `tomato_mq_mysql_queue`(
    `id` BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
    `topic_id` BIGINT NOT NULL COMMENT '属于哪个topic',
    `index` INT(11) NOT NULL COMMENT '队列索引, 0,1,2 .... N-1',
    `db_id` BIGINT NOT NULL COMMENT '属于哪个数据库资源',
    `table_name` VARCHAR(512) NOT NULL COMMENT '关联的消息表名',
    `status` VARCHAR(32) NOT NULL COMMENT '队列状态: OFFLINE 不可读写,READ_ONLY 只读, WRITE_ONLY 只写, ONLINE 允许读写',
    `gmt_create` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `gmt_modified` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY `uni_topic_index`(`topic_id`,`index`),
    INDEX `idx_dbid`(`db_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- broker与topic底层队列的绑定关系
CREATE TABLE `tomato_mq_broker_queue_relation`(
    `id` BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
    `broker_group` VARCHAR(128) NOT NULL COMMENT 'broker分组',
    `broker_name` VARCHAR(128) NOT NULL COMMENT 'broker名称',
    `queue_id` BIGINT NOT NULL COMMENT '绑定的队列ID',
    `gmt_create` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `gmt_modified` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY (`queue_id`,`broker_name`,`broker_group`),
    INDEX `idx_broker`(`broker_name`, `broker_group`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 创建broker资源库 本地创建4个
CREATE DATABASE tomato_mq_msg_0;
CREATE DATABASE tomato_mq_msg_1;
CREATE DATABASE tomato_mq_msg_2;
CREATE DATABASE tomato_mq_msg_3;


