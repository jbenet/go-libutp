// Copyright 2014 Juan Batiz-Benet.  All rights reserved.
// Copyright 2009 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utp

import (
  "net"
  "strconv"
)

// ipEmptyString is like ip.String except that it returns
// an empty string when ip is unset.
// Unexported item from golang/pkg/net/ip.go
func ipEmptyString(ip net.IP) string {
	if len(ip) == 0 {
		return ""
	}
	return ip.String()
}

// This file implements UTPAddr and UTPConn,
// to match net.UDPAddr and net.UDPConn
type UTPConn struct {
  context *UTPContext
  socket *UTPSocket
}

// UTPAddr represents the address of a uTP end point.
type UTPAddr struct {
	IP   net.IP
	Port int
	Zone string // IPv6 scoped addressing zone
}

// Network returns the address's network name, "utp".
func (a *UTPAddr) Network() string { return "utp" }

func (a *UTPAddr) String() string {
	if a == nil {
		return "<nil>"
	}
	ip := ipEmptyString(a.IP)
	if a.Zone != "" {
		return net.JoinHostPort(ip+"%"+a.Zone, strconv.Itoa(a.Port))
	}
	return net.JoinHostPort(ip, strconv.Itoa(a.Port))
}

func (a *UTPAddr) toAddr() net.Addr {
	if a == nil {
		return nil
	}
	return a
}


func newUTPConn(ctx *UTPContext, sock *UTPSocket) *UTPConn {
  return &UTPConn{
    context: ctx,
    socket: sock,
  }
}


/*func DialUTP(net string, laddr, raddr *UTPAddr) (*UTPConn, error) {

}

func (c *TCPConn) Close() error
    func (c *TCPConn) CloseRead() error
    func (c *TCPConn) CloseWrite() error
    func (c *TCPConn) File() (f *os.File, err error)
    func (c *TCPConn) LocalAddr() Addr
    func (c *TCPConn) Read(b []byte) (int, error)
    func (c *TCPConn) ReadFrom(r io.Reader) (int64, error)
    func (c *TCPConn) RemoteAddr() Addr
    func (c *TCPConn) SetDeadline(t time.Time) error
    func (c *TCPConn) SetKeepAlive(keepalive bool) error
    func (c *TCPConn) SetKeepAlivePeriod(d time.Duration) error
    func (c *TCPConn) SetLinger(sec int) error
    func (c *TCPConn) SetNoDelay(noDelay bool) error
    func (c *TCPConn) SetReadBuffer(bytes int) error
    func (c *TCPConn) SetReadDeadline(t time.Time) error
    func (c *TCPConn) SetWriteBuffer(bytes int) error
    func (c *TCPConn) SetWriteDeadline(t time.Time) error
    func (c *TCPConn) Write(b []byte) (int, error)*/
