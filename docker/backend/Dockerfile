FROM golang:1.26-alpine3.24 AS backend-builder

WORKDIR /app
RUN apk add --no-cache git

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./

ARG VERSION=unknown
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X 'main.Version=${VERSION}'" -o /app/backend .

FROM alpine

RUN apk --no-cache add ca-certificates \
    && adduser -D -S -g '' app

WORKDIR /app
COPY --from=backend-builder /app/backend /app/backend

USER app
EXPOSE 8080
ENV PORT=8080
ENV HOSTNAME="0.0.0.0"
CMD ["/app/backend"]
