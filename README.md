tlswrap
==========
[简体中文](README-zh_cn.md)

A TLS termination proxy intended to solve the "last mile" problem for the
[TLS client authentication](https://blog.cloudflare.com/introducing-tls-client-auth/).


## Overview

The feature of this tool is connect to TLS client authentication protected
service, and expose the server locally for other client that leaks of the
support of TLS client authentication.

I write this to bring modern TLS support to some legacy http client, such like
WebDAV.


## Build & Install

```
go get github.com/haruue/tlswrap
```

or download pre-built binary at
[GitHub Release](https://github.com/haruue/tlswrap/releases/latest)

## Usage

Run `tlswrap -help` for the usage.

Example:

```shell
./tlswrap \
        -ca ca.crt \
        -cert client.crt \
        -key client.key \
        -listen 127.0.0.1:2333 \
        -remote ssl.example.com:22443 \
        -sni secret.example.com \
        -host secret.example.com:22443 \
        -mode http
```


## Security Problem

1. tlswrap doesn't support encrypted key yet, a non-encrypted key is more like
   to be stolen.
2. if you expose the http server to a non-loopback network, your service will
   be access by everyone without authentication.


## License

MIT


