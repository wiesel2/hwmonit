package resource

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
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
	// RTDISK: "disk",  // not supported
	// RTNET:  "net",   // not supported
	RTSHM: "shm",
	RTPRO: "process",
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
	LastResult *ResourceResult
	RrcType    ResourceType
	C          Collector
}

func AllResource() map[ResourceType]string {
	return rtNameMap
}

func rtToName(rt ResourceType) (string, error) {
	name, find := rtNameMap[rt]
	if find == true {
		return name, nil
	}
	return "", errors.New(fmt.Sprintf("Not find rt: %s", rt))
}

type resbeat struct {
	R *Resource
	B *Beat
}

func (rb *resbeat) Name() string {
	return rb.R.Name
}

type ResourceManager struct {
	rbs   map[string]*resbeat
	lock  sync.Mutex
	state int
}

func (rm *ResourceManager) Start() {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	if rm.state == rmStateRun || rm.state == rmStateStop {
		return
	}
	for _, rb := range rm.rbs {
		go func(bb *Beat) { bb.Tick() }(rb.B)
	}
}

func NewResourceManager() *ResourceManager {
	rm := &ResourceManager{
		rbs:   make(map[string]*resbeat),
		state: rmStateInit,
	}
	return rm
}

func (rm *ResourceManager) AddResource(r *Resource, hb *Beat) {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	// bind callback
	hb.Callback = func(t Tick) {
		rr, err := r.C.GetInfo()
		if err != nil {
			return
		}
		r.LastResult = rr
	}
	rm.rbs[r.Name] = &resbeat{
		R: r,
		B: hb,
	}
}

func (rm *ResourceManager) Stop() {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	wg := sync.WaitGroup{}
	for _, v := range rm.rbs {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			v.B.Stop()
			wg.Done()
		}(&wg)
	}
	// wait all close
	wg.Wait()
}

func (rm *ResourceManager) GetAll() map[string]interface{} {
	res := make(map[string]interface{})
	for n, v := range rm.rbs {
		if v.R.LastResult != nil {
			b, err := json.Marshal(*v.R.LastResult)
			if err != nil {
				continue
			}
			res[n] = b
		}
	}
	return res
}

func ResourceBuilder(name string) *Resource {
	name = strings.ToLower(name)
	var t ResourceType
	var collector Collector
	switch name {
	case "cpu":
		t = RTCPU
		collector = &CPU{}
	case "mem":
		t = RTMEM
		collector = &MEM{}
	case "shm":
		t = RTSHM
		collector = &Shm{}
	case "process":
		t = RTPRO
		collector = &Process{}
	default:
		return nil
	}
	r := Resource{
		Name:    name,
		RrcType: t,
	}
	r.C = collector
	return &r
}

// not finished
type DISK struct{}

// not finished
type NET struct{}
