package utp

/*
#cgo CFLAGS: -I libutp
#cgo LDFLAGS: -lutp -L libutp
#include <stdlib.h>
#include <utp.h>
*/
import "C"

import (
  "github.com/jbenet/go-sockaddr"
  "net"
  "sync"
  "syscall"
  "unsafe"
)

const UTP_VERSION = 2

type UTPContext struct {
  raw *C.utp_context
  lock sync.Mutex
}

func NewUTPContext() *UTPContext {
  return &UTPContext{
    raw: C.utp_init(C.int(UTP_VERSION)),
    lock: sync.Mutex{},
  }
}

func (c *UTPContext) Close() {
  C.utp_destroy(c.raw)
}


// func (c *UTPContext) SetCallback(cb_name int, fn *UTPCallback) {
//   C.utp_set_callback(c.raw, cb_name, func)
// }

// void*			utp_context_set_userdata		(utp_context *ctx, void *userdata);
// void*			utp_context_get_userdata		(utp_context *ctx);

// UTPContext option ids
const (
  UTP_LOG_NORMAL = C.UTP_LOG_NORMAL
  UTP_LOG_MTU    = C.UTP_LOG_MTU
  UTP_LOG_DEBUG  = C.UTP_LOG_DEBUG
  UTP_SNDBUF      = C.UTP_SNDBUF
  UTP_RCVBUF      = C.UTP_RCVBUF
  UTP_TARGET_DELAY = C.UTP_TARGET_DELAY
)

// sets context options
func (c *UTPContext) SetOption(opt, val int) int {
  return int(C.utp_context_set_option(c.raw, C.int(opt), C.int(val)))
}

func (c *UTPContext) GetOption(opt int) int {
  return int(C.utp_context_get_option(c.raw, C.int(opt)))
}

// Processes a UDP packet. For when you're handling UDP yourself.
func (c *UTPContext) ProcessUDP (buf []byte, len int,
    to *syscall.RawSockaddr, tolen int) int {


  bufptr := (*C.byte)(unsafe.Pointer(&buf[0]))
  adrptr := (*C.struct_sockaddr)(unsafe.Pointer(to))
  return int(C.utp_process_udp(c.raw, bufptr, C.size_t(len),
    adrptr, C.socklen_t(tolen)))
}

func (c *UTPContext) CheckTimeouts() {
  C.utp_check_timeouts(c.raw)
}

func (c *UTPContext) IssueDeferredAcks() {
  C.utp_issue_deferred_acks(c.raw)
}

type ContextStats struct {
  Recv [5]uint32
  Send [5]uint32
}

func (c *UTPContext) GetContextStats() *ContextStats {
  ccs := C.utp_get_context_stats(c.raw)
  gcs := &ContextStats{}

  for i := 0; i < 5; i++ {
    gcs.Recv[i] = uint32((*ccs)._nraw_recv[i])
    gcs.Send[i] = uint32((*ccs)._nraw_recv[i])
  }

  return gcs
}

// UTPSocket type, tracks individual UTP connections
type UTPSocket struct {
  ctx *UTPContext
  raw *C.utp_socket
  lock sync.Mutex
}

func NewUTPSocket(c *UTPContext) *UTPSocket {
  return &UTPSocket{
    ctx: c,
    raw: C.utp_create_socket(c.raw),
    lock: sync.Mutex{},
  }
}

// void*			utp_set_userdata				(utp_socket *s, void *userdata);
// void*			utp_get_userdata				(utp_socket *s);


func (s *UTPSocket) SetSockopt(opt, val int) int {
  return int(C.utp_setsockopt(s.raw, C.int(opt), C.int(val)))
}

func (s *UTPSocket) GetSockopt(opt int) int {
  return int(C.utp_getsockopt(s.raw, C.int(opt)))
}

func (s *UTPSocket) Connect(addr *UTPAddr) (int, error) {
// func (s *UTPSocket) Connect(to *syscall.RawSockaddr, tolen int) int {

  if addr == nil {
    return 0, net.InvalidAddrError("No address given.")
  }

  sa, err := addr.Sockaddr()
  if err != nil {
    return 0, err
  }

  rsa, err := sockaddr.NewRawSockaddr(&sa)
  if err != nil {
    return 0, err
  }

  ptr := (*C.struct_sockaddr)(unsafe.Pointer(&rsa.Raw))
  ret := int(C.utp_connect(s.raw, ptr, C.socklen_t(rsa.Len)))
  return ret, nil
}

func (s *UTPSocket) Write(buf []byte, len int) int {
  ptr := unsafe.Pointer(&buf[0])
  return int(C.utp_write(s.raw, ptr, C.size_t(len)))
}

func (s *UTPSocket) Close() {
  C.utp_close(s.raw)
}

func (s *UTPSocket) Context() *UTPContext {
  return s.ctx
}


// int				utp_getpeername					(utp_socket *s, struct sockaddr *addr, socklen_t *addrlen);
// void			utp_read_drained				(utp_socket *s);
// int				utp_get_delays					(utp_socket *s, uint32 *ours, uint32 *theirs, uint32 *age);
// utp_socket_stats* utp_get_stats					(utp_socket *s);
// utp_context*	utp_get_context					(utp_socket *s);
