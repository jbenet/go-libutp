package utp

/*
#cgo CFLAGS: -I libutp
#cgo LDFLAGS: -lutp -L libutp
#include <stdlib.h>
#include <utp.h>
*/
import "C"

import (
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

  return int(C.utp_process_udp(c.raw, (*C.byte)(unsafe.Pointer(&buf[0])),
    C.size_t(len), (*C.struct_sockaddr)(to), C.socklen_t(tolen)))
}

func (c *UTPContext) CheckTimeouts() {
  C.utp_check_timeouts(c.raw)
}

func (c *UTPContext) IssueDeferredAcks() {
  C.utp_issue_deferred_acks(c.raw)
}

type ContextStats struct {
  recv [5]uint32
  send [5]uint32
}

func (c *UTPContext) GetContextStats() *ContextStats {
  cs := C.utp_get_context_stats(c.raw)
  return &ContextStats{
    recv: [5]uint32(cs._nraw_recv),
    send: [5]uint32(cs._nraw_send),
  }
}

func (c *UTPContext) CreateSocket() *UTPSocket {
  return &UTPSocket{
    raw: C.utp_create_socket(c.raw),
    lock: sync.Mutex{},
  }
}

// UTPSocket type, tracks individual UTP connections
type UTPSocket struct {
  raw *C.utp_socket
  lock sync.Mutex
}

// void*			utp_set_userdata				(utp_socket *s, void *userdata);
// void*			utp_get_userdata				(utp_socket *s);

/*int				utp_setsockopt					(utp_socket *s, int opt, int val);
int				utp_getsockopt					(utp_socket *s, int opt);
int				utp_connect						(utp_socket *s, const struct sockaddr *to, socklen_t tolen);
ssize_t			utp_write						(utp_socket *s, void *buf, size_t count);
ssize_t			utp_writev						(utp_socket *s, struct utp_iovec *iovec, size_t num_iovecs);
int				utp_getpeername					(utp_socket *s, struct sockaddr *addr, socklen_t *addrlen);
void			utp_read_drained				(utp_socket *s);
int				utp_get_delays					(utp_socket *s, uint32 *ours, uint32 *theirs, uint32 *age);
utp_socket_stats* utp_get_stats					(utp_socket *s);
utp_context*	utp_get_context					(utp_socket *s);
void			utp_close						(utp_socket *s);*/
