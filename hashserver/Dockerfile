FROM golang:1.16 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0  GOOS=linux GOARCH=amd64 go build -o ./hashservice ./cmd/grpcserver/server.go

FROM scratch
COPY --from=builder /app/hashservice /usr/bin/hashservice
ENTRYPOINT ["/usr/bin/hashservice"]