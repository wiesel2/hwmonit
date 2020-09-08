package resource

import (
	"context"
	"errors"
	"fmt"
	"hwmonit/common"
	"hwmonit/network"
	"hwmonit/resource/base"
	"hwmonit/resource/cpu"
	"hwmonit/resource/disk"
	"hwmonit/resource/mem"
	"hwmonit/resource/net"
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
		go func(wg *sync.WaitGroup, rb *resbeat) {
			rb.B.Stop()
			wg.Done()
		}(&wg, v)
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

func transResourceResult(rr *base.ResourceResult) *network.Results {
	contents := make([]*network.Content, 0, 10)
	for _, v := range rr.Result {
		switch i := v.(type) {
		case map[string]string:
			contents = append(contents, &network.Content{Content: i})
		case string:

			tmp := make(map[string]string)
			for k1, v1 := range rr.Result {
				tmp[k1] = v1.(string)
			}
			contents = append(contents, &network.Content{Content: tmp})
			return contentMaker(rr, contents)
		}
	}
	return contentMaker(rr, contents)
}

func contentMaker(rr *base.ResourceResult, contents []*network.Content) *network.Results {
	nr := &network.Results{
		Name:      rr.Name,
		Timestamp: rr.Timestamp.String(),
		Results:   contents,
	}
	return nr
}

// get all resource static
func (rm *Manager) getAllR() []*network.Results {
	res := make([]*network.Results, 0, 10)
	for _, v := range rm.rbs {
		rr, err := getResourceResult(v.R)
		if err != nil {
			continue
		}
		ress := transResourceResult(rr)
		res = append(res, ress)
	}

	return res
}

// GetAll Export
// Export for http and gRPC invoke
func (rm *Manager) GetAll(ctx context.Context, in *empty.Empty) (*network.Resp, error) {
	res := rm.getAllR()
	var suc int32 = 0
	var msg string
	var err error
	if len(res) == 0 && len(rm.rbs) != 0 {
		suc = 1
		msg = "resource get failed."
		err = errors.New(msg)
	}
	// data []*network.Results =
	return &network.Resp{
		State:   suc,
		Data:    res,
		Message: msg,
	}, err
}

func (rm *Manager) getR(name string) []*network.Results {
	res := make([]*network.Results, 0, 10)
	v, ok := rm.rbs[name]
	if ok == true {
		rr, err := getResourceResult(v.R)
		if err == nil {
			res = append(res, transResourceResult(rr))
		}
	}
	return res
}

// Get Export
// Export for http and gRPC invoke
func (rm *Manager) Get(ctx context.Context, in *network.GetReq) (*network.Resp, error) {
	res := rm.getR(in.Name)
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
		Data:    res,
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
	case "disk":
		t = base.RTDISK
		collector = &disk.Disk{}
	case "process":
		t = base.RTPRO
		collector = &process.Process{}
	case "net":
		t = base.RTNET
		collector = &net.Net{}
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
