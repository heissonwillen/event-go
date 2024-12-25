# Event Go

Broadcast server sent events by sending HTTP calls.

![Docker Image Version](https://img.shields.io/docker/v/heissonwillen/event-go?sort=semver&label=Docker%20Image%20Version&logo=docker)
![Docker Pulls](https://img.shields.io/docker/pulls/heissonwillen/event-go)
![GitHub branch check runs](https://img.shields.io/github/check-runs/heissonwillen/event-go/main)
![Coveralls](https://img.shields.io/coverallsCoverage/github/heissonwillen/event-go)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/heissonwillen/event-go)

<!-- TODO: add demo video -->
<!-- TODO: add demo env -->
---

## Table of contents

- [Event Go](#event-go)
  - [Table of contents](#table-of-contents)
  - [Features](#features)
  - [Docker deployment](#docker-deployment)
  - [Installation](#installation)
    - [Prerequisites](#prerequisites)
    - [Clone the repository](#clone-the-repository)
  - [Build and run](#build-and-run)
    - [Local setup](#local-setup)
    - [Running tests](#running-tests)
  - [Configuration](#configuration)
    - [Environment variables](#environment-variables)
  - [Endpoints](#endpoints)
  - [Contributing](#contributing)
  - [License](#license)
  - [Community and Support](#community-and-support)

---

## Features

- Broadcast events: Server sent events are broadcasted once an authorized user sends an HTTP call.
- Event storage: New clients immediately receive the latest events once they connect.
- Dockerized Deployment: Simple setup with Docker and Docker Compose.  
- Production and Debug Modes: Easily switch between production and debug builds.  


---

## Docker deployment
Pulls the image from the Docker registry and runs it locally using [`compose.sample.yaml`](./compose.sample.yaml)
```bash
docker compose --file compose.sample.yaml up -d
```

---
## Installation

### Prerequisites

- Docker: [install Docker](https://docs.docker.com/get-docker/)  
- Docker Compose: [Install Docker Compose](https://docs.docker.com/compose/install/)  

### Clone the repository

```bash
git clone https://github.com/heissonwillen/event-go.git
cd event-go
```
---

## Build and run

### Local setup

To download the dependencies, and have the hot-reloading locally
```bash
go mod tidy && air
```
Alternatively, it's possible to run without hot-reloading (`air`)
```bash
go mod tidy && go run -v cmd/server/main.go
```

### Running tests

```bash
go test ./... -v
```

## Configuration

### Environment variables

| Variable              | Description                          | Default Value         |
| --------------------- | ------------------------------------ | --------------------- |
| `LISTEN_ADDR`         | API listen address                   | `:8080`               |
| `SQLITE_DB_PATH`      | SQLite DB path                       | `tmp/event-go.sqlite` |
| `BASIC_AUTH_USER`     | Username for `/authorized` endpoints | `admin`               |
| `BASIC_AUTH_PASSWORD` | Password for `/authorized` endpoints | `pass`                |

---

## Endpoints

| Method | Endpoint             | Description                                                                                                                                   |
| ------ | -------------------- | --------------------------------------------------------------------------------------------------------------------------------------------- |
| `POST` | `/authorized/events` | Broadcast a new event to all clients. The event is stored on the DB as a side-effect. It uses basic-auth - see environment variables section. |
| `GET`  | `/events`            | Receive stream of events. The latest event from the DB is broadcasted when a new client connects.                                             |


---

## Contributing

1. Fork the repository.  
2. Create a new branch: `git checkout -b my-feature-branch`  
3. Make your changes and add tests.  
4. Submit a pull request.  

---

## License

event-go is licensed under the [MIT License](LICENSE).

---

## Community and Support

- Issues: [GitHub issues](https://github.com/heissonwillen/event-go/issues)  
- Discussions: [GitHub discussions](https://github.com/heissonwillen/event-go/discussions)  
