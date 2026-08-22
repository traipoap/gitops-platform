FROM docker.io/library/golang:1.26.6-alpine3.24 AS backend-builder

WORKDIR /backend

RUN apk add --no-cache git gcc musl-dev sqlite-dev

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./

ARG VERSION=unknown

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w -X 'main.Version=${VERSION}'" -o /backend/app .

# ===================================================
FROM alpine
RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /backend

COPY --from=backend-builder /backend/app /backend/app
COPY --from=backend-builder /backend/exports/ /backend/exports/

EXPOSE 8080
ENV GIN_MODE=release
ENV PORT=8080
CMD ["/backend/app"]
