# Golang with Chi and SQLC

- Simple dummy API using SQLC, MySQL and Chi router
- Utility commands defined using Taskfile

### Tasks

#### `upsql`
- Builds and starts a MySQL docker container
- Custom Dockerfile pulls in a `schema.sql` file to set up tables and seed dummy data

#### `killsql`
- Stops and destroys the MySQL container.
- If destroyed and rebuilt, the setup script will re-run 

#### `sqlcli`
- Connects to MySQL CLI inside the docker container

#### `run`
- Runs the Go app

#### `ping`
- Ping Go app on localhost to check response

### SQLC
- `query.sql` and `schema.sql` are used to generate database access boilerplate via SQLC.
- Generated files are written into ./sqldb
- To regenerate after editing/adding to the sql files, use `sqlc generate` from project root















