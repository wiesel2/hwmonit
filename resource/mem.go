package resource

import (
	"regexp"
	"strconv"
	"strings"
)

// 包含3种类型
//	memory, swap, shm
//
//

type MEM struct{}

var memMap = map[string]string{
	"MemTotal": "total",
	"MemFree":  "free",
	// used = MemTotal - MemFree
	"Buffers": "buffers",
}

type Swap struct{}

var swapMap = map[string]string{
	// swap
	"SwapTotal": "total",
	"SwapFree":  "free",
	// used = SwapTotal - SwapFree
	"SwapCached:": "cached",
}

func (m *MEM) GetInfo() (*ResourceResult, error) {
	cmdRes, err := execSysCmd(5, "cat", "/proc/meminfo")
	if err != nil {
		return nil, err
	}
	data := parseInfo(cmdRes, memMap)

	n, _ := rtToName(RTMEM)
	return NewResourceResult(n, data), nil

}

func (s *Swap) GetInfo() (*ResourceResult, error) {
	cmdRes, err := execSysCmd(5, "cat", "/proc/meminfo")
	if err != nil {
		return nil, err
	}
	data := parseInfo(cmdRes, swapMap)
	n, _ := rtToName(RTMEM)
	return NewResourceResult(n, data), nil

}

func parseInfo(memStr string, m map[string]string) map[string]string {
	res := make(map[string]string)
	memlines := strings.Split(memStr, "\n")
	for _, line := range memlines {
		lowerStr := strings.ToLower(line)
		for k, v := range m {
			reg := regexp.MustCompile(k + `:\s+(\d+)\sKB`)
			finds := reg.FindAllStringSubmatch(lowerStr, -1)
			if len(finds) == 0 {
				continue
			}
			res[v] = finds[0][1]
		}
	}
	t1, te := res["total"]
	f1, ue := res["free"]
	if te == ue == true {
		t2, e1 := strconv.ParseInt(t1, 10, 10)
		f2, e2 := strconv.ParseInt(f1, 10, 10)
		if e1 == nil && e2 == nil {
			res["used"] = string(t2 - f2)
		}
	}
	return res
}

type Shm struct{}

func (shm *Shm) GetInfo() (*ResourceResult, error) {
	cmdRes, err := execSysCmd(5, "df")
	if err != nil {
		return nil, err
	}
	data := parseSHM(cmdRes)
	n, _ := rtToName(RTSHM)
	return NewResourceResult(n, data), nil

}

func parseSHM(s string) map[string]string {
	res := make(map[string]string)
	memlines := strings.Split(s, "\n")
	for _, line := range memlines {
		lowerStr := strings.ToLower(line)
		if strings.Index(lowerStr, "shm") == 0 {
			reg := regexp.MustCompile(`\t`)
			spr := reg.Split(lowerStr, 5)
			res["filesystem"] = spr[0]
			res["blocks"] = spr[1]
			res["used"] = spr[2]
			res["available"] = spr[3]
			res["use"] = spr[4]
			res["mounted"] = spr[5]
			break
		}
	}
	return res
}