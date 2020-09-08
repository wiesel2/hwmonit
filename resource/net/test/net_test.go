package test

import (
	"fmt"
	"hwmonit/resource/net"
	"io/ioutil"
	"testing"
)

func TestNet(t *testing.T) {
	b, err := ioutil.ReadFile("net.txt")
	if err != nil {
		panic("File not exited")
	}
	res, err := net.ParseNetInfo(&b)
	fmt.Printf("%v", res)
}
