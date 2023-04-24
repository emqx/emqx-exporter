## Dashboard Instruction
Grafana dashboard has been split multipart to display different metrics.
* **EMQX** the main dashboard, it shows the overall metrics of the EMQX cluster.
* **Authentication** client connection metrics statistics.
* **Authorization** client ACL metrics statistics.
* **Client Events** client events that handled successfully by system. 
* **Messages** client packets that system received.
* **Rule Engine** message flow and device event processing and response metrics statistics.

The main dashboard `EMQX` shows the overall metrics of the EMQX cluster, you can switch data source and cluster by the global variable at the top of the dashboard.  

In addition, some panels have its own set of links that are shown in the upper left corner of the panel, the links will jump to relevant dashboard in a new tab, you can filter the metrics based on global variables, such as `cluster`,`node`, and so on.

## Metrics Instruction
Here we mainly introduce the metrics in the main dashboard `EMQX`.

### Cluster Status
集群健康状态

### License
License 信息，包括剩余天数、到期日期以及最大连接数

### Active Connections
活跃连接数

### Cluster Message Rate
集群每秒消息流入流出数

### Nodes Running
集群当前正在运行的节点数与已停止的节点数

### Connections
集群的总连接数与各节点的连接数

### Subscriptions
集群当前的总订阅数

### Sessions
集群当前的session数量，包括活跃连接与已断开连接，但未清除session的数量

### Rule Engine Last 5m Exec Rate
规则引擎最近五分钟的执行速率

### Connect Auth
* **Auth Success** 认证成功的客户端 CONNECT 报文数
* **Auth Failure** 认证失败的客户端 CONNECT 报文数

### ACL Auth
* **Publish ACL Failure** 认证失败的客户端 PUBLISH 报文数量
* **Sub ACL Failure** 认证失败的客户端 SUBSCRIBE 报文数量

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

### Erlang VM Run Queue
The value of `run_queue_length` can be used to evaluate the current Erlang scheduler load. If the value of run_queue_length is high, it means that there are many processes waiting to be executed in the run queue, and the scheduler may need more processor cores or a better scheduling strategy to better handle these processes

### Erlang VM Process

### Erlang VM Memory Used
Erlang VM 的内存占用

### Mnesia(built-in database) Memory Usage
在 EMQX5 中表示 Mnesia 数据库的磁盘空间使用