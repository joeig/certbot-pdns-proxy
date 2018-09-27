# CertBot PowerDNS Proxy

Authentication with regex-based authorization for PowerDNS 4.1, designed for CertBot.

[![Build Status](https://travis-ci.org/joeig/certbot-pdns-proxy.svg?branch=master)](https://travis-ci.org/joeig/certbot-pdns-proxy)
[![Go Report Card](https://goreportcard.com/badge/github.com/joeig/certbot-pdns-proxy)](https://goreportcard.com/report/github.com/joeig/certbot-pdns-proxy)

## Setup

### Install from source

You need `go` and `GOBIN` in your `PATH`. Once that is done, install dyndns-pdns using the following command:

~~~ bash
go get -u github.com/joeig/certbot-pdns-proxy
~~~

After that, copy [`config.dist.yml`](config.dist.yml) to `config.yml`, replace the default settings and run the binary:

~~~ bash
certbot-pdns-proxy -config=/path/to/config.yml
~~~

If you're intending to add the application to your systemd runlevel, you may want to take a look at [`scripts/certbot-pdns-proxy.service`](scripts/certbot-pdns-proxy.service).

## Usage

### Use in combination with CertBot

Deploy `scripts/authenticator.sh` and `scripts/cleanup.sh` on your servers and change the proxy URL.

You need to add your API credentials to `~/.netrc` as following:

~~~ text
machine 127.0.0.1
  login foo
  password bar
~~~ 

Pass the scripts to CertBot:

~~~ bash
certbot certonly --manual --preferred-challenges=dns --manual-auth-hook /path/to/authenticator.sh --manual-cleanup-hook /path/to/cleanup.sh -d secure.example.com
~~~

## FAQ

- **Q: How can I increase the SOA's serial automatically?**  
  A: Set the `SOA-EDIT-API` metadata to a value of your choice, for instance `pdnsutil set-meta example.com SOA-EDIT-API INCEPTION-INCREMENT`. There might be a [default setting](https://github.com/PowerDNS/pdns/issues/6173) in the future.
