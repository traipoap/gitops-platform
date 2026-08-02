FROM golang:1.26-alpine3.24 AS backend-builder

WORKDIR /backend
RUN apk add --no-cache git

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./

ARG VERSION=unknown
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X 'main.Version=${VERSION}'" -o /backend/app .

FROM alpine
RUN apk --no-cache add ca-certificates sqlite-libs \
    && addgroup -S app \
    && adduser -S app -G app

WORKDIR /backend
COPY --from=backend-builder --chown=app:app /backend/app /backend/app
COPY --from=backend-builder --chown=app:app /backend/users.db /backend/users.db
COPY --from=backend-builder --chown=app:app /backend/exports/ /backend/exports/

USER app
EXPOSE 8080
ENV GIN_MODE=release
ENV PORT=8080
CMD ["/backend/app"]
