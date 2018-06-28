#!/bin/bash

set -e

PROXY_URL="https://127.0.0.1:8000/v1"
ACME_DOMAIN="_acme-challenge.${CERTBOT_DOMAIN}"

curl -s -X DELETE -n "${PROXY_URL}/cleanup?fqdn=${ACME_DOMAIN}"
