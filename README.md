# TLS SNI router

[![license](https://img.shields.io/github/license/google/tlsrouter.svg?maxAge=2592000)](https://github.com/google/tlsrouter/blob/master/LICENSE) [![Travis](https://img.shields.io/travis/google/tlsrouter.svg?maxAge=2592000)](https://travis-ci.org/google/tlsrouter)  [![api](https://img.shields.io/badge/api-unstable-red.svg)](https://godoc.org/go.universe.tf/tlsrouter)

TLSRouter is a TLS proxy that routes connections to backends based on the TLS SNI (Server Name Indication) of the TLS handshake. It carries no encryption keys and cannot decode the traffic that it proxies.

This is not an official Google project.

## Installation

Install TLSRouter via `go get`:

```shell
go get go.universe.tf/tlsrouter
```

## Usage

TLSRouter requires a configuration file that tells it what backend to
use for a given hostname. The config file looks like:

```
# Basic hostname -> backend mapping
go.universe.tf localhost:1234

# DNS wildcards are understood as well.
*.go.universe.tf 1.2.3.4:8080

# DNS wildcards can go anywhere in name.
google.* 10.20.30.40:443

# RE2 regexes are also available
/(alpha|beta|gamma)\.mon(itoring)?\.dave\.tf/ 100.200.100.200:443

# If your backend supports HAProxy's PROXY protocol, you can enable
# it to receive the real client ip:port.

fancy.backend 2.3.4.5:443 PROXY

# Now there can be variables in config
# lucus.d1.tyd.us -> internal.nakagawa.lucus:443
# But remember, acme challenge will not work if this feature is used
/(?P<name>.*)\.d1\.tyd\.us/ /internal.nakagawa.${name}:443/

# Yes, there is also anonymous variable, but be careful with the brackets
# lucus.d2.tyd.us -> lucus.internal.nakagawa:443
/(.*)\.d2\.tyd\.us/ /$1.internal.nakagawa:443/
# lucus.d3.tyd.us -> lucus.internal.nakagawa:443
/((.*)\.d3\.tyd\.us)/ /$2.internal.nakagawa:443/

# Therefore, you can do some intersting things
# 443.lucus.d4.tyd.us -> internal.nakagawa.lucus:443
/(?P<port>^[0-9]{1,4})\.(?P<host>.*)\.d4\.tyd\.us/ /internal.nakagawa.${host}:${port}/
```

TLSRouter takes one mandatory commandline argument, the configuration file to use:

```shell
tlsrouter -conf tlsrouter.conf
```

Optional flags are:

 * `-listen <addr>`: set the listen address (default `:443`)
 * `-hello-timeout <duration>`: how long to wait for the start of the
   TLS handshake (default `3s`)
