package cpu

import (
	"errors"
	"fmt"
	"hwmonit/common"
	base "hwmonit/resource/base"
	"regexp"
	"runtime"
	"strings"
	"unsafe"
)

type CPU struct {
}

func (c *CPU) GetInfo() (*base.ResourceResult, error) {

	// limitation
	// current cpu load average, 1, 5, 10
	var cmdRes *[]byte
	var err error
	var filter string
	var key string
	switch common.GetOSType() {
	case common.OSWin:
		err = errors.New("Windows not supported")
	case common.OSLinux:
		cmdRes, err = common.ExecSysCmd(5, "top", "-n", "1")
		filter = `cpu[\(]?[s]?[\)]?:.*`
		key = `(\d+[.]?[\d+]?)[%]?\s+`
	case common.OSMac:
		cmdRes, err = common.ExecSysCmd(5, "top", "-l", "1")
		filter = `cpu usage:`
		key = `(\d+[.]?[\d+]?)[%]?\s+`
	default:
		err = fmt.Errorf("Other %s not supported ", runtime.GOOS)
	}
	if err != nil {
		return nil, err
	}

	data, err := ParseTOP(cmdRes, filter, key, cupMap)
	// TODO
	if common.IsOnDocker() {
		// 获取cpu占用时间
		// 获取cpu分配的周期
		// 从docker cgroup 文件读取
	} else {
	}
	n, _ := base.RtToName(base.RTCPU)
	return base.NewResourceResult(n, data), nil
}

var cupMap = map[string]string{
	"usr":  "us", // user
	"us":   "us",
	"sys":  "sy", // system
	"sy":   "sy",
	"nic":  "ni", // nice
	"ni":   "ni",
	"idle": "id", // idel
	"id":   "id",
	"wa":   "wa", // wait io
	"hi":   "hi", // headware interupt
	"irq":  "hi",
	"si":   "si", // software interupt
	"sirq": "si",
	"st":   "st", // steal

}

func ParseTOP(s *[]byte, filter_str, key string, keymap map[string]string) (map[string]string, error) {
	//
	// output {
	//	"us": "12.6", "sy": "2.6", "ni": "0.0", "id": "84.4", "wa": "0.2", "hi": "0.0", "si": "0.3", "st": "0.0"
	// }

	res := make(map[string]string)
	cpuStr := *(*string)(unsafe.Pointer(s))
	cpulines := strings.Split(cpuStr, "\n")
	for _, line := range cpulines {
		lowerStr := strings.ToLower(line)
		m, err := regexp.MatchString(filter_str, lowerStr)
		if err != nil {
			continue
		}
		if m {
			for k, v := range keymap {
				if _, e := res[v]; e == true {
					// skip existed
					continue
				}
				reg := regexp.MustCompile(key + k)
				finds := reg.FindAllStringSubmatch(lowerStr, -1)
				if len(finds) == 0 {
					continue
				}
				res[v] = finds[0][1]
			}
		}
	}

	if len(res) == 0 {
		return res, errors.New("Get CPU  info empty")
	}

	return res, nil
}
