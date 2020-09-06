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

// func TestBeat321(t *testing.T) {
// 	beatDif(3, 2, 1)
// }

// func TestBeat121(t *testing.T) {
// 	beatDif(1, 2, 1)
// }

// func TestBeat111(t *testing.T) {
// 	beatDif(1, 1, 1)
// }

// func TestBeat120(t *testing.T) {
// 	beatDif(1, 2, 0)
// }

// func TestBeat110(t *testing.T) {
// 	beatDif(1, 1, 0)
// }

func beatDif(bi, si, stopi int) {
	b := resource.NewBeatWithInterval("test beat", bi, func(t resource.Tick) {
		fmt.Printf("%v\n", t)
	})
	b.Start()
	time.Sleep(time.Duration(si) * time.Second)
	b.Stop()
	time.Sleep(time.Duration(stopi) * time.Second)
}
