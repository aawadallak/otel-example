It’s necessary to install a few artifacts to make Open Telemetry works right on our local development setup.

1. Prometheus
2. SigNoz
3. HotROD

To install Prometheus, we gonna use helm, follow the steps below

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
```

```bash
helm install -n platform prometheus -f prometheus/values.yaml \
prometheus-community/prometheus
```

```bash
export POD_NAME=$(kubectl get pods --namespace platform -l "app=prometheus,component=server" -o jsonpath="{.items[0].metadata.name}")
kubectl --namespace platform port-forward $POD_NAME 9090
```

```bash
helm install --namespace platform destroyable signoz/signoz --values signoz/values.yaml
```

To install Signoz, we can use helm too, but it's necessary to add a new Otel collector config, so we can send metrics and logs to SigNoz.

```bash
helm repo add signoz https://charts.signoz.io
```

```bash
helm --namespace destroyable install oap signoz/signoz --set-file otelCollectorMetrics.config=signoz/values.yaml --dry-run --debug
```

```bash
helm --namespace destroyable install oap signoz/signoz --dry-run --debug
```

```bash
helm uninstall -n destroyable oap
```

```bash
export SERVICE_NAME=$(kubectl get svc --namespace platform -l "app.kubernetes.io/component=frontend" -o jsonpath="{.items[0].metadata.name}")
kubectl --namespace platform port-forward svc/$SERVICE_NAME 3301:3301
```