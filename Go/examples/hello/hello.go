package main

import (
	"log"
	"time"
	"xpc"
)

func main() {
	timeout, err := time.ParseDuration("100ms")
	if err != nil {
		panic(err)
	}

	h := xpc.Host{
		XPHost:  "localhost",
		XPPort:  49009,
		Timeout: timeout,
	}

	if err := xpc.Hello(h); err != nil {
		log.Fatal(err)
	}
}
