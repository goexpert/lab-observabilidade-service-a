module github.com/goexpert/lab-observabilidade-service-b

go 1.22.5

require (
	github.com/joho/godotenv v1.5.1
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.56.0
	go.opentelemetry.io/otel v1.31.0
	go.opentelemetry.io/otel/exporters/stdout/stdoutlog v0.7.0
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.31.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.31.0
	go.opentelemetry.io/otel/log v0.7.0
	go.opentelemetry.io/otel/sdk v1.31.0
	go.opentelemetry.io/otel/sdk/log v0.7.0
	go.opentelemetry.io/otel/sdk/metric v1.31.0
)

require (
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	go.opentelemetry.io/otel/metric v1.31.0 // indirect
	go.opentelemetry.io/otel/trace v1.31.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
)
