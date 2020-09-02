package common

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

type systype int
type ostype int

var once3 sync.Once
var once2 sync.Once

const (
	Normal systype = iota // 物理机
	Docker                // docker
)

var SysType systype
var OSType ostype

const (
	OSLinux ostype = iota //OSUnix
	OSMac
	OSWin
)

func checksystype() {
	once2.Do(func() {
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

func ExecSysCmd(timeout int, cmdStr ...string) (*[]byte, error) {
	// create new goroutine to exec cmd
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	var result []byte
	var err error
	cmd := exec.CommandContext(ctx, cmdStr[0], cmdStr[1:]...)
	stdout := &bytes.Buffer{}
	cmd.Stdout = stdout
	if err = cmd.Run(); err != nil {
		return nil, err
	}
	result = stdout.Bytes()

	if err != nil {
		return nil, err
	}
	return &result, err
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

func checkOSType() {
	once3.Do(
		func() {
			osName := runtime.GOOS
			switch osName {
			case "darwin":
				OSType = OSMac
			case "windows":
				OSType = OSWin
			default:
				OSType = OSLinux
			}
		})
}

func GetOSType() ostype {

	checkOSType()
	return OSType

}
