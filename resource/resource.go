package resource

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hwmonit/network"
	"hwmonit/resource/base"
	"hwmonit/resource/cpu"
	"hwmonit/resource/mem"
	"hwmonit/resource/process"
	"strings"
	"sync"

	"github.com/golang/protobuf/ptypes/empty"
)

type resbeat struct {
	R *base.Resource
	B *Beat
}

func (rb *resbeat) Name() string {
	return rb.R.Name
}

const (
	rmStateInit int = iota
	rmStateRun
	rmStateStop
)

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

func (rm *ResourceManager) AddResource(name string, interval int) {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	r := ResourceBuilder(name)
	b := NewBeatWithInterval(name, int(interval))
	rm.addResource(r, b)

}

func (rm *ResourceManager) addResource(r *base.Resource, hb *Beat) {

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
func getResourceResult(r *base.Resource) (*base.ResourceResult, error) {
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

func ResourceBuilder(name string) *base.Resource {
	name = strings.ToLower(name)
	var t base.ResourceType
	var collector base.Collector
	switch name {
	case "cpu":
		t = base.RTCPU
		collector = &cpu.CPU{}
	case "mem":
		t = base.RTMEM
		collector = &mem.MEM{}
	case "shm":
		t = base.RTSHM
		collector = &mem.Shm{}
	case "process":
		t = base.RTPRO
		collector = &process.Process{}
	default:
		return nil
	}
	r := base.Resource{
		Name:    name,
		RrcType: t,
	}
	r.C = collector
	return &r
}
