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

## V1 API

All API paths are prefixed with `/api/v1`.

### Authentication

Authenticate by sending your API token in the `Authorization` header.

```
Authorization: Bearer <token>
```

### Todos

#### List Todos

`GET /todos`

##### Query Paramters

- `completed` - Optional. Retrieve todos based on their completion. Valid values are `true` and `false`.
- `priority` - Optional. Retrive todos with this priority. Valid values are `low`, `medium`, and `high`.

##### Request Body

N/A

##### Response Body

```json
[
  {
    "id": "00000000-0000-0000-0000-000000000000",
    "text": "Example todo.",
    "priority": "low",
    "completed": false,
    "time_created": "2022-12-23T20:48:08.273566323Z",
    "time_updated": "2022-12-23T20:48:08.273566323Z"
  },
  {
    "id": "00000000-0000-0000-0000-000000000000",
    "text": "Example todo.",
    "priority": "low",
    "completed": false,
    "time_created": "2022-12-23T20:48:08.273566323Z",
    "time_updated": "2022-12-23T20:48:08.273566323Z"
  }
]
```

#### Create a Todo

`POST /todos`

##### Request Body

```json
{
  "text": "Example todo.",
  "priority": "low"
}
```

##### Response Body

```json
{
  "id": "00000000-0000-0000-0000-000000000000",
  "text": "Example todo.",
  "priority": "low",
  "completed": false,
  "time_created": "2022-12-23T20:48:08.273566323Z",
  "time_updated": "2022-12-23T20:48:08.273566323Z"
}
```

#### Get a Todo

`GET /todos/:id`

##### Request Body

N/A

##### Response Body

```json
{
  "id": "00000000-0000-0000-0000-000000000000",
  "text": "Example todo.",
  "priority": "low",
  "completed": false,
  "time_created": "2022-12-23T20:48:08.273566323Z",
  "time_updated": "2022-12-23T20:48:08.273566323Z"
}
```

#### Update a Todo

`PATCH /todos/:id`

##### Request Body

```json
{
  "text": "Example todo updated.",
  "priority": "high",
  "completed": true
}
```

##### Response Body

```json
{
  "id": "00000000-0000-0000-0000-000000000000",
  "text": "Example todo updated.",
  "priority": "high",
  "completed": true,
  "time_created": "2022-12-23T20:48:08.273566323Z",
  "time_updated": "2022-12-23T20:48:08.273566323Z"
}
```

#### Delete a Todo

`DELETE /todos/:id`

##### Request Body

N/A

##### Response Body

N/A

### Version 

#### Get Version

`GET /version`

##### Request Body

N/A

##### Response Body

```json
{
  "version": "development"
}
```
