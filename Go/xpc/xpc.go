package xpc

import (
	"fmt"
	"net"
	"time"
)

type Host struct {
	XPHost  string
	XPPort  uint
	Timeout time.Duration
}

func Hello(host Host) error {
	d := net.Dialer{Timeout: host.Timeout}
	conn, err := d.Dial("udp", fmt.Sprintf("%s:%d", host.XPHost, host.XPPort))
	if err != nil {
		return fmt.Errorf("dial XPC host: %w", err)
	}
	defer conn.Close()

	return nil
}
