# grpctasks

gRPC сервис задач на Go. Код без комментариев.

## Требования
- Go 1.22+
- protoc
- protoc-gen-go
- protoc-gen-go-grpc

## Генерация
```
protoc -I proto --go_out=proto --go-grpc_out=proto proto/task.proto
```

## Запуск
```
go run ./cmd/server
```

## API
- CreateTask
- GetTask
- SetCompleted
- ListTasks
