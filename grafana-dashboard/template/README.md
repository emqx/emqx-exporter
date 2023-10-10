## Other Languages

[中文](./README_zh_CH.md)

## Dashboard Instruction

Grafana dashboard has been split multipart to display different metrics.

* **EMQX** the main dashboard, it shows the overall metrics of the EMQX cluster.
* **Authentication** client connection metrics statistics.
* **Authorization** client ACL metrics statistics.
* **Client Events** the events statistics that handled successfully in rule engine.
* **Messages** client packets that system received.
* **Rule Engine** message flow and device event processing and response metrics statistics.

The main dashboard `EMQX` shows the overall metrics of the EMQX cluster, you can switch data source and cluster by the global variable at the top of the dashboard.

In addition, some panels have its own set of links that are shown in the upper left corner of the panel, the links will jump to relevant dashboard in a new tab, you can filter the metrics based on global variables, such as `cluster`,`node`, and so on.

## Metrics Instruction

Here we mainly introduce the metrics in the main dashboard `EMQX`.

### Cluster Status

The healthy status of cluster

### License

The license table contains days remaining, expiration day and the max client connection num.

### Active Connections

The total active connections

### Cluster Message Rate

The messages input/output per second(TPS).

### Nodes Running

The running/stopped nodes of cluster.

### Connections

The total connections that contains active connections and the remaining sessions after client disconnected

### Subscriptions

The total subscriptions of cluster.

### Sessions

The total connections that contains active connections and the remaining sessions after client disconnected.

### Rule Engine Last 5m Exec Rate

The last 5m exec rate of rule engine

## Data Bridge Queuing

The Count of messages that are currently queuing in the rule engine

### Connect Auth

* **Auth Success** Number of CONNECT packets from authenticated clients.
* **Auth Failure** Number of CONNECT packets from unauthenticated clients.

### ACL Auth

* **Publish ACL Failure** Number of PUBLISH packets from authorized clients.
* **Sub ACL Failure** Number of SUBSCRIBE packets from unauthorized clients.

### Client Connection Events

* **Connections** `client.connected`hook trigger times
* **Disconnections** `client.disconnected`hook trigger times

### Client Sub Events

* **Subscribes** `client.subscribe`hook trigger times
* **Unsubscribes** `client.unsubscribe`hook trigger times

### Client Connect Auth Events

* **Connect Auth** `client.authenticate`hook trigger times
* **Anonymous Auth** Number of successful anonymous authentication

### Client ACL Auth Events

* **ACL Auth** It represents `client.check_acl`hook trigger times in EMQX 4, and`client.authorize`hook trigger times in EMQX 5

### Packets Connections

* **Packets Connect** Number of received CONNECT packets
* **Packets Connack Sent** Number of sent CONNACK packets
* **Packets Connack Error** Number of sent CONNACK packets where reason code is not 0x00.

### Packets Disconnections

* **Packets Disconnect Sent** Number of sent DISCONNECT packets
* **Packets Disconnect Received** Number of received DISCONNECT packets

### Packets Publish

* **Packets Publish Sent** Number of sent PUBLISH packets
* **Packets Publish Received** Number of received PUBACK packets
* **Packets Publish Dropped** Number of PUBLISH packets that were discarded due to the receiving limit
* **Packets Publish Error** Number of received PUBLISH packets that cannot be published

### Packets Subscribe/Unsubscribe

* **Packets Subscribe Received** Number of received SUBSCRIBE packets
* **Packets Suback Sent** Number of sent SUBACK packets
* **Packets Subscribe Error** Number of received SUBSCRIBE packets with failed subscriptions
* **Packets Unsubscribe Received** Number of received UNSUBSCRIBE packets
* **Packets Unsubscribe Error** Number of received UNSUBSCRIBE packets with failed unsubscriptions

### Messages Count

* **Messages Received** Number of messages received from the client
* **Messages Sent** The number of messages sent to the client

### Messages QOS Received

The incremental change trend of QOS0~QOS1 messages received by the cluster.

### Cluster Traffic Statistics

* **Bytes received** Number of received bytes
* **Bytes sent** Number of send bytes

### Data Bridge Status

The status of rule engine resources and the current queue length of pending requests.

### Rule Engine Execute Count

Statistics of the rule engine, including topic hits, maximum execution rate, success count, failure count, etc.

### Rule Engine Current Exec Rate

The current execution rate of the rule engine.

### Data Bridge Queuing

The current queue length of pending requests of rule engine.

### Authenticate Count

Connection status and statistics of the client authentication (Connect) plugin.

### Authenticate Current Exec Rate

The current execution rate of the client authentication (Connect) plugin.

### Authenticate Last 5m Exec Rate

The execution rate of the client authentication (Connect) plugin in the last 5 minutes.

### Authorize Count

Connection status and statistics of the client authorization (Pub&Sub) plugin.

### Authorize Current Exec Rate

The current execution rate of the client authorization (Pub&Sub) plugin.

### Authorize Last 5m Exec Rate

Connection status and statistics of the client authorization (Pub&Sub) plugin.

### Last 1m CPU Load

The last 1m CPU load of node

### Last 5m CPU Load

The last 5m CPU load of node

### Last 15m CPU Load

The last 15m CPU load of node
