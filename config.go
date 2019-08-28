package main

type tlsWrapConfig struct {
	remote        string
	listen        string
	serverName    string
	rootCAPath    string
	clientCrtPath string
	clientKeyPath string
	proxy         string
}
