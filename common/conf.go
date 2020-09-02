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
	"tls",
	"mem",
	"shm",
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
	t bool   // tls enabled
	p string // port
	l string // log name
	d string // log dir
	v string = "0.1"
	m string // mem
	s string // shm
	c string // cpu
	w string // process
)

func setbool(T bool) string {
	if T == true {
		return "true"
	}
	return "false"
}

func readArgs() {
	flag.BoolVar(&h, "h", false, "help")
	flag.BoolVar(&t, "t", false, "enable tls")
	setEnv("tls", setbool(t))
	flag.StringVar(&m, "m", "0", "memory stat")
	setEnv("mem", m)
	flag.StringVar(&s, "s", "0", "shm tls")
	setEnv("shm", s)
	flag.StringVar(&c, "c", "30", "cpu stat")
	setEnv("cpu", c)
	flag.StringVar(&w, "w", "0", "process stat")
	setEnv("process", w)

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
