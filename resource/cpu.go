package resource

import (
	"errors"
	"regexp"
	"strings"
)

type CPU struct {
}

func (c *CPU) GetInfo() (*ResourceResult, error) {

	// limitation
	// current cpu load average, 1, 5, 10

	cmdRes, err := execSysCmd(5, "top", "-n", "1")
	if err != nil {
		return nil, err
	}

	data, err := parseTOP(cmdRes, `cpu[\(]?[s]?[\)]?:.*`, `(\d+[.]?[\d+]?)[%]?\s+`)
	// TODO
	if IsOnDocker() {
		// 获取cpu占用时间
		// 获取cpu分配的周期
		// 从docker cgroup 文件读取
	} else {
	}
	n, _ := rtToName(RTCPU)
	return NewResourceResult(n, data), nil
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

func getCPUMap() map[string]string {
	return cupMap
}

func parseTOP(cpuStr, filter_str, key string) (map[string]string, error) {
	//
	// output {
	//	"us": "12.6", "sy": "2.6", "ni": "0.0", "id": "84.4", "wa": "0.2", "hi": "0.0", "si": "0.3", "st": "0.0"
	// }

	res := make(map[string]string)
	cpulines := strings.Split(cpuStr, "\n")
	for _, line := range cpulines {
		lowerStr := strings.ToLower(line)
		m, err := regexp.MatchString(filter_str, lowerStr)
		if err != nil {
			continue
		}
		if m {
			for k, v := range getCPUMap() {
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

// Export for
func ParseTOPCPU(s string) (map[string]string, error) {
	return parseTOP(s, `cpu[\(]?[s]?[\)]?:.*`, `(\d+[.]?[\d+]?)[%]?\s+`)
}
