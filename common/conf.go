package common

import (
	"flag"
	"fmt"
	"os"
	"sync"
)

// Configs
//  port: local listen port
//

var config map[string]string
var lock sync.Mutex

// ["key:default"]
var configList = []string{
	"port",
	"log",
	"log-dir",
	"mem",
	"disk",
	"net",
	"cpu",
	"process",
}

func readConfig() {
	lock.Lock()
	defer lock.Unlock()
	config = make(map[string]string)
	for _, v := range configList {

		osVal := os.Getenv(v)
		if osVal == "" {
			panic(fmt.Sprintf("Config not set: %s", v))
		}
		config[v] = osVal
	}
}

//GetConfig export
func GetConfig() map[string]string {
	return config
}

var (
	h bool   // help
	p string // port
	l string // log name
	d string // log dir
	v string = "0.1"
	m string // mem
	s string // disk
	c string // cpu
	r string // process
	w string // swap
	n string // net
)

func setbool(T bool) string {
	if T == true {
		return "true"
	}
	return "false"
}

func readArgs() {
	flag.BoolVar(&h, "h", false, "help")

	flag.StringVar(&m, "m", "30", "memory stat collect time")
	flag.StringVar(&s, "s", "30", "disk stat collect time")
	flag.StringVar(&c, "c", "30", "cpu stat collect time")
	flag.StringVar(&r, "r", "30", "process stat collect time")
	flag.StringVar(&w, "w", "30", "swap stat collect time")
	flag.StringVar(&n, "n", "30", "net stat collect time")

	// common config
	flag.StringVar(&p, "p", "10241", "port")
	flag.StringVar(&l, "l", "hwmonit.log", "log file name")
	flag.StringVar(&d, "d", "/tmp", "log dir")
	flag.Parse()
	setEnv("port", p)
	setEnv("log", l)
	setEnv("net", w)
	setEnv("swap", w)
	setEnv("process", r)
	setEnv("cpu", c)
	setEnv("disk", s)
	setEnv("mem", m)
	setEnv("log-dir", d)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `hwmonit. 
Usage: ./hwmonit [-h] [-t] [-p 10241] [-l logname] [-d logdir] [-m] [-s] [-c] [-w] 
Options:
`)
		flag.PrintDefaults()
	}
	// force update config
	readConfig()

	if h == true {
		flag.Usage()
		os.Exit(0)
	}
}

func setEnv(k, v string) {
	if v != "" {
		os.Setenv(k, v)
	}
}
