FROM golang:1.19-bullseye AS builder

WORKDIR /usr/src/todo-api

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -v -o /usr/local/bin/todo-api ./app/todo-api

FROM ubuntu:jammy

COPY --from=builder /usr/local/bin/todo-api /usr/local/bin/todo-api

CMD ["/usr/local/bin/todo-api"]
