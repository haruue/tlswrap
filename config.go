package main

type tlsWrapConfig struct {
	mode          string
	remote        string
	listen        string
	serverName    string
	host          string
	rootCAPath    string
	clientCrtPath string
	clientKeyPath string
}

const (
	modeHttp = "http"
	modeTunnel = "tunnel"
)