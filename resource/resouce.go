package resource

import (
	"errors"
	"fmt"
	"sync"
)

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
	RTMEM:  "memory",
	RTSWAP: "swap",
	RTDISK: "disk",
	RTNET:  "net",
	RTSHM:  "shm",
	RTPRO:  "process",
}

const (
	rmStateInit int = iota
	rmStateRun
	rmStateStop
)

type Collector interface {
	GetInfo() (*ResourceResult, error)
}

type Resource struct {
	Name string
	// ResChan    chan [20]ResourceResult
	LastResult ResourceResult
	RrcType    ResourceType
	c          Collector
}

func rtToName(rt ResourceType) (string, error) {
	name, find := rtNameMap[rt]
	if find == true {
		return name, nil
	}
	return "", errors.New(fmt.Sprintf("Not find rt: %s", rt))
}

func (r *Resource) Run() {
}

type ResourceManager struct {
	rs    map[string]*Resource
	beats map[string]*Beat
	lock  sync.Mutex
	state int
}

func (rm *ResourceManager) Start() {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	if rm.state == rmStateRun || rm.state == rmStateStop {
		return
	}
	for _, b := range rm.beats {
		b.Run()
		go func(bb *Beat{bb.Tick()}(b)
	}
}

func NewResourceManager() *ResourceManager {
	rm := &ResourceManager{
		beats: make(map[string]*Beat),
		rs:    make(map[string]*Resource),
		state: rmStateInit,
	}
	return rm
}

func (rm *ResourceManager) AddResource(r *Resource, hb *Beat) {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	rm.rs[r.Name] = r
	rm.beats[r.Name] = hb
}

func (rm *ResourceManager) Stop() {
	rm.lock.Lock()
	defer rm.lock.Unlock()

}

// not finished
type DISK struct{}

// not finished
type NET struct{}
