#!/bin/sh

./cloudflared -v &
# sudo cloudflared service install $CLOUDFLARE_TOKEN

/usr/bin/caddy run --config /etc/caddy/Caddyfile
