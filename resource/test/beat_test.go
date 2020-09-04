package test

import (
	"fmt"
	"hwmonit/resource"
	"testing"
	"time"
)

func TestBeat(t *testing.T) {
	b := resource.NewBeatWithInterval("test beat", 5, func(t resource.Tick) {
		fmt.Printf("%v", t)
	})
	b.Start()
	time.Sleep(10 * time.Second)
	b.Stop()
	// time.Sleep(10 * time.Second)

}
