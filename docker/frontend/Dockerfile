FROM node:26.5.1-alpine3.24 AS builder

WORKDIR /app

COPY frontend/package*.json ./
RUN npm ci

COPY frontend/ ./
RUN npm run build

# ===================================================
FROM node:26.7.0-alpine3.24

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/dist /app/dist
COPY --from=builder /app/node_modules /app/node_modules
COPY --from=builder /app/package*.json /app/
COPY --from=builder /app/astro.config.mjs /app/

ENV PORT=4321
ENV HOSTNAME="0.0.0.0"

EXPOSE 4321
RUN npm install -g npm@latest

CMD ["npm", "run", "preview", "--", "--host", "0.0.0.0"]
