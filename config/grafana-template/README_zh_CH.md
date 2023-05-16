## Dashboard Instruction

EMQX 指标分成多个不同的 dashboard 进行展示。

* **EMQX** 主看板, 用于展示集群整体的指标数据。
* **Authentication** 客户端认证相关指标，可按集群、节点和第三方组件进行筛选查看。
* **Authorization** 客户端 ACL 授权相关指标，可按集群、节点和第三方组件进行筛选查看。
* **Client Events** EMQX 规则引擎中对相关事件的执行统计，可按集群和节点进行筛选查看。
* **Messages** EMQX 集群接收到的 MQTT 报文统计，可按集群和节点进行筛选查看。
* **Rule Engine** EMQX 集群规则引擎相关指标，可按集群、节点和数据处理组件进行筛选查看。

Dashboard `EMQX` 展示的是集群整体的指标。在 dashboard 的顶部可以根据数据源和集群名称进行切换。
此外，部分面板设置是跳转链接（面板左上角），点击后跳转到相关的 dashboard 页面，在新跳转的 dashboard 上可以根据集群、节点等类型进行筛选查看。

## Metrics Instruction

本文仅对 `EMQX` dashboard 中的指标进行简单的介绍

### Cluster Status

集群健康状态。当集群中任意个节点异常时，该看板将显示 `Unhealthy`，背景色变为红色。

### License

License 信息，包括剩余天数、到期日期以及最大连接数。当剩余天数不足90天时，背景色变为黄色，不足30天时，则为红色。

### Active Connections

活跃连接数。显示当前活跃的客户端连接数的仪表盘，当连接数超过 license 最大连接数的80%时，仪表盘显示为黄色，超过90%时显示为红色。

### Cluster Message Rate

集群每秒消息流入流出数(TPS)。

### Nodes Running

集群当前正在运行的节点数与已停止的节点数

### Connections

集群的总连接数与各节点的连接数

### Subscriptions

集群当前的总订阅数

### Sessions

集群当前的 session 数量，包括活跃连接与已断开连接，但未清除session的数量

### Rule Engine Last 5m Exec Rate

规则引擎最近五分钟的执行速率

## Data Bridge Queuing

当前规则引擎中被缓存到磁盘队列的消息个数

### Connect Auth

* **Auth Success** 认证成功的客户端 CONNECT 报文数
* **Auth Failure** 认证失败的客户端 CONNECT 报文数

### ACL Auth

* **Publish ACL Failure** 授权失败的客户端 PUBLISH 报文数量
* **Sub ACL Failure** 授权失败的客户端 SUBSCRIBE 报文数量

### Client Connection Events

* **Connections** `client.connected`钩子触发次数
* **Disconnections** `client.disconnected`钩子触发次数

### Client Sub Events

* **Subscribes** `client.subscribe`钩子触发次数
* **Unsubscribes** `client.unsubscribe`钩子触发次数

### Client Connect Auth Events

* **Connect Auth** `client.authenticate`钩子触发次数
* **Anonymous Auth** 客户端以匿名方式认证成功次数

### Client ACL Auth Events

* **ACL Auth** 在 EMQX4 中表示`client.check_acl`钩子触发次数， 在 EMQX5 中表示`client.authorize`钩子触发次数

### Packets Connections

* **Packets Connect** 接收的 CONNECT 报文数量
* **Packets Connack Sent** 发送的 CONNACK 报文数量
* **Packets Connack Error** 发送的原因码不为 0x00 的 CONNACK 报文数量

### Packets Disconnections

* **Packets Disconnect Sent** 发送的 DISCONNECT 报文数量
* **Packets Disconnect Received** 接收的 DISCONNECT 报文数量

### Packets Publish

* **Packets Publish Sent** 发送的 PUBLISH 报文数量
* **Packets Publish Received** 接收的 PUBLISH 报文数量
* **Packets Publish Dropped** 超出接收限制而被丢弃的 PUBLISH 报文数量
* **Packets Publish Error** 接收的无法被发布的 PUBLISH 报文数量

### Packets Subscribe/Unsubscribe

* **Packets Subscribe Received** 接收的 SUBSCRIBE 报文数量
* **Packets Suback Sent** 发送的 SUBACK 报文数量
* **Packets Subscribe Error** 接收的订阅失败的 SUBSCRIBE 报文数量
* **Packets Unsubscribe Received** 接收的 UNSUBSCRIBE 报文数量
* **Packets Unsubscribe Error** 接收的取消订阅失败的 UNSUBSCRIBE 报文数量

### Messages Count

* **Messages Received** 接收来自客户端的消息数量
* **Messages Sent** 发送给客户端的消息数量

### Erlang VM Messages Queue

Erlang VM 中未处理的消息队列长度

### Messages QOS Received

集群接收到的QOS0~QOS1的消息的增量变化趋势

### Cluster Traffic Statistics

* **Bytes received** 集群接收到的消息字节数
* **Bytes sent** 集群发送的消息字节数

### Data Bridge Status

规则引擎第三方资源的连接状态及当前待处理请求的队列长度

### Rule Engine Execute Count

规则引擎的统计数据，包括 topic 命中数、最大执行速率、成功数、失败数等。

### Rule Engine Current Exec Rate

规则引擎的当前执行速率

### Data Bridge Queuing

规则引擎第三方资源的待处理请求的队列长度

### Authenticate Count

客户端认证(Connect)插件的连接状态及统计数据

### Authenticate Current Exec Rate

客户端认证(Connect)插件的当前执行速率

### Authenticate Last 5m Exec Rate

客户端认证(Connect)插件的最近5分钟的执行速率

### Authorize Count

客户端授权(Pub&Sub)插件的连接状态及统计数据

### Authorize Current Exec Rate

客户端授权(Pub&Sub)插件的当前执行速率

### Authorize Last 5m Exec Rate

客户端授权(Pub&Sub)插件的最近5分钟的执行速率

### Last 1m CPU Load

The last 1m CPU load of node

### Last 5m CPU Load

The last 5m CPU load of node

### Last 15m CPU Load

The last 15m CPU load of node
