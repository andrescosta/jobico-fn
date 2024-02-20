# End to end testing

## Implementation overview

The testing framework for Jobico is designed to provide comprehensive and efficient testing capabilities while offering several advantages:

- **In-Memory Execution**: The testing framework operates entirely in memory, utilizing [grpc.bufconn](https://pkg.go.dev/google.golang.org/grpc/test/bufconn) for communication. This approach ensures fast and efficient testing without the need for external resources or dependencies.

- **In-Memory Database and Queue**: To further enhance speed and simplicity, the testing framework employs in-memory implementations for the database and queue. This eliminates the need for external databases or message brokers during testing, streamlining the testing process.

- **Flexible Error Injection**: The testing framework allows for the injection of errors during test execution. This feature enables testers to simulate various error scenarios and edge cases, ensuring robustness and reliability in real-world conditions.

By leveraging these advantages, the testing framework ensures thorough testing of Jobico's functionality while maintaining efficiency and flexibility.

### High level architecture

![alt](docs/img/testing.svg?)

### Execution

Execute Jobico test cases directly from the command line with ease. Follow these simple steps to validate the functionality of Jobico's components and ensure seamless operation of your deployment.

```bash
cd jobico

# Start Jovico using a bash shell
./scripts/startall.sh
## ./scripts/startall.ps1 in PowerShell
## make dckr_up for Docker

# no coverage report
make test

# with a coverage report
make test_coverage

# with an HTML coverage report
make test_html
```

# Performance

Harness the power of performance testing with Jobico's K6 implementation. This section guides you through the process of conducting comprehensive performance tests, enabling you to evaluate Jobico's scalability and performance under varying workloads.

## Execution
```bash
cd jobico

# Generate the K6 executable
make k6

# Execute scenario 1
make perf1

# Execute scenario 2
make perf2
```

### Adjusting test parameters

Fine-tune your performance tests by customizing the number of virtual users (VUs) and iterations. 

```javascript
// in events.js or eventsandstream.js change:

export let options = {
  vus: 1, // The number of virtual users.
  iterations: 1 //  the total number of times the test script will be executed.
};
```

### Profiling
Turn on pprof before executing tests to collect valuable profiling information. By enabling pprof, you can gather insights into the performance of Jobico's components and identify potential areas for optimization.

The following environment variables control the profiling capabilities:

| Parameter | Description |
| --- | --- |
| prof.enabled | Enable or disable pprof profiler service. |
| pprof.addr | Address for the pprof profiler. |

### Telemetry

Jobico supports integration with Prometheus and Jaeger for monitoring and tracing capabilities.

The following environment variables control this integration:

| Parameter | Description |
| --- | --- |
| obs.enabled | Enable or disable the observability stack. |
| obs.exporter.trace.grpc.host | Grpc's otel expoter host address. |
| obs.exporter.metrics.http.host | HTTP's otel expoter host address. |
| obs.exporter.metrics.host.path | HTTP's otel expoter path. |
| obs.metrics.host | Enable or disable host's metrics. |
| obs.metrics.runtime | Enable or disable runtime's metrics. |


### More information on Telemetry and Profiling configuration can be found in this [guide](OPERATING.md)

