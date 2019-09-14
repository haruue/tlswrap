tlswrap
==========
[English](README.md)

一个旨在解决 [TLS 客户端认证](https://blog.cloudflare.com/introducing-tls-client-auth/)
中 「最后一公里」 问题的 TLS 终端代理。


## 概览

这个工具的主要功能是与被 TLS 客户端认证保护的服务器通信，
并将对应的服务在本地暴露成一个未加密的 HTTP 服务。

我编写这个工具来为那些不支持 TLS 客户端认证的老式应用提供现代的 TLS 支持。
例如各种 WebDAV 客户端。


## 编译和安装

```
go get github.com/haruue/tlswrap
```

或者在
[GitHub Release](https://github.com/haruue/tlswrap/releases/latest)
下载预编译的二进制可执行文件。

## 用法

运行 `tlswrap -help` 以获取用法说明：

```
Usage of ./tlswrap:
  -ca string
        验证服务端证书所使用的 CA， 默认会使用系统内置的。
  -cert string
        客户端认证使用的证书文件路径。
  -key string
        客户端认证使用的私钥文件路径。
  -listen string
        addr:port 在本地监听的地址和端口。
  -remote string
        addr:port 服务器所在的远程地址和端口。
  -sni string
        服务器名称指示。 这个名称会在 Client Hello 时使用以让服务器决定发送哪个证书。
        如果 -host 没有指定， 那么这个名称也会被用于 "Host:" 字段中。
  -host string
        服务名称， 用于 HTTP 头中的 "Host" 字段， 默认使用与 -sni 相同的名称。
  -dns string
        addr:port 使用的 DNS 服务器。
                  在 Android 上， 必须指定此选项， 否则解析将无法工作。
  -mode string
        http/tunnel 两种模式， http 模式会重写 HTTP 头中各字段的内容， 不支持长连接。
                    tunnel 字段在建立连接之后即将数据原样传输， 支持长连接。
                    默认 http 模式。
```

一个完整的例子：

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


## 安全问题

1. tlswrap 暂时不支持加密的密钥， 而存储解密的密钥看起来更容易导致密钥被窃取。

2. 如果你 -listen 在非回环网络上，
   那么任何能连接到你的计算机的人都能在无须授权的情况下访问你的服务。


## 开源许可

MIT


