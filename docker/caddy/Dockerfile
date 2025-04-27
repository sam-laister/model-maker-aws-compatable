FROM dunglas/mercure

LABEL maintainer="Sam Laister <saml@everbit.dev>"

RUN apk add --no-cache curl

CMD [ "/usr/bin/caddy", "run", "--config", "/etc/caddy/Caddyfile" ]
