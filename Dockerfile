#build stage
FROM golang:1.23-alpine AS builder
WORKDIR /go/src/xmpp-llm-bridge
COPY . .
RUN apk add --no-cache git
RUN go get ./...
RUN go build -o app ./cmd/app/main.go

#final stage
FROM alpine:latest
RUN apk --no-cache add wget ca-certificates
COPY --from=builder /go/src/xmpp-llm-bridge/app /app
COPY --from=builder /go/src/xmpp-llm-bridge/configs/default.yml /default.yml

LABEL Name=xmpp-llm-bridge
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD wget --timeout=5 http://localhost:8080/health || exit 1

CMD ["/app"]
