FROM caddy:2.5.1-builder AS builder

RUN xcaddy build \
	--with github.com/caddyserver/cache-handler

FROM caddy:2.5.1

COPY --from=builder /usr/bin/caddy /usr/bin/caddy
