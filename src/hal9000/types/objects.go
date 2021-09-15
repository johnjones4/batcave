package types

import "time"

type Nameable interface {
	GetNames() []string
}

type Person interface {
	GetID() string
	GetNames() []string
	GetPhoneNumber() string
	GetOriginName() string
}

type Device interface {
	GetNames() []string
	GetID() string
	GetType() string
	GetDevices(runtime Runtime) []Device
}

type Event interface {
	GetStartTime() time.Time
	GetEndTime() time.Time
	GetName() string
}

type Displayable interface {
	GetNames() []string
	GetURL() string
	GetType() string
	GetSource() string
}

type Job interface {
	GetID() string
	GetInterval() time.Duration
	GetName() string
}

type JobState string

const (
	JobStateNormal   = "normal"
	JobStateAbnormal = "abnormal"
	JobStateLate     = "late"
)

type JobStatusInfo struct {
	State       JobState `json:"state"`
	Description string   `json:"description"`
}
type JobStatus struct {
	Info       JobStatusInfo `json:"info"`
	LastUpdate time.Time     `json:"lastUpdate"`
}
