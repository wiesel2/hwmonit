package resource

import (
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

type systype int
type ostype int

var once sync.Once

const (
	Normal systype = iota // 物理机
	Docker                // docker
)

var SysType systype

func checksystype() {
	once.Do(func() {
		// 读取CPU文件，proc/self/cgroup，并判断是否有docker字段
		cpufile := "/proc/self/cgroup"
		cf, err := os.Stat(cpufile)
		if os.IsExist(err) == false || cf.IsDir() {
			return
		}
		// 读取mount文件，proc/self/mountinfo，并判断是否有docker字段
		mountfile := "/proc/self/mountinfo"
		mf, err := os.Stat(mountfile)
		if os.IsExist(err) == false || mf.IsDir() {
			return
		}
		cbytes, _ := ioutil.ReadFile(cpufile)
		suc := 0
		if strings.Contains(string(cbytes), "docker") {
			suc++
		}
		mbybtes, _ := ioutil.ReadFile(mountfile)
		if strings.Contains(string(mbybtes), "docker") {
			suc++
		}
		if suc == 2 {
			SysType = Docker
		}
	})
}

func GetSysType() systype {
	checksystype()
	return SysType
}

func IsOnDocker() bool {
	checksystype()
	if SysType == Docker {
		return true
	}
	return false
}
