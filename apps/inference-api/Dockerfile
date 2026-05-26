FROM golang:1.23-alpine AS builder
WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
RUN go test ./... && go build -o /out/inference-api ./cmd/server

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /out/inference-api /app/inference-api
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/app/inference-api"]
