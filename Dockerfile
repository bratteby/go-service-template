FROM golang:1.19 AS builder

WORKDIR /go/src/example-service

COPY . .
RUN make lint test

# Build migrations cli
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/migrate ./cmd/migrate

# Build app
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/examplesvc ./cmd/example

# FROM alpine:3.17 
FROM scratch

COPY --from=builder /bin/migrate /bin/migrate
COPY --from=builder /go/src/example-service/migrations /migrations

COPY --from=builder /bin/examplesvc /app/examplesvc

CMD ["/app/examplesvc"]