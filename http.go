package cloudhttp

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/netip"
	"os"
	"strconv"

	"github.com/KarpelesLab/cloudinfo"
	"golang.org/x/crypto/acme/autocert"
)

// ServerHost is filled when calling Serve() and will contain the hostname that was generated for this machine
var ServerHost string

// Serve is a simple function that will listen on http/https ports, and redirect all
// http requests to the right https host.
func Serve(h http.Handler) error {
	info, err := cloudinfo.Load()
	if err != nil {
		log.Printf("failed to load cloud info: %s", err)
		// setup local info
		info = cloudinfo.Sysinfo()
	}
	ip, found := info.PublicIP.GetFirstV4()
	if !found {
		log.Printf("no public IPv4, using localhost...")
		ip = netip.AddrFrom4([...]byte{127, 0, 0, 1})
	}

	ServerHost = b32e.EncodeToString(ip.AsSlice()) + ".g-dns.net"

	l, err := getSocketFallback(443)
	if err != nil {
		return fmt.Errorf("failed to listen on https: %w", err)
	}

	m := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache(os.TempDir()),
		HostPolicy: autocert.HostWhitelist(ServerHost),
	}
	cfg := m.TLSConfig()
	l = tls.NewListener(l, cfg)
	go http.Serve(l, h)

	log.Printf("cloudhttp: listening on %s for https connections on host %s", l.Addr(), ServerHost)

	l, err = getSocketFallback(80)
	if err == nil {
		go http.Serve(l, HttpsRedirector(ServerHost))
		log.Printf("cloudhttp: listening on %s for http connections", l.Addr())
	}
	return nil
}

// getSocketFallback will attempt to get a socket at the given port
func getSocketFallback(port int) (net.Listener, error) {
	jump := 8000
	c := 0

	for {
		l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
		if err == nil {
			return l, err
		}
		c += 1
		if c > 10 {
			return nil, err
		}
		port += jump
		jump = 1
	}
}
