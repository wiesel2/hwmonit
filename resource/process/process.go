package process

import (
	"errors"
	"fmt"
	"hwmonit/common"
	"hwmonit/resource/base"
	"hwmonit/resource/cpu"
	"runtime"
)

type Process struct{}

var processMap = map[string]string{
	"total":    "total",
	"running":  "running",
	"sleeping": "sleeping",
	"stopped":  "stopped",
	"zombie":   "zombie",
}

func (m *Process) GetInfo() (*base.ResourceResult, error) {

	var cmdRes *[]byte
	var err error
	var filter string
	var key string
	switch common.GetOSType() {
	case common.OSWin:
		err = errors.New("Windows not supported")
	case common.OSLinux:
		cmdRes, err = common.ExecSysCmd(5, "top", "-n", "1")
		filter = `^tasks:.*`
		key = `(\d+)`
	case common.OSMac:
		cmdRes, err = common.ExecSysCmd(5, "top", "|", "head", "-20")
		filter = `processes:`
		key = `(\d+)\s*`
	default:
		err = fmt.Errorf("Other %s not supported ", runtime.GOOS)
	}

	if err != nil {
		return nil, err
	}
	data, _ := cpu.ParseTOP(cmdRes, filter, key, processMap)

	n, _ := base.RtToName(base.RTPRO)
	return base.NewResourceResult(n, data), nil

}
