# Configuration

In Jobico-fn, various features and behaviors are controlled through OS environment variables. This configuration variables can be set not only using standard OS mechanisms but also through command line parameters and .env files.

![alt](docs/img/config.svg?)

## .env files

The loading of .env files follows the following principles: 

| File | Description |
| --- | --- |
|.env.[environment].local| Local overrides of environment variables. |
|.env.local| Local overrides.|
|.env.[environment]| Variables specific to each environment. |
|.env| Variables shared by all environments |

The current environment is determined by the value of the **APP_ENV** variable, which can take on the following values:
- development
- production
- test

These rules were defined with the help of https://github.com/bkeepers/dotenv#what-other-env-files-can-i-use

## Command line

Configuration variables can also be supplied during service startup as command line arguments, following this format:

```bash
[command] --env:[variable]=[value]
```

## List of configuration variables

### General
| Parameter | Description |
| --- | --- |
| APP_ENV | Environment mode for the application. Supported values: development, production,  test .|
| workdir | Directory where services store their information. |
| basedir | Directory from which services read .env configuration files. |
| [svc].addr | Address of the service. |
| [svc].host | Host port of the service. |
| [svc].dir| Name of the service's data directory. |
| dial.timeout | Timeout duration for dialing connections.|
| metadata.enabled | Enable or disable service metadata. |

### Profiling
| Parameter | Description |
| --- | --- |
| prof.enabled | Enable or disable pprof profiler service. |
| pprof.addr | Address for the pprof profiler. |

### Service specifics

#### Executor
| Parameter | Description |
| --- | --- |
|executor.timeout| Time the executor waits before fetching new events from the queue. |
|executor.maxproc| Number of processes to run in parallel when processing new events. |

# Observability

The observability stack in Jobico is implemented on top the OpenTelemetry client libraries and the Zerolog framework.Currently, metrics are sent to Prometheus, while traces are routed to Jaeger. 

![alt](docs/img/observability.svg)

## Environment Variables Overview

The observability stack is managed through a set of configuration variables. 

### Metrics and traces
| Parameter | Description |
| --- | --- |
| obs.enabled | Enable or disable the observability stack. |
| obs.exporter.trace.grpc.host | Grpc's otel expoter host address. |
| obs.exporter.metrics.http.host | HTTP's otel expoter host address. |
| obs.exporter.metrics.host.path | HTTP's otel expoter path. |
| obs.metrics.host | Enable or disable host's metrics. |
| obs.metrics.runtime | Enable or disable runtime's metrics. |

### Logging
| Parameter | Description |
| --- | --- |
| log.level | [Log's level.](https://github.com/rs/zerolog#leveled-logging)  |
| log.caller | Log the caller information. |
| log.console.enabled | Enable or disable console's logging. |
| log.console.exclude.timestamp | Exclude the timestamp when logging.|
| log.file.enabled | Enable or disable logging to a file. |
| log.file.name | File where logs will be written  |
| log.file.max.size | It is the maximum size in megabytes of the log file before it gets rotated. |
| log.file.max.backups | It is the maximum number of old log files to retain. |
| log.file.max.age | It is the maximum number of days to retain old log files. |

# Health Check

Jobico services exposes an API to check its internal health. This API will return the status code 200 if everything is correct or 500 if there is an error.

![alt](docs/img/healthchk.svg?)

Additionally, the HTTP health interfaces will return the follogin information in JSON format:

| Parameter | Type | Description |
| --- | --- | --- |
|Status| string | OK - If everything is normal. ERROR - If there is a problem with the service. |
|StartedAt| string | When the service as started. |
|Error| string | More details en case Status equal to ERROR|

## Configuration

| Parameter | Description |
| --- | --- |
| grpc.healthcheck.freq | Frequency duration for validating the service's health. |