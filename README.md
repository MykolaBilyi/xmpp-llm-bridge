# XMPP LLM Bridge

## Tasks

### run

```sh
go run cmd/app/main.go
```

### build-image

run: once
Inputs: TAG
Environment: TAG=dev

```sh
docker build -t ghcr.io/mykolabilyi/xmpp-llm-bridge:$TAG .
```

### run-image

Inputs: TAG
Environment: TAG=dev
Requires: build-image

```sh
docker run --rm --name xmpp-llm-bridge-dev \
    -e XMPP_JID=$XMPP_JID \
    -e XMPP_PASSWORD=$XMPP_PASSWORD \
    -p 8080:8080 \
    ghcr.io/mykolabilyi/xmpp-llm-bridge:$TAG
```
