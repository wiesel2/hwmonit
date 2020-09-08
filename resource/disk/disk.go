package disk

import (
	"hwmonit/common"
	"hwmonit/resource/base"
	"regexp"
	"strings"
	"unsafe"
)

var logger = common.GetLogger()

// Disk export
type Disk struct{}

// GetInfo export, implementation of Collector interface
func (d *Disk) GetInfo() (*base.ResourceResult, error) {
	cmdRes, err := common.ExecSysCmd(5, "df")
	if err != nil {
		return nil, err
	}
	data := ParseDF(cmdRes)
	return base.NewResourceResult(base.RTDISK, data)
}

// ParseDF export
func ParseDF(b *[]byte) map[string]interface{} {
	res := make(map[string]interface{})
	memlines := strings.Split(*(*string)(unsafe.Pointer(b)), "\n")
	for _, line := range memlines {
		lowerStr := strings.ToLower(line)
		if strings.Index(lowerStr, "filesystem") == 0 {
			// skip header
			continue
		}
		reg := regexp.MustCompile(`\s+`)
		spr := reg.Split(lowerStr, 5)
		if len(spr) < 5 {
			continue
		}
		if len(spr) < 6 {
			ss := reg.Split(spr[4], -1)
			spr = spr[0:4]
			spr = append(spr, ss...)
		}
		tmp := map[string]string{
			"filesystem": spr[0],
			"blocks":     spr[1],
			"used":       spr[2],
			"available":  spr[3],
			"use":        spr[4],
			"mounted":    spr[len(spr)-1],
		}
		res[spr[len(spr)-1]] = tmp
	}
	return res
}
