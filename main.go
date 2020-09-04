package main

import (
	"hwmonit/common"
	"hwmonit/network"
	"hwmonit/resource"
)

var logger = common.GetLogger()

func main() {
	rm := resource.NewManager()
	rm.LoadAllResources(common.GetConfig())
	rm.Start()
	config := common.GetConfig()
	s := network.Server{}
	if err := s.Serve(":"+config["port"], rm); err != nil {
		rm.Stop()
	}
}
