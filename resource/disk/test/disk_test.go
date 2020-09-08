package test

import (
	"fmt"
	"hwmonit/resource/disk"
	"io/ioutil"
	"testing"
)

func TestMem(t *testing.T) {

	b, err := ioutil.ReadFile("df.txt")
	if err != nil {
		panic("File not exited")
	}

	disinfo := disk.ParseDF(&b)
	fmt.Printf("%v", disinfo)
}
