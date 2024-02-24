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

### For a comprehensive overview of Jobico's architecture, [click here](ARCHITECTURE.md)

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
make dckr_build
```
## Service Management

1. Windows

```powershell
scripts/startall.ps1
```
```powershell
scripts/stopall.ps1
```

2. Linux

```bash
scripts/startall.sh
```
```bash
scripts/stopall.sh
```

3. Docker

``` bash
make dckr_up
```
``` bash
make dckr_stop
```

## Running Tests with K6
After compiling and starting the services locally, you can run a set of basic test cases using K6:

``` bash
make k6
cd perf
./k6.exe run events.js
```

And for a more comprensive set, run:

```bash
./k6.exe run eventsandstream.js
```
### For testing details on Jobico, [click here](TESTING.md)

# Platform Setup and Administration Guide

Explore how to work with the Jobico platform by checking out our comprehensive tutotial. Learn how to set up new jobs, upload wasm and schema files, and utilize administrative tools. If you're interested in diving deeper, click the link below to open the manual.

[Learn more with the In-Depth Guide](GUIDE.md)

# Operating Jobico: Managing Your Deployment

Explore key aspects of managing and maintaining your Jobico deployment. From monitoring job execution to configuring environment variables and performing health checks, this section covers essential topics to ensure the smooth operation of your Jobico environment.

[Learn More about Operating Jobico](OPERATING.md)


# Roadmap

The roadmap can be accessed or queried at this location:

https://github.com/users/andrescosta/projects/3/views/1


### Short Term
- Complete demo of a Jobicolet
- Extend the capabilities of the Testing framework
- Errors management

### Mid term
- Improvements to the Wasm runtime

### Long Term
- Distributed storage for the queue and control services
- Durable computing exploration

# Contact

For questions, feedback reach out to us at jobicowasm@gmail.com
