# A QUIC implementation in pure Go

<img src="docs/quic.png" width=303 height=124>

[![Documentation](https://img.shields.io/badge/docs-quic--go.net-red?style=flat)](https://quic-go.net/docs/)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/quic-go/quic-go)](https://pkg.go.dev/github.com/quic-go/quic-go)
[![Code Coverage](https://img.shields.io/codecov/c/github/quic-go/quic-go/master.svg?style=flat-square)](https://codecov.io/gh/quic-go/quic-go/)
[![Fuzzing Status](https://oss-fuzz-build-logs.storage.googleapis.com/badges/quic-go.svg)](https://issues.oss-fuzz.com/issues?q=quic-go)

quic-go is an implementation of the QUIC protocol ([RFC 9000](https://datatracker.ietf.org/doc/html/rfc9000), [RFC 9001](https://datatracker.ietf.org/doc/html/rfc9001), [RFC 9002](https://datatracker.ietf.org/doc/html/rfc9002)) in Go. It has support for HTTP/3 ([RFC 9114](https://datatracker.ietf.org/doc/html/rfc9114)), including QPACK ([RFC 9204](https://datatracker.ietf.org/doc/html/rfc9204)) and HTTP Datagrams ([RFC 9297](https://datatracker.ietf.org/doc/html/rfc9297)).

In addition to these base RFCs, it also implements the following RFCs:

* Unreliable Datagram Extension ([RFC 9221](https://datatracker.ietf.org/doc/html/rfc9221))
* Datagram Packetization Layer Path MTU Discovery (DPLPMTUD, [RFC 8899](https://datatracker.ietf.org/doc/html/rfc8899))
* QUIC Version 2 ([RFC 9369](https://datatracker.ietf.org/doc/html/rfc9369))
* QUIC Event Logging using qlog ([draft-ietf-quic-qlog-main-schema](https://datatracker.ietf.org/doc/draft-ietf-quic-qlog-main-schema/) and [draft-ietf-quic-qlog-quic-events](https://datatracker.ietf.org/doc/draft-ietf-quic-qlog-quic-events/))
* QUIC Stream Resets with Partial Delivery ([draft-ietf-quic-reliable-stream-reset](https://datatracker.ietf.org/doc/html/draft-ietf-quic-reliable-stream-reset-07))

Support for WebTransport over HTTP/3 ([draft-ietf-webtrans-http3](https://datatracker.ietf.org/doc/draft-ietf-webtrans-http3/)) is implemented in [webtransport-go](https://github.com/quic-go/webtransport-go).

Detailed documentation can be found on [quic-go.net](https://quic-go.net/docs/).

## Projects using quic-go

| Project                                                   | Description                                                                                                                                                       | Stars                                                                                               |
| ---------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------- |
| [AdGuardHome](https://github.com/AdguardTeam/AdGuardHome) | Free and open source, powerful network-wide ads & trackers blocking DNS server.                                                                                   | ![GitHub Repo stars](https://img.shields.io/github/stars/AdguardTeam/AdGuardHome?style=flat-square) |
| [algernon](https://github.com/xyproto/algernon)           | Small self-contained pure-Go web server with Lua, Markdown, HTTP/2, QUIC, Redis and PostgreSQL support                                                            | ![GitHub Repo stars](https://img.shields.io/github/stars/xyproto/algernon?style=flat-square)        |
| [caddy](https://github.com/caddyserver/caddy/)            | Fast, multi-platform web server with automatic HTTPS                                                                                                              | ![GitHub Repo stars](https://img.shields.io/github/stars/caddyserver/caddy?style=flat-square)       |
| [cloudflared](https://github.com/cloudflare/cloudflared)  | A tunneling daemon that proxies traffic from the Cloudflare network to your origins                                                                               | ![GitHub Repo stars](https://img.shields.io/github/stars/cloudflare/cloudflared?style=flat-square)  |
| [frp](https://github.com/fatedier/frp)                    | A fast reverse proxy to help you expose a local server behind a NAT or firewall to the internet                                                                   | ![GitHub Repo stars](https://img.shields.io/github/stars/fatedier/frp?style=flat-square)            |
| [go-libp2p](https://github.com/libp2p/go-libp2p)          | libp2p implementation in Go, powering [Kubo](https://github.com/ipfs/kubo) (IPFS) and [Lotus](https://github.com/filecoin-project/lotus) (Filecoin), among others | ![GitHub Repo stars](https://img.shields.io/github/stars/libp2p/go-libp2p?style=flat-square)     |
| [gost](https://github.com/go-gost/gost)                   | A simple security tunnel written in Go                                                                                                                        | ![GitHub Repo stars](https://img.shields.io/github/stars/go-gost/gost?style=flat-square)            |
| [Hysteria](https://github.com/apernet/hysteria)           | A powerful, lightning fast and censorship resistant proxy                                                                                                         | ![GitHub Repo stars](https://img.shields.io/github/stars/apernet/hysteria?style=flat-square)        |
| [Mercure](https://github.com/dunglas/mercure)             | An open, easy, fast, reliable and battery-efficient solution for real-time communications                                                                         | ![GitHub Repo stars](https://img.shields.io/github/stars/dunglas/mercure?style=flat-square)         |
| [OONI Probe](https://github.com/ooni/probe-cli)           | Next generation OONI Probe. Library and CLI tool.                                                                                                                 | ![GitHub Repo stars](https://img.shields.io/github/stars/ooni/probe-cli?style=flat-square)          |
| [reverst](https://github.com/flipt-io/reverst)            | Reverse Tunnels in Go over HTTP/3 and QUIC                                                                                                                        | ![GitHub Repo stars](https://img.shields.io/github/stars/flipt-io/reverst?style=flat-square) |
| [RoadRunner](https://github.com/roadrunner-server/roadrunner) | High-performance PHP application server, process manager written in Go and powered with plugins | ![GitHub Repo stars](https://img.shields.io/github/stars/roadrunner-server/roadrunner?style=flat-square) |
| [syncthing](https://github.com/syncthing/syncthing/)      | Open Source Continuous File Synchronization                                                                                                                       | ![GitHub Repo stars](https://img.shields.io/github/stars/syncthing/syncthing?style=flat-square)     |
| [traefik](https://github.com/traefik/traefik)             | The Cloud Native Application Proxy                                                                                                                                | ![GitHub Repo stars](https://img.shields.io/github/stars/traefik/traefik?style=flat-square)         |
| [v2ray-core](https://github.com/v2fly/v2ray-core)         | A platform for building proxies to bypass network restrictions                                                                                                    | ![GitHub Repo stars](https://img.shields.io/github/stars/v2fly/v2ray-core?style=flat-square)        |
| [YoMo](https://github.com/yomorun/yomo)                   | Streaming Serverless Framework for Geo-distributed System                                                                                                         | ![GitHub Repo stars](https://img.shields.io/github/stars/yomorun/yomo?style=flat-square)            |

If you'd like to see your project added to this list, please send us a PR.

## Release Policy

quic-go always aims to support the latest two Go releases.

## Contributing

We are always happy to welcome new contributors! We have a number of self-contained issues that are suitable for first-time contributors, they are tagged with [help wanted](https://github.com/quic-go/quic-go/issues?q=is%3Aissue+is%3Aopen+label%3A%22help+wanted%22). If you have any questions, please feel free to reach out by opening an issue or leaving a comment.
