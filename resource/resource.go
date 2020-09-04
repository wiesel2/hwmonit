package resource

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hwmonit/common"
	"hwmonit/network"
	"hwmonit/resource/base"
	"hwmonit/resource/cpu"
	"hwmonit/resource/mem"
	"hwmonit/resource/process"
	"strconv"
	"strings"
	"sync"

	"github.com/golang/protobuf/ptypes/empty"
)

var logger = common.GetLogger()

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

// Manager Export
type Manager struct {
	rbs   map[string]*resbeat
	lock  sync.Mutex
	state int
}

// Start export
// Start, non-blocking
func (rm *Manager) Start() {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	if rm.state == rmStateRun || rm.state == rmStateStop {
		return
	}
	for _, rb := range rm.rbs {
		rb.B.Start()
	}
	rm.state = rmStateRun
}

// NewManager Export
// Creator of *Manager
func NewManager() *Manager {
	rm := &Manager{
		rbs:   make(map[string]*resbeat),
		state: rmStateInit,
	}
	return rm
}

// LoadAllResources Export
// map: "source name": "interval". interval must be greater than 0!
/* sample {
	"cpu": "30",	// every 30 seconds
	"mem": "20",	// every 20 seconds
	"shm" : "dddd"	// invalid, disable shm
	"process": "-1"	// < 0, disbale process
}
*/
func (rm *Manager) LoadAllResources(config map[string]string) {
	for _, name := range base.AllResource() {
		v1, ok := config[name]
		if ok == false {
			continue
		}
		interval, err := strconv.ParseInt(v1, 10, 10)
		if err != nil {
			continue
		}

		if interval > 0 {
			rm.addResource(name, int(interval))
		}
	}
}

// add create a new resrouce and bind a heat to peroid collect hareware info(GetInfo)
func (rm *Manager) addResource(name string, interval int) {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	r := resourceBuilder(name)
	b := NewBeatWithInterval(name, int(interval), func(t Tick) {
		rr, err := r.C.GetInfo()
		if err != nil {
			return
		}
		r.LastResult = rr
	})
	rm.rbs[r.Name] = &resbeat{
		R: r,
		B: b,
	}
}

// Stop Export
// Stop all source beat, blocked until all beat stopped.
func (rm *Manager) Stop() {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	rm.state = rmStateStop
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
func getResourceResult(r *base.Resource) (*base.ResourceResult, error) {
	if r.LastResult != nil {
		return r.LastResult, nil
	}
	return nil, fmt.Errorf("Resource %s is result is empty", r.Name)
}

// get all resource static
func (rm *Manager) getAllR() map[string]interface{} {
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

// GetAll Export
// Export for http and gRPC invoke
func (rm *Manager) GetAll(ctx context.Context, in *empty.Empty) (*network.Resp, error) {
	res, _ := json.Marshal(rm.getAllR())
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

func (rm *Manager) getR(name string) map[string]interface{} {
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

// Get Export
// Export for http and gRPC invoke
func (rm *Manager) Get(ctx context.Context, in *network.GetReq) (*network.Resp, error) {
	res, _ := json.Marshal(rm.getR(in.Name))
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

// resource factory to build CPU, MEM, SWAP, SHM and process etc.
func resourceBuilder(name string) *base.Resource {
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
	case "swap":
		t = base.RTSWAP
		collector = &mem.Swap{}
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
