# XMPP LLM Bridge

[![GitHub Release](https://img.shields.io/github/v/release/mykolabilyi/xmpp-llm-bridge)](https://github.com/mykolabilyi/xmpp-llm-bridge/releases) [![xc compatible](https://xcfile.dev/badge.svg)](https://xcfile.dev)

## Description

XMPP LLM Bridge is a Go service that connects XMPP chat networks with OpenAI's Large Language Models, enabling AI-powered conversations through XMPP clients. The service acts as an intelligent bot that can participate in XMPP conversations by processing incoming messages through OpenAI's API and responding with contextually relevant replies.

### Key Features

- **XMPP Integration**: Full XMPP client implementation using mellium.im/xmpp library with SASL authentication
- **OpenAI LLM Integration**: Seamlessly forwards chat messages to OpenAI's API and returns responses
- **Configurable**: YAML-based configuration with environment variable overrides
- **Docker Ready**: Containerized deployment with multi-stage builds

## Usage

The XMPP LLM Bridge can be deployed in several ways. All methods require three essential environment variables:

- `XMPP_JID`: Your bot's XMPP account (e.g., `bot@example.com`)
- `XMPP_PASSWORD`: Password for the XMPP account
- `OPENAI_API_KEY`: Your OpenAI API key

### Quick Start with Docker

Pull and run the latest image:

```bash
docker run \
    -e XMPP_JID=bot@example.com \
    -e XMPP_PASSWORD=your_password \
    -e OPENAI_API_KEY=sk-your-openai-key \
    ghcr.io/mykolabilyi/xmpp-llm-bridge:latest
```

### Docker Compose (Recommended)

Use this `docker-compose.yml`:

```yaml
services:
  bridge:
    image: ghcr.io/mykolabilyi/xmpp-llm-bridge:latest
    environment:
      - XMPP_JID=bot@example.com
      - XMPP_PASSWORD=your_password
      - OPENAI_API_KEY=sk-your-openai-key
    restart: unless-stopped
```

Start the service:

```bash
docker compose up -d
```

## Development

### Architecture

The service implements a hexagonal architecture pattern:

- **Core Business Logic**: Located in `/internal/app/` with dependency injection via context-based providers
- **Ports**: Interfaces in `/internal/ports/` define contracts for external services
- **Adapters**: External integrations in `/internal/adapters/` handle XMPP, OpenAI, and HTTP communications
- **Entities**: Domain models in `/internal/entities/` represent core data structures

The application runs two concurrent services:

1. **XMPP Client**: Handles incoming chat messages and processes them through the LLM
2. **HTTP Server**: Provides health check endpoint on port 8080 for monitoring

## Tasks

[xc] is recommended to simplify running tasks

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

[xc]: https://xcfile.dev/
