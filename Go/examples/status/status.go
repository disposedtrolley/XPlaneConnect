package main

import (
	"fmt"
	"time"
	"xpc"
)

func getAircraftStatus(conn *xpc.Conn, aircraft uint) {
	ctrl, err := xpc.GetCTRL(conn, aircraft)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", ctrl)

	posi, err := xpc.GetPOSI(conn, aircraft)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", posi)

	values, err := xpc.GetDREFs(conn,
		"sim/flightmodel/misc/h_ind",
		"sim/cockpit2/gauges/indicators/heading_electric_deg_mag_pilot")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n\n", values)
}

func main() {
	conn, err := xpc.Dial(xpc.Host{
		XPHost:  "localhost",
		XPPort:  49009,
		Timeout: time.Millisecond * 50,
	})

	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(100 * time.Millisecond)

	for {
		select {
		case <-ticker.C:
			getAircraftStatus(conn, 0)
		}
	}
}
