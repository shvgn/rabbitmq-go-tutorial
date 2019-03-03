# How to run

## 1. Launch RabbitMQ server

```sh
docker run -p 5672:5672 rabbitmq:3.7-management-alpine
```

## 2. Launch listener

```sh
go run receive.go
```

## 3. Send messages

```sh
go run send.go
```