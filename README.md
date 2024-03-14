# Jobico

[![Go Report Card](https://goreportcard.com/badge/github.com/andrescosta/jobico)](https://goreportcard.com/report/github.com/andrescosta/jobico)

## Introduction

Jobico: Multitenant Job Processing with WebAssembly

Jobico is an open-source platform designed for defining and processing jobs as a collection of events. Its core features revolve around multi-tenancy and language-agnostic execution via WebAssembly, making it suitable for various use cases requiring scalability and flexibility.

## Key Characteristics

- **Exploratory Nature**: Jobico serves as an exploratory project, providing a platform for investigating different approaches to asynchronous computing technologies.

- **Multi-Tenancy Focus**: Jobico's architecture is specifically designed to facilitate multi-tenancy, enabling the simultaneous operation of multiple isolated tenants on the platform.

- **Event-driven processing** : Break down jobs into smaller events for easier orchestration and control flow management.

- **Event Definition with JSON Schema**: Tenants can define events through JSON Schema, allowing for structured and dynamic event handling. Incoming requests undergo validation against the specified schema.

- **WASM-Compatible Language Support**: Implement processing logic in any language that compiles to WebAssembly, promoting platform independence..

# Software Stack

![alt](docs/img/stack.svg?)

# Architecture
Explore the architecture of Jobico to gain insights into its design principles and components. This section provides an overview of the system's structure and how its various elements interact to deliver powerful capabilities. plex workflows.

[Learn More about Jobico's architecture](ARCHITECTURE.md)

# Getting Started with Jobico

This guide provides instructions for compiling, starting, testing, and running.

## Local Go Environment

If you have [Go installed](https://go.dev/doc/install) on your machine, you can compile Jobico directly from the source code:

``` bash
git clone 
cd jobico
make local
```

Alternatively, you can compile using [Docker](https://docs.docker.com/engine):

``` bash
git clone https://github.com/andrescosta/jobico.git
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

Requirements

for testinginstall gcc

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

[Learn more how testing works in Jobico](TESTING.md)

# Platform Setup and Administration Guide

Explore how to work with the Jobico platform by checking out our comprehensive tutotial. Learn how to set up new jobs, upload wasm and schema files, and utilize administrative tools. If you're interested in diving deeper, click the link below to open the manual.

[Learn more with the In-Depth Guide](GUIDE.md)

# Operating Jobico: Managing Your Deployment

Explore key aspects of managing and maintaining your Jobico deployment. From monitoring job execution to configuring environment variables and performing health checks, this section covers essential topics to ensure the smooth operation of your Jobico environment.

[Learn More about Operating Jobico](OPERATING.md)

# Kubernetes

Running Jobico in Kubernetes marks the initial phase of a more ambitious project to develop a highly available WebAssembly-based serverless platform. This project includes implementing a quorum-replicated data store, enhancing caching components, fortifying GRPC communication for robustness, and exploring a deeper integration by creating a K8s operator. 

## Getting started

### Kind

To run Jobico locally using Kind, ensure you have the following dependencies installed:

- [Docker](https://docs.docker.com/engine/install/)
- [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)
- [Make](https://www.gnu.org/software/make/)


**Update the 'host' file:**

1- Add the following entries:

```
   127.0.0.1 ctl
   127.0.0.1 recorder
   127.0.0.1 repo
   127.0.0.1 listener
   127.0.0.1 queue
   127.0.0.1 prometheus
   127.0.0.1 jaeger
```
2- Create a cluster and test

```bash
   git clone https://github.com/andrescosta/jobico
   cd jobico

   # 1- Self signed certificates
   ## Creates the certifcates at k8s/certs
   make create-certs
   # Adds the certificates to the local storage
   make upload-certs-linux # Adds the certificates to the local store.
   #windows: make add-certs-windows (Warning: this command run as the admin user(opens the UAC dialog) and requires the user to accepts the changes.
   
   # 2- Creates the cluster and deploy the application
   make kind
   
   # 3- Local Test
   ## Builds k6 in perf/
   make k6
   ## Runs a basic scenario locally
   make perf1-k8s

   # 4- Deletes the  cluster
   make kind-delete
```

# Roadmap

The roadmap can be accessed or queried at this location:

https://github.com/users/andrescosta/projects/3/views/1


### Short Term
- More complex Jobicolet examples
- Extend the capabilities of the Testing framework
- Improve error management

### Mid term
- Improvements to the Wasm runtime

### Long Term
- Quorum based replicated storage 
- Durable computing exploration

# Contact

For questions, feedback reach out to us at jobicowasm@gmail.com
