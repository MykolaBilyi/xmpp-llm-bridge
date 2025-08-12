#build stage
FROM golang:1.21-alpine AS builder
#FIXME: change to your project name
WORKDIR /go/src/__template__ 
COPY . .
RUN apk add --no-cache git
RUN go get ./...
RUN go build -o app ./cmd/app/main.go

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN apk add --no-cache bash
#FIXME: change to your project name
COPY --from=builder /go/src/__template__/app /app
COPY --from=builder /go/src/__template__/configs/default.yml /default.yml

#FIXME: change to your project name
LABEL Name=__template__
EXPOSE 8080
CMD ["/app"]
