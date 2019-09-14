package main

import (
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

type tlsWrapContext struct {
	wrapConfig *tlsWrapConfig
	tlsConfig  tls.Config
}

func (c *tlsWrapContext) start() {
	listener, err := net.Listen("tcp", c.wrapConfig.listen)
	if err != nil {
		log.Fatalf("fatal: can't listen on %s\n", c.wrapConfig.listen)
	}
	log.Printf("info: listen on %s, hit CTRL-C to stop\n", c.wrapConfig.listen)

	if c.wrapConfig.mode == modeTunnel {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("error: accept() failed: %v\n", err)
				continue
			}
			go c.handleConn(conn)
		}
	} else {
		log.Fatal(http.Serve(listener, c))
	}
}

func (c *tlsWrapContext) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.handleHTTP(w, r)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func (c *tlsWrapContext) filterHeaders(headers http.Header, src, dst string) {
	for i, vi := range headers {
		for j, vj := range vi {
			headers[i][j] = strings.ReplaceAll(vj, src, dst)
		}
	}
}

func (c *tlsWrapContext) handleHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	tr := http.Transport{
		TLSClientConfig:        &context.tlsConfig,
		MaxIdleConns:       20,
		IdleConnTimeout:    60 * time.Second,
	}

	r.Host = c.wrapConfig.host
	r.URL.Scheme = "https"
	r.URL.Host = c.wrapConfig.remote
	if c.wrapConfig.host != "" {
		c.filterHeaders(r.Header, c.wrapConfig.listen, c.wrapConfig.host)
	}

	resp, err := tr.RoundTrip(r)

	if err != nil {
		w.WriteHeader(503)
		log.Printf("error: failed to request the remote %s: %v\n", c.wrapConfig.remote, err)
		return
	}

	if c.wrapConfig.host != "" {
		c.filterHeaders(resp.Header, c.wrapConfig.host, c.wrapConfig.listen)
	}
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}

func (c *tlsWrapContext) handleConn(conn net.Conn) {
	log.Printf("info: accepted %s\n", conn.RemoteAddr().String())
	rconn, err := net.Dial("tcp", c.wrapConfig.remote)
	if err != nil {
		log.Printf("error: failed to dial the remote %s\n", c.wrapConfig.remote)
		_ = conn.Close()
		return
	}
	tlsConn := tls.Client(rconn, &c.tlsConfig)
	go pipe(conn, tlsConn)
	go pipe(tlsConn, conn)
}

func pipe(src, dst net.Conn) {
	defer func() {
		_ = src.Close()
		_ = dst.Close()
	}()
	buff := make([]byte, 65535)
	for {
		n, err := src.Read(buff)
		if err != nil {
			log.Printf("error: failed to read from %s for %s: %v\n", src.RemoteAddr().String(), dst.RemoteAddr().String(), err)
			return
		}
		b := buff[:n]

		_, err = dst.Write(b)
		if err != nil {
			log.Printf("error: failed to write to %s for %s: %v\n", dst.RemoteAddr().String(), src.RemoteAddr().String(), err)
			return
		}
	}
}
