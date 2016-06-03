## 内存优先消息队列实现， 基于idmm2 java版重写

### 下面是功能列表


### 一、 DONE

* 内存优先级排序， 分组按时间顺序
 * 结构体定义
 * 优先级链表实现， 一个基本的双向链表， 一个带sort key的排序双向链表
 * PUSH 新增消息
 * PULL 拉取消息， 消息状态变为 "在途", 消息重试次数+1
  * 锁定时间设置
  * -- 生效时间 EFFECTIVE_TIME
  * -- 失效时间 EXPIRE_TIME
 * COMMIT 消费确认, 把 "在途" 的对应消息从队列中删除
 * ROLLBACK 回滚， "在途"状态改为0， 只处理在途的消息
 * DELETE   删除消息，非在途也可以处理
 * FAIL     消费失败， 不再处理， 只处理在途的
 * 最大重试次数， 达到最大重试次数后， 消息按FAIL处理


* 实现TCP接口的包发送和接收
 * 包结构解析， 相关常量和结构体定义， Reader上的报文读取并还原为结构体
 * 结构体写入 Writer
 * 请求响应函数实现，按报文类型分配到对应的处理函数

* 网络并发请求与内存操作的单线程之间的隔离， 使用chan来实现， 与java版方法相同
* 与 java 版 broker 联调成功下列的消息类型
 * SEND-COMMIT， 支持单个消息和批量消息id(BATCH_MESSAGE_ID)的处理
 * PULL
 * PULL COMMIT
 * PULL COMMIT_AND_NEXT
 * PULL ROLLBACK
 * PULL ROLLBACK_AND_NEXT
 * PULL ROLLBACK_AND_RETRY
 * -- DELETE
 * -- UNLOCK

### 二、 TODO

#### 内存队列benchmark测试用例
#### mysql数据库配置读取和更新
 * 连接mysql
 * 配置结构体定义
 * 读取配置表并生成配置数据结构
 * 根据配置构建消息队列， 以及绑定端口
 * 在线重新加载配置， 并更新到运行环境
#### zookeeper 处理
 * 连接 zk， 并读取配置版本号
 * 建立ble的id临时节点， 发布端口号

#### mysql索引数据持久化实现
 * 消息处理过程中， 先更新db， 再更新内存
  * 流程定义
  * 对应结构体定义
#### 消息处理过程持久化
  * SEND-COMMIT
  * PULL
  * COMMIT
  * ROLLBACK
  * FAIL
  * DELETE
#### 重启后从mysql中恢复数据到内存
#### jmx 端口上的运行情况监控
 * 内存队列积压情况
 * 手工发送 删除 消息

#### 优先级映射

