package main

import (
	"flag"
	"github.com/terloo/nfs-csi/driver"
)

var (
	endpoints = flag.String("endpoint", "", "CSI endpoint")
)

func main() {
	flag.Parse()
	server := driver.NewNonBlockGRPCServer()
	server.RunNodeServer(*endpoints)
}
