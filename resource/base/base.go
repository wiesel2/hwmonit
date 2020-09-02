package base

import (
	"errors"
	"fmt"
	"net"
	"time"
)

type HostInfo struct {
	Addr net.IPAddr
	Name string
}

type ResourceResult struct {
	Name      string            `json:"name"`
	Timestamp time.Time         `json:"timestamp"`
	Result    map[string]string `json:"result"`
}

func NewResourceResult(n string, r map[string]string) *ResourceResult {
	return &ResourceResult{
		Timestamp: time.Now(),
		Name:      n,
		Result:    r,
	}
}

type ResourceType int

const (
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

type Collector interface {
	GetInfo() (*ResourceResult, error)
}

type Resource struct {
	Name string
	// ResChan    chan [20]ResourceResult
	LastResult *ResourceResult
	RrcType    ResourceType
	C          Collector
}

func AllResource() map[ResourceType]string {
	return rtNameMap
}

func RtToName(rt ResourceType) (string, error) {
	name, find := rtNameMap[rt]
	if find == true {
		return name, nil
	}
	return "", errors.New(fmt.Sprintf("Not find rt: %s", rt))
}

// not finished
type DISK struct{}

// not finished
type NET struct{}
