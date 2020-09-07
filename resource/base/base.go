package base

import (
	"fmt"
	"net"
	"time"
)

// HostInfo export
type HostInfo struct {
	Addr net.IPAddr
	Name string
}

// ResourceResult export
type ResourceResult struct {
	Name      string            `json:"name"`
	Timestamp time.Time         `json:"timestamp"`
	Result    map[string]string `json:"result"`
}

// NewResourceResult export
func NewResourceResult(t ResourceType, r map[string]string) (*ResourceResult, error) {
	n, err := rtToName(t)
	if err != nil {
		return nil, err
	}
	return &ResourceResult{
		Timestamp: time.Now(),
		Name:      n,
		Result:    r,
	}, nil
}

// ResourceType export
type ResourceType int

const (
	// RTCPU Export
	RTCPU ResourceType = iota // CPU
	RTMEM
	RTSWAP
	RTSHM
	RTDISK
	RTNET
	RTPRO
)

var rtNameMap = map[ResourceType]string{
	RTCPU:  "cpu",
	RTMEM:  "mem",
	RTSWAP: "swap",
	// RTDISK: "disk",  // not supported
	// RTNET:  "net",   // not supported
	RTSHM: "shm",
	RTPRO: "process",
}

// Collector export, interface of all resource
// GetInfo to gather info of resource
type Collector interface {
	GetInfo() (*ResourceResult, error)
}

// Resource Export
type Resource struct {
	Name string
	// ResChan    chan [20]ResourceResult
	LastResult *ResourceResult
	RrcType    ResourceType
	C          Collector
}

// AllResource export
func AllResource() map[ResourceType]string {
	return rtNameMap
}

func rtToName(rt ResourceType) (string, error) {
	name, find := rtNameMap[rt]
	if find == true {
		return name, nil
	}
	return "", fmt.Errorf("Not find rt: %d", rt)
}

// not finished
type DISK struct{}

// not finished
type NET struct{}
