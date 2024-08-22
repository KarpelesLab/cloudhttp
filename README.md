[![GoDoc](https://godoc.org/github.com/KarpelesLab/cloudhttp?status.svg)](https://godoc.org/github.com/KarpelesLab/cloudhttp)

# cloudhttp

Easily setup a HTTPs server on most cloud providers (aws, gcp, etc) without having to worry about much of anything.

Each server will automatically get a dynamic .g-dns.net domain for the SSL certificate.

## Usage

```go
mux := http.NewServeMux()
// setup mux handlers
err := cloudhttp.Serve(mux)
if err != nil {
    // ...
}
```

With just this code, a http/https server will start and will have a valid SSL certificate, provided the host's IP address is publicly accessible and port 443 isn't blocked from outside.
