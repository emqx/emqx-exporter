## User Manual
Grafana dashboard has been split multipart to display different metrics.
* **EMQX** the main dashboard, it shows the overall metrics of the EMQX cluster.
* **Authentication** client connection metrics statistics.
* **Authorization** client ACL metrics statistics.
* **Client Events** client events that handled successfully by system. 
* **Messages** client packets that system received.
* **Rule Engine** message flow and device event processing and response metrics statistics.

The main dashboard `EMQX` shows the overall metrics of the EMQX cluster, you can switch data source and cluster by the global variable at the top of the dashboard.  

In addition, some panels have its own set of links that are shown in the upper left corner of the panel, the links will jump to relevant dashboard in a new tab, you can filter the metrics based on global variables, such as `cluster`,`node`, and so on.    
