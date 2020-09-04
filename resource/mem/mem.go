package mem

import (
	"hwmonit/common"
	"hwmonit/resource/base"
	"regexp"
	"strconv"
	"strings"
	"unsafe"
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

func (m *MEM) GetInfo() (*base.ResourceResult, error) {
	cmdRes, err := common.ExecSysCmd(5, "cat", "/proc/meminfo")
	if err != nil {
		return nil, err
	}
	data := ParseInfo(cmdRes, memMap)

	n, _ := base.RtToName(base.RTMEM)
	return base.NewResourceResult(n, data), nil
}

type Swap struct{}

var swapMap = map[string]string{
	// swap
	"SwapTotal": "total",
	"SwapFree":  "free",
	// used = SwapTotal - SwapFree
	"SwapCached:": "cached",
}

func (s *Swap) GetInfo() (*base.ResourceResult, error) {
	cmdRes, err := common.ExecSysCmd(5, "cat", "/proc/meminfo")
	if err != nil {
		return nil, err
	}
	data := ParseInfo(cmdRes, swapMap)
	n, _ := base.RtToName(base.RTMEM)
	return base.NewResourceResult(n, data), nil

}

func ParseInfo(mem *[]byte, m map[string]string) map[string]string {
	res := make(map[string]string)
	memStr := (*string)(unsafe.Pointer(mem))
	memlines := strings.Split(*memStr, "\n")
	for _, line := range memlines {
		lowerStr := strings.ToLower(line)
		for k, v := range m {
			reg := regexp.MustCompile(k + `:\s+(\d+)\skb`)
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

func (shm *Shm) GetInfo() (*base.ResourceResult, error) {
	cmdRes, err := common.ExecSysCmd(5, "df")
	if err != nil {
		return nil, err
	}
	data := ParseSHM(cmdRes)
	n, _ := base.RtToName(base.RTSHM)
	return base.NewResourceResult(n, data), nil

}

func ParseSHM(b *[]byte) map[string]string {
	res := make(map[string]string)
	memlines := strings.Split(*(*string)(unsafe.Pointer(b)), "\n")
	for _, line := range memlines {
		lowerStr := strings.ToLower(line)
		if strings.Index(lowerStr, "shm") == 0 {
			reg := regexp.MustCompile(`\s+`)
			spr := reg.Split(lowerStr, 5)
			if len(spr) < 6 {
				ss := strings.Split(spr[4], " ")
				spr = spr[0:4]
				spr = append(spr, ss...)
			}
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
