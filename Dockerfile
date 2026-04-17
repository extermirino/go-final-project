FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/todo-app main.go

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app/todo-app .
COPY --from=builder /app/web ./web

EXPOSE 7540

ENV TODO_PORT=7540
ENV TODO_DBFILE=scheduler.db

CMD ["./todo-app"]