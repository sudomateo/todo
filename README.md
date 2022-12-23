# todo

An API to manage todo items. Primarily used for educational purposes.

## Running

```
# Native
TODO_DATABASE_URL='postgres://username:password@host:5432/database?sslmode=disable' -e TODO_AUTH_TOKEN='foo' go run ./app/todo-api

# Docker
docker run --rm -it -e TODO_DATABASE_URL='postgres://username:password@host:5432/database?sslmode=disable' -e TODO_AUTH_TOKEN='foo' -p 7836:7836 sudomateo/todo:latest

# Podman 
podman run --rm -it -e TODO_DATABASE_URL='postgres://username:password@host:5432/database?sslmode=disable' -e TODO_AUTH_TOKEN='foo' -p 7836:7836 sudomateo/todo:latest
```

## Configuration

This service is configured using environment variables.

```
# Required: The authentication token to use for API requests.
TODO_AUTH_TOKEN='foo'

# Required: The URL of the database to connect to.
TODO_DATABASE_URL='postgres://username:password@host:5432/database?sslmode=disable'

# Optional: The address for the API to listen on in the format [IP]:PORT
TODO_ADDR=':7836'

# Optional: The log level to use. Valid levels are "trace", "debug", "info",
# "error", "warn".
TODO_LOG_LEVEL='info'

# Optional: The version of the application.
TODO_VERSION='development'
```
