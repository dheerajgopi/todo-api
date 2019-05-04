# A NON TRIVIAL TODO APPLICATION

### Requirements
- go1.12
- mysql 8
- golang-migrate cli tool (for database migrations)
`go get -tags 'mysql' -u github.com/golang-migrate/migrate/cmd/migrate`

### Run database migrations
- migrate up
`migrate -path migrations/sql -database {{mysql db connection string}} up`

- migrate down
`migrate -path migrations/sql -database {{mysql db connection string}} down`

- create new migration
`migrate create -dir migrations/sql -ext sql {{migration file name}}`

### Running the application
The application can be started using the `go run` command (`go run main.go`),
or by directly running the executable created using `go install` or `go build` command.
