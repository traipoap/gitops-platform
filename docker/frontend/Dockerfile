FROM node:26.5.1-alpine3.24 AS builder

WORKDIR /app

COPY frontend/package*.json ./
RUN npm ci

COPY frontend/ ./
RUN npm run build

# ===================================================
FROM node:26.7.0-alpine3.24

RUN apk --no-cache add ca-certificates \
 && addgroup -g 10001 app \
 && adduser -u 10001 -G app -s /sbin/nologin -D app

WORKDIR /app

COPY --from=builder --chown=app:app /app/dist /app/dist
COPY --from=builder --chown=app:app /app/node_modules /app/node_modules
COPY --from=builder --chown=app:app /app/package*.json /app/
COPY --from=builder --chown=app:app /app/astro.config.mjs /app/

USER app

ENV PORT=4321
ENV HOSTNAME="0.0.0.0"

ARG ALLOWED_HOSTS
ENV ALLOWED_HOSTS=${ALLOWED_HOSTS:-localhost,127.0.0.1}
ARG API_PROXY_TARGET
ENV API_PROXY_TARGET=${API_PROXY_TARGET:-}

EXPOSE 4321

CMD ["npm", "run", "preview", "--", "--host", "0.0.0.0"]
