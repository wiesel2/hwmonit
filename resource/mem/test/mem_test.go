package test

import (
	"fmt"
	"hwmonit/resource/mem"
	"io/ioutil"
	"testing"
)

func TestMem(t *testing.T) {

	b, err := ioutil.ReadFile("mem1.txt")
	if err != nil {
		panic("File not exited")
	}
	var memMap = map[string]string{
		"memtotal": "total",
		"memfree":  "free",
		// used = MemTotal - MemFree
		"buffers": "buffers",
	}
	res := mem.ParseInfo(&b, memMap)
	fmt.Printf("%v", res)

	var swapMap = map[string]string{
		// swap
		"swaptotal": "total",
		"swapfree":  "free",
		// used = SwapTotal - SwapFree
		"swapcached:": "cached",
	}

	reswap := mem.ParseInfo(&b, swapMap)
	fmt.Printf("%v", reswap)
}

// func TestShm(t *testing.T) {
// 	b, err := ioutil.ReadFile("shm1.txt")
// 	if err != nil {
// 		panic("File not exited")
// 	}

// 	res := mem.ParseSHM(&b)
// 	fmt.Printf("%v", res)

// }
