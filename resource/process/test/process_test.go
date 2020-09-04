package test

import (
	"fmt"
	"hwmonit/resource/process"
	"testing"
)

func TestProOnMac(t *testing.T) {
	p := process.Process{}
	res, _ := p.GetInfo()
	fmt.Printf("res,: %v", res)
}
