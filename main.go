package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

var (
	proxy   string
	dest    string
	port    string
	verbose bool
)

func init() {
	if proxy == "" {
		// sets defaults if manually (no make) built with `go build`
		flag.StringVar(&proxy, "proxy", "http://localhost:8181", "Upstream proxy URI")
	} else {
		// if built with `make`, proxy will be prepopulated
		flag.StringVar(&proxy, "proxy", proxy, "Upstream proxy URI")
	}

	if dest == "" {
		flag.StringVar(&dest, "dest", "https://example.com", "Final destination URL")
	} else {
		flag.StringVar(&dest, "dest", dest, "Final destination URL")
	}

	if port == "" {
		flag.StringVar(&port, "port", "8080", "Listen port")
	} else {
		flag.StringVar(&port, "port", port, "Listen port")
	}
	flag.BoolVar(&verbose, "v", false, "Verbose")
	flag.Parse()

	port = fmt.Sprintf(":%s", port)
}

func main() {
	origin, err := url.Parse(dest)
	if err != nil {
		log.Fatal(err)
	}

	transport := &http.Transport{
		Proxy: func(*http.Request) (*url.URL, error) { return url.Parse(proxy) },
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		ExpectContinueTimeout: 1 * time.Second,
	}

	director := func(req *http.Request) {
		// req.Header.Add("X-Forwarded-Host", req.Host)
		// req.Header.Add("X-Origin-Host", origin.Host)
		req.URL.Scheme = "https"
		req.URL.Host = origin.Host
	}

	reverseProxy := &httputil.ReverseProxy{
		Transport: transport,
		Director:  director,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		reverseProxy.ServeHTTP(w, r)
	})

	mlog("Listening on %s", port)
	mlog("Upstream proxy: %s", proxy)
	mlog("Forwarding to: %s", dest)

	log.Fatal(http.ListenAndServe(port, nil))
}

func mlog(line, f string) {
	if verbose {
		log.Printf(line, f)
	}
}
