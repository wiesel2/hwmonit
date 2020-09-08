package net

import (
	"errors"
	"hwmonit/common"
	"hwmonit/resource/base"
	"regexp"
	"strings"
	"unsafe"
)

// Net export
type Net struct{}

// GetInfo export, a implementation for Collector interface
func (n *Net) GetInfo() (*base.ResourceResult, error) {
	var fName string
	if common.IsOnDocker() {
		// read data from /proc/self/net/dev
		fName = "/proc/self/net/dev"
	} else {
		// read data from /proc/net/dev
		fName = "/proc/net/dev"
	}

	res, err := common.ExecSysCmd(5, "cat", fName)
	if err != nil {
		return nil, err
	}
	data, err := ParseNetInfo(res)
	if err != nil {
		return nil, err
	}
	return base.NewResourceResult(base.RTNET, data)
}

// ParseNetInfo export
func ParseNetInfo(s *[]byte) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	netStr := *(*string)(unsafe.Pointer(s))
	netlines := strings.Split(netStr, "\n")
	for _, line := range netlines {
		lowerStr := strings.ToLower(line)
		regF := regexp.MustCompile(`(\w+):.*`)
		m := regF.FindAllStringSubmatch(lowerStr, -1)
		if len(m) == 0 {
			// Not matched data row
			continue
		}
		reg := regexp.MustCompile(`(\d+)`)
		m2 := reg.FindAllStringSubmatch(lowerStr, -1)
		if len(m2) < 16 {
			continue
		}
		var tmp = map[string]string{
			"receive-bytes":        m2[0][0],
			"receive-packets":      m2[1][0],
			"receive-errs":         m2[2][0],
			"receive-drop":         m2[3][0],
			"receive-fifo":         m2[4][0],
			"receive-frame":        m2[5][0],
			"receive-compressed":   m2[6][0],
			"receive-multicast":    m2[7][0],
			"transimit-bytes":      m2[8][0],
			"transimit-packets":    m2[9][0],
			"transimit-errs":       m2[10][0],
			"transimit-drop":       m2[11][0],
			"transimit-fifo":       m2[12][0],
			"transimit-colls":      m2[13][0],
			"transimit-carrier":    m2[14][0],
			"transimit-compressed": m2[15][0],
		}
		res[m[0][1]] = tmp
	}
	if len(res) > 0 {
		return res, nil
	}
	return nil, errors.New("Empty net data parsed")

}
