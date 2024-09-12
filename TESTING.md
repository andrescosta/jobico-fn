# End to end testing

## Implementation overview

- **In-Memory Execution**: The testing framework operates entirely in memory, utilizing [grpc.bufconn](https://pkg.go.dev/google.golang.org/grpc/test/bufconn) for communication. 

- **In-Memory Database and Queue**: The testing framework employs in-memory implementations for the database and queue. This eliminates the need for external databases or message brokers during testing, streamlining the testing process.

- **Flexible Error Injection**: The testing framework allows for the injection of errors during test execution. This feature enables testers to simulate various error scenarios and edge cases, ensuring robustness and reliability in real-world conditions.

### High level architecture

![alt](docs/img/testing.svg?)

### Execution

```bash
cd jobico

# from Linux
./scripts/startall.sh
## ./scripts/startall.ps1 from PowerShell
## make dckr_up from Docker

# no coverage report
make test

# with coverage report
make test_coverage

# with HTML coverage report
make test_html
```

# Performance

This section guides you through the process of conducting comprehensive performance tests using K6.

## Execution
```bash
# Generate the K6 executable.
make k6

# Execute a simple scenario.
make perf1

# Execute a more complex scenario that involves streaming.
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

Jobico-fn supports integration with Prometheus and Jaeger for monitoring and tracing capabilities.

[Learn how to configure the observability stack](OPERATING.md#observability)