FROM --platform=linux/amd64 node:25.8.1-alpine3.23 AS builder

WORKDIR /app

COPY frontend/package*.json ./
RUN npm ci

COPY frontend/ ./
RUN npm run build

FROM --platform=linux/amd64 node:25.8.1-alpine3.23

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/dist /app
COPY --from=builder /app/node_modules /app/node_modules
COPY --from=builder /app/package*.json /app/

ENV PORT=4321
ENV HOSTNAME="0.0.0.0"

# Create src/pages directory for Astro (required for dev/preview)
RUN mkdir -p /app/src/pages
RUN mkdir -p /app/dist

EXPOSE 4321
RUN npm install -g npm@11.15.0
CMD ["npm", "run", "preview"]
