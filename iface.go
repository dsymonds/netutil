package netutil

import (
	"fmt"
	"net"
)

// RestrictAddrToInterface takes a host:port and rewrites it such that listening on
// the returned address will only accept connections on the named interface.
//
// This is useful for serving a network service on a more private interface,
// such as tailscale0.
func RestrictAddrToInterface(hostPort, ifaceName string) (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("getting network interfaces: %w", err)
	}
	var addrs []net.Addr
	for _, iface := range ifaces {
		if iface.Name != ifaceName {
			continue
		}
		addrs, err = iface.Addrs()
		if err != nil {
			return "", fmt.Errorf("getting network addresses for interface %q: %w", iface.Name, err)
		}
		break
	}
	if addrs == nil {
		return "", fmt.Errorf("unknown or address-free network interface %q", ifaceName)
	}
	var ip net.IP
	for _, a := range addrs {
		ipn, ok := a.(*net.IPNet)
		if !ok {
			continue
		}
		ip = ipn.IP.To4() // pick out the IPv4 address
		if ip != nil {
			break
		}
	}
	if ip == nil {
		return "", fmt.Errorf("network interface %q does not have any IPv4 addresses", ifaceName)
	}

	_, port, err := net.SplitHostPort(hostPort)
	if err != nil {
		return "", fmt.Errorf("splitting %q: %w", hostPort, err)
	}
	return net.JoinHostPort(ip.String(), port), nil
}
