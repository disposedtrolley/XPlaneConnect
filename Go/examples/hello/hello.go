package main

import (
	"fmt"
	"time"
	"xpc"
)

func main() {
	timeout, err := time.ParseDuration("100ms")
	if err != nil {
		panic(err)
	}

	conn, err := xpc.Dial(xpc.Host{
		XPHost:  "localhost",
		XPPort:  49009,
		Timeout: timeout,
	})

	for {
		ctrl, err := xpc.GetCTRL(conn, 0)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Elevator: %.4f	Aileron: %.4f	Rudder: %.4f\n", ctrl.Elevator, ctrl.Aileron, ctrl.Rudder)

		time.Sleep(100 * time.Millisecond)
	}
}
