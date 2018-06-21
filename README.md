# CertBot PowerDNS Proxy

Authentication with regex-based authorization for PowerDNS 4.1, designed for CertBot.

[![Build Status](https://travis-ci.org/joeig/certbot-pdns-proxy.svg?branch=master)](https://travis-ci.org/joeig/certbot-pdns-proxy)

## Usage

### Daemon

Copy the configuration template `config.dist.yaml` and launch the daemon:

~~~ bash
./certbot-pdns-proxy --config=/path/to/config.yaml
~~~

### CertBot

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
