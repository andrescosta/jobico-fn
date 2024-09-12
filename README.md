# Jobico-fn

[![Go Report Card](https://goreportcard.com/badge/github.com/andrescosta/jobico)](https://goreportcard.com/report/github.com/andrescosta/jobico)

## Introduction

Jobico-fn is a multi-tenant compute service that enables the asynchronous execution of WebAssembly (WASM) functions in response to event triggers, offering scalable and efficient event-driven processing.

## Key Characteristics

- **Exploratory Nature**: Jobico-fn serves as an exploratory project, providing a platform for investigating different approaches to asynchronous computing technologies.

- **Multi-Tenancy Focus**: Jobico-fn's architecture is specifically designed to facilitate multi-tenancy, enabling the simultaneous operation of multiple isolated tenants on the platform.

- **Event-driven processing** : Break down jobs into smaller events for easier orchestration and control flow management.

- **Event Definition with JSON Schema**: Tenants can define events through JSON Schema, allowing for structured and dynamic event handling. Incoming requests undergo validation against the specified schema.

- **WASM-Compatible Language Support**: Implement processing logic in any language that compiles to WebAssembly, promoting platform independence..

# Software Stack

![alt](docs/img/stack.svg?)

# Architecture

![alt](docs/img/architecture.svg?)

For a more detailed overview of the system design and key architectural components, please refer to the [ARCHITECTURE.md](./ARCHITECTURE.md) file.

# Jobicolets

A **Jobicolet** is a WebAssembly (WASM) function designed to process an event and generate a result within the platform. It represents the executable logic that is dynamically loaded and executed by the platform. 

## SDK

The SDKs for creating Jobicolets are currently available in Go, Rust, JavaScript, and Python. They provide essential tools and functionality for developing functions.

- [Python](https://github.com/andrescosta/jobicolet-sdk-python)
- [JavaScript](https://github.com/andrescosta/jobicolet-sdk-js)
- [GO](https://github.com/andrescosta/jobicolet-sdk-go)
- [Rust](https://github.com/andrescosta/jobicolet-sdk-rust)


# Getting Started with Jobico-fn

This section provides instructions for compiling, starting, testing, and running.

## Local Go Environment

If you have [Go installed](https://go.dev/doc/install) on your machine, you can compile Jobico-fn directly from the source code:

``` bash
git clone 
cd jobico
make local
```

Alternatively, you can compile using [Docker](https://docs.docker.com/engine):

``` bash
git clone https://github.com/andrescosta/jobico-fn.git
cd jobico
make docker-build
```
## Service Management

1. Local
```bash
# Starting the services
scripts/startall.sh
#powershell: scripts\startall.ps1

# Stopping the services
scripts/stopall.sh
#powershell: scripts\stopall.ps1
```

2. Docker
``` bash
# Starting the environment
make docker-up

# Stopping the environment
make docker-stop
```

## Release

To release(build, e2e tests and lints) Jobico-fn locally, ensure you have the following dependencies installed:

- [Gcc](https://gcc.gnu.org/install/)
- [Go](https://go.dev/)
- [Make](https://www.gnu.org/software/make/)

And execute:

``` bash
make release
```

## Running Tests
After compiling and starting the services locally, you can run a set of happy path scenarios:

1. Install k6

``` bash
make k6
```

2. Run the test cases
   
``` bash
make perf1/local
```

``` bash
make perf2/local
```

[Learn more how testing works in Jobico-fn](TESTING.md)

# Platform Setup and Administration Guide

For instructions on platform setup and administration, please refer to the [GUIDE.md](./GUIDE.md) file.

# Operating Jobico-fn: Managing Your Deployment

For details on configuring and managing Jobico-fn, please refer to the [OPERATING.md](./OPERATING.md) file.

# Kubernetes

## Getting started

### Locally using Kind

**Requirements**

- [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)
- [Make](https://www.gnu.org/software/make/)

**Update the 'host' file.** 

Add the following entries:

```
   127.0.0.1 ctl
   127.0.0.1 recorder
   127.0.0.1 repo
   127.0.0.1 listener
   127.0.0.1 queue
   127.0.0.1 prometheus
   127.0.0.1 jaeger
```

**Self signed certificates**

```bash
   git clone https://github.com/andrescosta/jobico
   cd jobico

   # 1- Self signed certificates
   ## Creates the certifcates at k8s/certs
   make create-certs
   # Adds the certificates to the local storage
   make upload-certs-linux # Adds the certificates to the local store.
   #windows: make add-certs-windows (Warning: this command run as the admin user(opens the UAC dialog) and requires the user to accepts the changes.
```

#### Docker

**Requirements**

- [Docker](https://docs.docker.com/engine/install/)

**Create a cluster and test**

```bash
   # 1- Creates the cluster and deploy the application
   make kind
   
   # 1.1 - Wait until all ingresses are ready
   make wait-ings

   # 2- Local Test
   ## Builds k6 in perf/
   make k6
   ## Runs a basic scenario locally
   make perf1-k8s

   # 3- Deletes the  cluster
   make kind-delete
```

#### Podman

**Set enviroment variable for Kind**
If Docker and Kind are installed on the same machine, and Kind auto-detects Docker, set this environment variable to use Podman instead: 

```bash

KIND_EXPERIMENTAL_PROVIDER="podman"

```

**Create a cluster and test**

```bash
   # 0- Install Podman (only for Windows and Debian/Ubuntu, for others check: https://podman.io/docs/installation)
   make podman-install

   # 1- Init the Podman machine
   make podman-init

   # 1.1- Start the Podman machine
   make podman-start

   # 2- Creates a cluster and deploy the application
   make kind-podman
   
   # 2.1 - Wait until all ingresses are ready
   make wait-ings

   # 3- Local Test
   ## Builds k6 in perf/
   make k6
   ## Runs a basic scenario locally
   make perf1-k8s

   # 3- Deletes the cluster
   make kind-delete

   # 4- Resets Podman
   make podman-reset
```

# Roadmap

The roadmap can be accessed or queried at this location:

https://github.com/users/andrescosta/projects/3/views/1


### Short Term
- More examples
- Extend the capabilities of the Testing framework
- Improve error management

### Mid term
- Improvements to the Wasm runtime

### Long Term
- Quorum based replicated storage
- Durable computing exploration