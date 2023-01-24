k -n platform port-forward deployment/oap-signoz-frontend 3301:3301
k -n platform port-forward pods/oap-signoz-otel-collector-677f97c5d9-jsj5d 4317:4317

go run ./backend-parent/cmd/*.go &
go run ./backend-child/cmd/*.go &
go run ./scripts/request/cmd/*.go &
wait