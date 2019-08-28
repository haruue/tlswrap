package main

import (
	"crypto/tls"
	"golang.org/x/net/proxy"
	"log"
	"net"
)

type tlsWrapContext struct {
	wrapConfig *tlsWrapConfig
	tlsConfig  tls.Config
	dialer     proxy.Dialer
}

func (c *tlsWrapContext) start() {
	listener, err := net.Listen("tcp", c.wrapConfig.listen)
	if err != nil {
		log.Fatalf("fatal: can't listen on %s\n", c.wrapConfig.listen)
	}
	log.Printf("info: listen on %s, hit CTRL-C to stop\n", c.wrapConfig.listen)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("error: accept() failed: %v\n", err)
			continue
		}
		go c.handleConn(conn)
	}
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

func (c *tlsWrapContext) handleConn(conn net.Conn) {
	log.Printf("info: accepted %s\n", conn.RemoteAddr().String())
	rconn, err := c.dialer.Dial("tcp", c.wrapConfig.remote)
	if err != nil {
		log.Printf("error: failed to dial the remote %s\n", c.wrapConfig.remote)
		_ = conn.Close()
		return
	}
	tlsConn := tls.Client(rconn, &c.tlsConfig)
	go pipe(conn, tlsConn)
	go pipe(tlsConn, conn)
}
