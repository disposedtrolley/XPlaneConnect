package xpc

import (
	"bytes"
	"fmt"
	"net"
	"time"

	"github.com/lunixbochs/struc"
)

type Host struct {
	XPHost  string
	XPPort  uint
	Port    uint
	Timeout time.Duration
}

type Conn struct {
	udp     *net.UDPConn
	timeout time.Duration
}

func (c *Conn) Write(in []byte) error {
	_ = c.udp.SetDeadline(time.Now().Add(c.timeout))
	_, err := c.udp.Write(in)
	return err
}

func (c *Conn) Read(out []byte, expectedBytes int) error {
	_ = c.udp.SetDeadline(time.Now().Add(c.timeout))
	n, _, err := c.udp.ReadFromUDP(out)
	if err != nil {
		return err
	}
	if n != expectedBytes {
		return fmt.Errorf("expected response length %d but got %d", expectedBytes, n)
	}
	return nil
}

func Dial(h Host) (*Conn, error) {
	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", h.XPHost, h.XPPort))
	if err != nil {
		return nil, fmt.Errorf("resolve server address: %w", err)
	}

	localAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("localhost:%d", h.Port))
	if err != nil {
		return nil, fmt.Errorf("resolve local address: %w", err)
	}

	conn, err := net.DialUDP("udp", localAddr, serverAddr)
	if err != nil {
		return nil, fmt.Errorf("dial udp: %w", err)
	}

	return &Conn{
		udp:     conn,
		timeout: h.Timeout,
	}, nil
}

type CTRL struct {
	Elevator   float32 `struc:"float32,little"`
	Aileron    float32 `struc:"float32,little"`
	Rudder     float32 `struc:"float32,little"`
	Throttle   float32 `struc:"float32,little"`
	Gear       int     `struc:"int8"`
	Flaps      float32 `struc:"float32,little"`
	Aircraft   uint    `struc:"uint8"`
	Speedbrake float32 `struc:"float32,little"`
}

func GetCTRL(conn *Conn, aircraft uint) (ctrl CTRL, err error) {
	type reqPacked struct {
		Command   string `struc:"[4]uint8,little"`
		Padding__ byte   `struc:"pad"`
		Aircraft  uint   `struc:"uint8"`
	}
	type respPacked struct {
		Header    string `struc:"[4]uint8,little"`
		Padding__ byte   `struc:"pad"`
		CTRL
	}

	var reqBuf bytes.Buffer
	if err := struc.Pack(&reqBuf, &reqPacked{
		Command:  "GETC",
		Aircraft: aircraft,
	}); err != nil {
		return ctrl, fmt.Errorf("pack request: %w", err)
	}

	if err := conn.Write(reqBuf.Bytes()); err != nil {
		return ctrl, fmt.Errorf("send request: %w", err)
	}

	respBuf := make([]byte, 31)
	if err := conn.Read(respBuf, len(respBuf)); err != nil {
		return ctrl, fmt.Errorf("read response: %w", err)
	}

	var r respPacked
	if err := struc.Unpack(bytes.NewReader(respBuf), &r); err != nil {
		return ctrl, fmt.Errorf("unpack response: %w", err)
	}

	return r.CTRL, nil
}

func SendCTRL(conn *Conn, ctrl *CTRL) error {
	type reqPacked struct {
		Command   string `struc:"[4]uint8,little"`
		Padding__ byte   `struc:"pad"`
		CTRL
	}

	var reqBuf bytes.Buffer
	if err := struc.Pack(&reqBuf, &reqPacked{
		Command: "CTRL",
		CTRL:    *ctrl,
	}); err != nil {
		return fmt.Errorf("pack request: %w", err)
	}

	return conn.Write(reqBuf.Bytes())
}
