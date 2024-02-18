# Monitoring using OpenTelemetry

[TODO]

# Configuration

## Enviroment variables

## .env files

https://github.com/bkeepers/dotenv#what-other-env-files-can-i-use

### Command line

The configuration variables con be provided during the startup of the service as command line argument using the following format:

```bash
--env:[variable]=[value]
```
# Configuration

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

### Profiling
| Parameter | Description |
| --- | --- |
| prof.enabled | Enable or disable pprof profiler service. |
| pprof.addr | Address for the pprof profiler. |

### Observability
| Parameter | Description |
| --- | --- |
| metadata.enabled | Enable or disable service metadata. |
| grpc.healthcheck.freq | Frequency duration for validating the service's health. |


#### OpenTelemetry
| Parameter | Description |
| --- | --- |
| obs.enabled | Enable or disable the observability stack. |
| obs.exporter.trace.grpc.host | Grpc's otel expoter host address. |
| obs.exporter.metrics.http.host | HTTP's otel expoter host address. |
| obs.exporter.metrics.host.path | HTTP's otel expoter path. |
| obs.metrics.host | Enable or disable host's metrics. |
| obs.metrics.runtime | Enable or disable runtime's metrics. |

#### Logging
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

### Service specifics

#### Executor
| Parameter | Description |
| --- | --- |
|executor.timeout| Time the executor waits before fetching new events from the queue. |
|executor.maxproc| Number of processes to run in parallel when processing new events. |
