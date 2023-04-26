package collector

type LicenseInfo struct {
	MaxClientLimit int64
	Expiration     int64
	RemainingDays  float64
}

type ClusterStatus struct {
	Status     int
	NodeUptime map[string]int64
	NodeMaxFDs map[string]int
	CPULoads   map[string]CPULoad
}

type CPULoad struct {
	Load1  float64
	Load5  float64
	Load15 float64
}

type Broker struct {
	MsgConsumeTimeCosts map[string]uint64
	MsgInputPeriodSec   int64
	MsgOutputPeriodSec  int64
}

type RuleEngine struct {
	// NodeName the name of emqx node
	NodeName string
	RuleID   string
	// TopicHitCount
	TopicHitCount      int64
	ExecPassCount      int64
	ExecFailureCount   int64
	NoResultCount      int64
	ExecRate           float64
	ExecLast5mRate     float64
	ExecMaxRate        float64
	ActionTotal        int64
	ActionSuccess      int64
	ActionFailed       int64
	ActionExecTimeCost map[string]uint64
}

type DataBridge struct {
	Type string
	Name string
	// Status define the status of the third-party resource. It's ok if the value is 2, else is not ready
	Status int

	// bridge Metrics
	Queuing    int64
	RateLast5m float64
	RateMax    float64
	Failed     int64
	Dropped    int64
}

type Authentication struct {
	// NodeName the name of emqx node
	NodeName       string
	ResType        string
	Total          int64
	AllowCount     int64
	DenyCount      int64
	ExecRate       float64
	ExecLast5mRate float64
	ExecMaxRate    float64
	ExecTimeCost   map[string]uint64
}

type Authorization struct {
	// NodeName the name of emqx node
	NodeName       string
	ResType        string
	Total          int64
	AllowCount     int64
	DenyCount      int64
	ExecRate       float64
	ExecLast5mRate float64
	ExecMaxRate    float64
	ExecTimeCost   map[string]uint64
}

type DataSource struct {
	ResType string
	Status  int
}

type Cluster interface {
	GetLicense() (*LicenseInfo, error)
	GetClusterStatus() (ClusterStatus, error)
	GetBrokerMetrics() (*Broker, error)
	GetRuleEngineMetrics() ([]DataBridge, []RuleEngine, error)
	GetAuthenticationMetrics() ([]DataSource, []Authentication, error)
	GetAuthorizationMetrics() ([]DataSource, []Authorization, error)
}
