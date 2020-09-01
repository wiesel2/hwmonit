package main

import (
	"hwmonit/resource"
	"os"
	"strconv"
)

func main() {
	// config := common.GetConfig()

	rm := resource.NewResourceManager()
	for _, name := range resource.AllResource() {
		v1, ok := os.LookupEnv(name)
		if ok == false {
			continue
		}
		interval, err := strconv.ParseInt(v1, 10, 10)
		if err != nil {
			interval = 30
		}
		if interval >= 0 {
			r := resource.ResourceBuilder(name)
			b := resource.NewBeatWithInterval(name, int(interval), )

		}

	}
}
