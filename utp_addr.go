package utp

import (
  "net"
  "strings"
  "strconv"
  "syscall"
)

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

func SockaddrToUTP(sa syscall.Sockaddr) net.Addr {
	switch sa := sa.(type) {
	case *syscall.SockaddrInet4:
		return &UTPAddr{IP: sa.Addr[0:], Port: sa.Port}
	case *syscall.SockaddrInet6:
		return &UTPAddr{IP: sa.Addr[0:], Port: sa.Port, Zone: zoneToString(int(sa.ZoneId))}
	}
	return nil
}

// Convert UTPAddr into a sockaddr.
func (a *UTPAddr) Sockaddr() (syscall.Sockaddr, error) {
	if a == nil {
		return nil, nil
	}

  family := syscall.AF_INET
  if len(a.IP) > 4 {
    family = syscall.AF_INET6
  }
	return IpToSockaddr(family, a.IP, a.Port, a.Zone)
}

func (a *UTPAddr) family() int {
	if a == nil || len(a.IP) <= net.IPv4len {
		return syscall.AF_INET
	}
	if a.IP.To4() != nil {
		return syscall.AF_INET
	}
	return syscall.AF_INET6
}

func (a *UTPAddr) isWildcard() bool {
	if a == nil || a.IP == nil {
		return true
	}
	return a.IP.IsUnspecified()
}

// ResolveUTPAddr parses addr as a UTP address of the form "host:port"
// or "[ipv6-host%zone]:port" and resolves a pair of domain name and
// port name on the network net, which must be "utp", "utp4" or
// "utp6".  A literal address or host name for IPv6 must be enclosed
// in square brackets, as in "[::1]:80", "[ipv6-host]:http" or
// "[ipv6-host%zone]:80".
func ResolveUTPAddr(net_, addr string) (*UTPAddr, error) {
  switch net_ {
  case "utp", "utp4", "utp6":
  case "":
    net_ = "utp"
  default:
    return nil, net.UnknownNetworkError(net_)
  }

  net_ = strings.Replace(net_, "d", "t", 1)
  a, err := net.ResolveUDPAddr(net_, addr)
  if err != nil {
    return nil, err
  }
  return &UTPAddr{
    IP: a.IP,
    Port: a.Port,
    Zone: a.Zone,
  }, nil
}


// Unexported IP Address functions copied from pkg/net.
//
// You would think that in a language as mature as go one wouldn't have to
// copy paste large portions of the core libraries. Yes, you would think that.

// Turn an IP Address into a
func IpToSockaddr(family int, ip net.IP, port int, zone string) (syscall.Sockaddr, error) {
	switch family {
	case syscall.AF_INET:
		if len(ip) == 0 {
			ip = net.IPv4zero
		}
		if ip = ip.To4(); ip == nil {
			return nil, net.InvalidAddrError("non-IPv4 address")
		}
		sa := new(syscall.SockaddrInet4)
		for i := 0; i < net.IPv4len; i++ {
			sa.Addr[i] = ip[i]
		}
		sa.Port = port
		return sa, nil
	case syscall.AF_INET6:
		if len(ip) == 0 {
			ip = net.IPv6zero
		}
		// IPv4 callers use 0.0.0.0 to mean "announce on any available address".
		// In IPv6 mode, Linux treats that as meaning "announce on 0.0.0.0",
		// which it refuses to do.  Rewrite to the IPv6 unspecified address.
		if ip.Equal(net.IPv4zero) {
			ip = net.IPv6zero
		}
		if ip = ip.To16(); ip == nil {
			return nil, net.InvalidAddrError("non-IPv6 address")
		}
		sa := new(syscall.SockaddrInet6)
		for i := 0; i < net.IPv6len; i++ {
			sa.Addr[i] = ip[i]
		}
		sa.Port = port
		sa.ZoneId = uint32(zoneToInt(zone))
		return sa, nil
	}
	return nil, net.InvalidAddrError("unexpected socket family")
}

func zoneToString(zone int) string {
	if zone == 0 {
		return ""
	}
	if ifi, err := net.InterfaceByIndex(zone); err == nil {
		return ifi.Name
	}
	return itod(uint(zone))
}

func zoneToInt(zone string) int {
	if zone == "" {
		return 0
	}
	if ifi, err := net.InterfaceByName(zone); err == nil {
		return ifi.Index
	}
	n, _, _ := dtoi(zone, 0)
	return n
}

// ipEmptyString is like ip.String except that it returns
// an empty string when ip is unset.
// Unexported item from golang/pkg/net/ip.go
func ipEmptyString(ip net.IP) string {
	if len(ip) == 0 {
		return ""
	}
	return ip.String()
}

// Convert i to decimal string.
func itod(i uint) string {
	if i == 0 {
		return "0"
	}

	// Assemble decimal in reverse order.
	var b [32]byte
	bp := len(b)
	for ; i > 0; i /= 10 {
		bp--
		b[bp] = byte(i%10) + '0'
	}

	return string(b[bp:])
}

// Bigger than we need, not too big to worry about overflow
const big = 0xFFFFFF

// Decimal to integer starting at &s[i0].
// Returns number, new offset, success.
func dtoi(s string, i0 int) (n int, i int, ok bool) {
	n = 0
	for i = i0; i < len(s) && '0' <= s[i] && s[i] <= '9'; i++ {
		n = n*10 + int(s[i]-'0')
		if n >= big {
			return 0, i, false
		}
	}
	if i == i0 {
		return 0, i, false
	}
	return n, i, true
}

// Unexported parts of pkg syscall
