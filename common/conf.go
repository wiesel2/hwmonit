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
	setEnv("mem", m)
	flag.StringVar(&s, "s", "30", "disk stat collect time")
	setEnv("disk", s)
	flag.StringVar(&c, "c", "30", "cpu stat collect time")
	setEnv("cpu", c)
	flag.StringVar(&r, "r", "30", "process stat collect time")
	setEnv("process", r)
	flag.StringVar(&w, "w", "30", "swap stat collect time")
	setEnv("process", w)
	flag.StringVar(&n, "n", "30", "net stat collect time")
	setEnv("net", w)

	// common config
	flag.StringVar(&p, "p", "10241", "port")
	setEnv("port", p)
	flag.StringVar(&l, "l", "hwmonit.log", "log file name")
	setEnv("log", l)
	flag.StringVar(&d, "d", "/tmp", "log dir")
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

	flag.Parse()
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
