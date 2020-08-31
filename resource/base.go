package resource

import (
	"encoding/json"
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
		Result:    make(map[string]string),
	}
}

type ResourceInfo interface {
	String() (string, error)
	Dump() ([]byte, error)
}

func (rr *ResourceResult) Dump() ([]byte, error) {
	r, err := json.Marshal(rr)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (rr *ResourceResult) String() (string, error) {
	r, err := rr.Dump()
	if err != nil {
		return "", err
	}
	return string(r), nil
}
