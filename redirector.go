package cloudhttp

import (
	"net"
	"net/http"
	"net/url"
)

// HttpsRedirector is a http.Handler that will redirect all requests it receives to the
// specified host as https.
type HttpsRedirector string

func (h HttpsRedirector) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// take request's URL, rewrite scheme and host
	u := &url.URL{}
	*u = *req.URL
	u.Scheme = "https"
	u.Host = string(h)

	http.Redirect(rw, req, u.String(), http.StatusMovedPermanently)
}

// HttpsRedirectorServer launches a redirector and serves serves requests in a separate goroutine.
func HttpsRedirectorServer(proto, ip, targetHost string) error {
	l, err := net.Listen(proto, ip)
	if err != nil {
		return err
	}

	go http.Serve(l, HttpsRedirector(targetHost))
	return nil
}
