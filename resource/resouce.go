package resource

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"hwmonit/network"

	"github.com/golang/protobuf/ptypes/empty"
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
		go func(bb *Beat) { bb.tick() }(rb.B)
	}
}

func NewResourceManager() *ResourceManager {
	rm := &ResourceManager{
		rbs:   make(map[string]*resbeat),
		state: rmStateInit,
	}
	return rm
}

func (rm *ResourceManager) AddResource(name string, interval int) {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	r := ResourceBuilder(name)
	b := NewBeatWithInterval(name, int(interval))
	rm.addResource(r, b)

}

func (rm *ResourceManager) addResource(r *Resource, hb *Beat) {

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
			v.B.stop()
			wg.Done()
		}(&wg)
	}
	// wait all close
	wg.Wait()
}

func getResourceResult(r *Resource) (*ResourceResult, error) {
	if r.LastResult != nil {
		return r.LastResult, nil
	}
	return nil, fmt.Errorf("Resource %s is result is empty", r.Name)
}

// Export
func (rm *ResourceManager) GetAllR() map[string]interface{} {
	res := make(map[string]interface{})
	for n, v := range rm.rbs {
		rr, err := getResourceResult(v.R)
		if err != nil {
			continue
		}
		res[n] = rr
	}
	return res
}

// Export for RPC/http
func (rm *ResourceManager) GetAll(ctx context.Context, in *empty.Empty) (*network.Resp, error) {
	res, _ := json.Marshal(rm.GetAllR())
	var suc int32 = 0
	var msg string
	var err error
	if len(res) == 0 && len(rm.rbs) != 0 {
		suc = 1
		msg = "resource get failed."
		err = errors.New(msg)
	}
	return &network.Resp{
		State:   suc,
		Data:    string(res),
		Message: msg,
	}, err
}

func (rm *ResourceManager) GetR(name string) map[string]interface{} {
	res := make(map[string]interface{})
	v, ok := rm.rbs[name]
	if ok == true {
		b, err := getResourceResult(v.R)
		if err == nil {
			res[name] = b
		}
	}
	return res
}

func (rm *ResourceManager) Get(ctx context.Context, in *network.GetReq) (*network.Resp, error) {
	res, _ := json.Marshal(rm.GetR(in.Name))
	var suc int32 = 0
	var msg string
	var err error
	if len(res) == 0 && len(rm.rbs) != 0 {
		suc = 1
		msg = "resource get failed."
		err = errors.New(msg)
	}
	return &network.Resp{
		State:   suc,
		Data:    string(res),
		Message: msg,
	}, err
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
