# todo

A todo web application. Used for educational purposes with
https://github.com/sudomateo/terraform-training.

## Running

```
docker compose up
```

## Configuration

This service is configured using environment variables.

```
# Database host in the format HOST:PORT. When set the database will be used to
# store todo items.
TODO_DATABASE_HOST='postgres:5432'

# Database user.
TODO_DATABASE_USER='todo'

# Database password.
TODO_DATABASE_PASSWORD='todo'

# Database name.
TODO_DATABASE_NAME='todo'

# Database parameters.
TODO_DATABASE_PARAMETERS='sslmode=disable'

# Address the application will listen on in the format [IP]:PORT
TODO_ADDR=':8080'

# Log level to use. Valid levels are "trace", "debug", "info", "error", "warn".
TODO_LOG_LEVEL='info'

# Application version.
TODO_VERSION='1.0.0'
```
