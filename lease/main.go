package main

import (
	"flag"
	"os"

	"github.com/k0kubun/pp"
	"github.com/thewraven/dhcpd-leases"
)

func main() {
	fn := flag.String("file", "leases", "file to be parsed")
	flag.Parse()
	f, err := os.Open(*fn)
	if err != nil {
		panic(err)
	}
	leases, err := leases.ParseLeases(f)
	if err != nil {
		panic(err)
	}
	pp.Println(leases)
}
