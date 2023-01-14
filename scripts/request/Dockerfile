FROM golang:alpine as builder

WORKDIR /usr/app
ADD go.mod go.sum ./ 
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o build/bin cmd/main.go

FROM golang:alpine
WORKDIR /usr/app
COPY --from=builder /usr/app/build/bin ./bin
CMD ["./bin"] 