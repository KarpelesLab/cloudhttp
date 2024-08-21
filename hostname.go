package cloudhttp

import (
	"encoding/base32"
	"net"
)

var b32e = base32.NewEncoding("abcdefghijklmnopqrstuvwxyz234567").WithPadding(base32.NoPadding)

// HostnameForIP returns a hostname for the given ipv4
func HostnameForIP(ip net.IP) string {
	if ip == nil {
		return ""
	}
	ip = ip.To4()
	if ip == nil {
		return ""
	}
	return b32e.EncodeToString(ip[:]) + ".g-dns.net"
}
