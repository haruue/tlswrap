package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"log"
	"os"
)

func init() {
	// enable tls1.3
	_ = os.Setenv("GODEBUG", os.Getenv("GODEBUG")+",tls13=1")
}

var context tlsWrapContext

func main() {
	config := tlsWrapConfig{}
	flag.StringVar(&config.mode, "mode", "http", "http/tunnel The \"http\" mode will modify the HTTP headers to specified host, " +
		"\"tunnel\" mode will open a long connection over TLS and transmit data as is")
	flag.StringVar(&config.listen, "listen", "", "addr:port Local addr and port to listen")
	flag.StringVar(&config.remote, "remote", "", "addr:port Remote server to connect")
	flag.StringVar(&config.serverName, "sni", "", "Server name indication")
	flag.StringVar(&config.host, "host", "", "Host name of the server used in the HTTP header, \"-sni\" will be used when unspecified, valid in \"http\" mode only")
	flag.StringVar(&config.rootCAPath, "ca", "", "CA certificate to verify peer against")
	flag.StringVar(&config.clientCrtPath, "cert", "", "Client certificate file")
	flag.StringVar(&config.clientKeyPath, "key", "", "Private key file name")
	flag.Parse()

	context.wrapConfig = &config

	if config.mode != modeHttp && config.mode != modeTunnel {
		log.Fatalf("fatal: unknown mode: %s\n", config.mode)
	}
	if config.remote == "" {
		log.Fatalf("fatal: no remote to connect\n")
	}
	if config.listen == "" {
		log.Fatalf("fatal: no local addr to listen\n")
	}
	if config.host == "" {
		config.host = config.serverName
	} else {
		if config.mode == "tunnel" {
			log.Printf("warn: -host %s is ignored in tunnel mode, please use -sni if you meant to specify the server name indication\n", config.host)
		}
	}

	context.tlsConfig.ServerName = config.serverName

	if config.rootCAPath != "" {
		bs, err := ioutil.ReadFile(config.rootCAPath)
		if err != nil {
			log.Fatalf("fatal: can't open CA certificate path %s to read: %v\n", config.rootCAPath, err)
		}
		rootCA := x509.NewCertPool()
		if ok := rootCA.AppendCertsFromPEM(bs); !ok {
			log.Fatalf("fatal: failed to add root CA certificate from %s, is it in X.509 PEM format?\n", config.rootCAPath)
		}
		context.tlsConfig.RootCAs = rootCA
	}

	if config.clientCrtPath != "" && config.clientKeyPath != "" {
		var certificates []tls.Certificate
		keyPair, err := tls.LoadX509KeyPair(config.clientCrtPath, config.clientKeyPath)
		if err != nil {
			log.Fatalf("fatal: failed to add client cert-key pair from %s and %s, are they in X.509 PEM format?\n",
				config.clientCrtPath, config.clientKeyPath)
		}
		certificates = append(certificates, keyPair)
		context.tlsConfig.Certificates = certificates
	} else {
		if config.clientCrtPath != "" {
			log.Fatalf("fatal: client certificate specified without a key\n")
		}
		if config.clientKeyPath != "" {
			log.Fatalf("fatal: client key specified without a certificate\n")
		}
	}

	context.start()
}
