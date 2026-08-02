FROM golang:1.26-alpine3.24 AS backend-builder

WORKDIR /backend
RUN apk add --no-cache git

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./

ARG VERSION=unknown
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X 'main.Version=${VERSION}'" -o /backend/app .

FROM alpine
RUN apk --no-cache add ca-certificates \
    && adduser -D -S -g '' app

WORKDIR /backend
COPY --from=backend-builder /backend/app /backend/app
COPY --from=backend-builder /backend/users.db /backend/users.db
COPY --from=backend-builder /backend/exports/ /backend/exports/

RUN chown -R app:app /backend
USER app
ENV GIN_MODE=release
EXPOSE 8080
ENV PORT=8080
ENV HOSTNAME="0.0.0.0"
CMD ["/backend/app"]
