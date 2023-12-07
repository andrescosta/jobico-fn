* Small (by priority)
- Bug queue service
- Documentation
    - Introduction
    - Architecture
    - How to run
- More wasm examples (https://wasmer.io/posts/onyxlang-powered-by-wasmer) https://blog.jetbrains.com/kotlin/2023/12/kotlin-for-webassembly-goes-alpha/?utm_campaign=kotlin-wasm-alpha&utm_medium=social&utm_source=twitter
- Tests

- make file
    (check tools section)
    - tools: https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/tools
    - Vanity imports github.com/jcchavezs/porto/cmd/porto
    - Docker
        - compose 
- With pattern for services

- improve otel (
     with Functional Options Pattern,
     ADD LOGS: https://opentelemetry.io/docs/specs/otel/logs/, 
     Grafana Stack, 
     go.opentelemetry.io/otel/metric/noop, 
     go.opentelemetry.io/otel/trace/noop, 
     https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/config, 
     https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/exporters/autoexport , 
     https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/samplers/probability/consistent 

     https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/instrgen
     https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/samplers/probability/consistent
     
     correlation between logs, metrics and traces using id
)


(GRPC:
- Test
- Grpc Streaming (ByDirection or unidirectional)
- Grpc error handling https://grpc.io/docs/guides/error/
- Deadlines https://grpc.io/docs/guides/deadlines/
- Loadbalancing https://grpc.io/docs/guides/custom-backend-metrics/ - check clustering sestion below
- https://github.com/grpc-ecosystem/grpc-gateway/ (Not needed but useful)
- Backoff https://github.com/cenkalti/backoff/
- https://github.com/hanakoa/alpaca/issues/45)


- Service discovery


  (check https://github.com/mikehelmick/go-functional/blob/main/Makefile)
    
    - linters 
        (check https://github.com/golangci/golangci-lint/tree/master,
        https://github.com/mikehelmick/go-functional/blob/main/.golangci.yaml)
    
    - protobuff checkers:
    (check https://github.com/bufbuild/buf)


- Terraform

- reliability
    - performance https://research.swtch.com/testing
    - https://go.dev/blog/pprof , https://pkg.go.dev/go.uber.org/goleak
    - GC (https://medium.com/safetycultureengineering/analyzing-and-improving-memory-usage-in-go-46be8c3be0a8)
    - unit
    - integration using TestContainers 


- security (TLS, AuthN, AuthZ)

- Tools:

_ "github.com/atombender/go-jsonschema"
	_ "github.com/client9/misspell/cmd/misspell"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/jcchavezs/porto/cmd/porto"
	_ "github.com/wadey/gocovmerge"
	_ "golang.org/x/exp/cmd/gorelease"
	_ "golang.org/x/tools/cmd/stringer"
	_ "golang.org/x/vuln/cmd/govulncheck"

* Large
- health check manager who can monitor and stop the service, see: https://github.com/google/exposure-notifications-server/blob/main/internal/middleware/maintenance.go
https://github.com/google/exposure-notifications-server/blob/main/pkg/server/healthz.go

- Clustering: (LB, distributed queue, etc.) https://grpc.io/blog/grpc-load-balancing/ https://github.com/grpc/grpc/blob/master/doc/load-balancing.md, https://mykidong.medium.com/howto-grpc-java-client-side-load-balancing-using-consul-8f729668d3f8 

- WASM and WASI 

- Kubernetes
    - Operator (check https://medium.com/developingnodes/mastering-kubernetes-operators-your-definitive-guide-to-starting-strong-70ff43579eb9)

- Research:
https://github.com/rogpeppe/go-internal