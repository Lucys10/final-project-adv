FROM golang:1.16 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0  GOOS=linux GOARCH=amd64 go build -o ./controller ./cmd/api/main.go

FROM scratch
COPY --from=builder /app/controller /usr/bin/controller
ENTRYPOINT ["/usr/bin/controller"]