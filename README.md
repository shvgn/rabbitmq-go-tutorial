# How to run

## 1. Launch RabbitMQ server

```sh
docker run -p 5672:5672 rabbitmq:3.7-management-alpine
```

## 2. Launch listeners

To test the delivery distribution and the durability, you cann kill all of them
and re-runs while tasks from `new_task.go` are launched.

```sh
go run worker.go
```

## 3. Send messages

```sh
go run new_task.go <arbitrary arguments>
```