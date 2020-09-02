package test

import (
	"fmt"
	"hwmonit/resource/cpu"
	"io/ioutil"
	"testing"
)

var fileList = []string{"./cpu1.txt", "./cpu2.txt", "./cpu3.txt", "./cpu4.txt"}

func TestCPU1(t *testing.T) {
	// TOP命令样例 1

	// 	cpu1 := `Mem: 195159944K used, 2820220K free, 4334096K shrd, 3200K buff, 134593368K cached
	// CPU:   6% usr   1% sys   0% nic  91% idle   0% io   0% irq   0% sirq
	// Load average: 3.31 3.85 4.07 5/7990 76
	//   PID  PPID USER     STAT   VSZ %VSZ CPU %CPU COMMAND`

	// 	// TOP命令样例 2

	// 	cpu2 := `top - 16:11:49 up 132 days, 18 min,  0 users,  load average: 6.51, 6.45, 6.66
	// Tasks:  22 total,   1 running,  21 sleeping,   0 stopped,   0 zombie
	// %Cpu(s): 11.3 us,  3.6 sy,  0.0 ni, 84.9 id,  0.0 wa,  0.0 hi,  0.2 si,  0.0 st
	// KiB Mem : 13191990+total,  6356936 free, 45617728 used, 79945240 buff/cache
	// KiB Swap: 16777212 total, 16777212 free,        0 used. 80394496 avail Mem

	//   PID USER      PR  NI    VIRT    RES    SHR S  %CPU %MEM     TIME+ COMMAND`

	// 	// TOP命令样例 3

	// 	cpu3 := `Mem: 194877864K used, 3102300K free, 4252184K shrd, 3200K buff, 134472480K cached
	// CPU:   9% usr   5% sys   0% nic  85% idle   0% io   0% irq   0% sirq
	// Load average: 5.04 5.06 4.50 8/8001 32
	//   PID  PPID USER     STAT   VSZ %VSZ CPU %CPU COMMAND`

	// 	// TOP命令样例 4

	// 	cpu4 := `top - 08:39:10 up 132 days, 45 min,  0 users,  load average: 9.56, 9.76, 9.37
	// Tasks:  14 total,   1 running,  13 sleeping,   0 stopped,   0 zombie
	// %Cpu(s): 10.3 us,  3.3 sy,  0.0 ni, 86.2 id,  0.1 wa,  0.0 hi,  0.2 si,  0.0 st
	// KiB Mem:  19798016+total, 18841305+used,  9567116 free,     2912 buffers
	// KiB Swap: 16777212 total,  1744300 used, 15032912 free. 12639403+cached Mem`

	for _, k := range fileList {
		b, err := ioutil.ReadFile(k)
		if err != nil {
			continue
		}
		r, err := cpu.ParseTOP(&b, `cpu[\(]?[s]?[\)]?:.*`, `(\d+[.]?[\d+]?)[%]?\s+`,
			map[string]string{
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
			})
		if err != nil {
			fmt.Print("ERROR")
		}
		fmt.Printf("%v", r)
	}
}
