kubectl apply -f - <<EOF
apiVersion: opentelemetry.io/v1alpha1
kind: OpenTelemetryCollector
metadata:
  name: otel-collector-side-car
spec:
  mode: sidecar
  config: |
    receivers:
      otlp:
        protocols:
        http:
          endpoint: 0.0.0.0:4318
        grpc:
          endpoint: 0.0.0.0:4317

    exporters:
      logging:
      
    processors:
      batch:

    extensions:
      health_check:
      pprof:
      zpages:

    service:
      telemetry:
        logs:
          level: "debug"

      extensions: [pprof, zpages, health_check]

      pipelines:
        traces:
          receivers: [otlp]
          processors: [batch]
          exporters: [otlp]

        metrics:
          receivers: [otlp, otlp]
          processors: [batch]
          exporters: [otlp]
EOF