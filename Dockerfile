FROM golang:1.20-bullseye AS builder

WORKDIR /usr/src/todo

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -v -o /usr/local/bin/todo .

FROM ubuntu:jammy

COPY --from=builder /usr/local/bin/todo /usr/local/bin/todo

CMD ["/usr/local/bin/todo"]
