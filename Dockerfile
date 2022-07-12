# STEP 1: BUILD
FROM golang:1.18-alpine as build
WORKDIR /app
COPY . .
RUN go mod download
RUN export CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build cmd/proxy_server/proxy_server.go
RUN go build cmd/ws_server/ws_server.go

# STEP 2: RUN
FROM alpine as run
WORKDIR /app
COPY --from=build /app/proxy_server proxy_server
COPY --from=build /app/ws_server ws_server
COPY --from=build /app/.env .env
COPY --from=build /app/web web

ENTRYPOINT ["sh", "-c", "/app/proxy_server & /app/ws_server"]
