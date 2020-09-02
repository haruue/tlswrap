你好！
很冒昧用这样的方式来和你沟通，如有打扰请忽略我的提交哈。我是光年实验室（gnlab.com）的HR，在招Golang开发工程师，我们是一个技术型团队，技术氛围非常好。全职和兼职都可以，不过最好是全职，工作地点杭州。
我们公司是做流量增长的，Golang负责开发SAAS平台的应用，我们做的很多应用是全新的，工作非常有挑战也很有意思，是国内很多大厂的顾问。
如果有兴趣的话加我微信：13515810775  ，也可以访问 https://gnlab.com/，联系客服转发给HR。
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
        -dns 1.1.1.1:53 \
        -mode http
```


## Security Problem

1. tlswrap doesn't support encrypted key yet, a non-encrypted key is more like
   to be stolen.
2. if you expose the http server to a non-loopback network, your service will
   be access by everyone without authentication.


## License

MIT


