package test

import (
	"fmt"
	"hwmonit/resource"
	"testing"
	"time"
)

func TestBeat311(t *testing.T) {
	beatDif(3, 1, 1)
}

func TestBeat321(t *testing.T) {
	beatDif(3, 2, 1)
}

func TestBeat121(t *testing.T) {
	beatDif(1, 2, 1)
}

func TestBeat111(t *testing.T) {
	beatDif(1, 1, 1)
}

func TestBeat120(t *testing.T) {
	beatDif(1, 2, 0)
}

func TestBeat110(t *testing.T) {
	beatDif(1, 1, 0)
}

func beatDif(bi, si, stopi int) {
	name := fmt.Sprintf("test:%d-%d-%d beat", bi, si, stopi)
	b := resource.NewBeatWithInterval(name, bi, func(t resource.Tick) {
		fmt.Printf("%v\n", t)
	})
	b.Start()
	fmt.Printf("%s start \n", name)
	time.Sleep(time.Duration(si) * time.Second)
	fmt.Printf("%s stop \n", name)
	b.Stop()
	fmt.Printf("%s stopped \n", name)
	time.Sleep(time.Duration(stopi) * time.Second)
	fmt.Printf("%s exit \n", name)
}
