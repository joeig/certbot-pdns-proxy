#!/bin/bash

set -e

PROXY_URL="https://127.0.0.1:8000/v1"
ACME_DOMAIN="_acme-challenge.${CERTBOT_DOMAIN}"

curl -s -X POST -n "${PROXY_URL}/authenticate?fqdn=${ACME_DOMAIN}&validation=${CERTBOT_VALIDATION}"

sleep 60
